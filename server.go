package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

//Adresses
const (
	IP   = "192.168.1.47" // IP local
	PORT = "3569"         // Port utilisé
)

func gestionErreur(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Lancement du serveur ...")
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", IP, PORT))
	gestionErreur(err)
	var clients []net.Conn // tableau de clients
	for {
		conn, err := ln.Accept()
		if err == nil {
			clients = append(clients, conn) //quand un client se connecte on le rajoute à notre tableau
			go clientHandler(conn)
		}
		gestionErreur(err)
		fmt.Println("Client connected from :", conn.RemoteAddr())
	}
}

func clientHandler(conn net.Conn) {
	fmt.Println(conn)
	read(conn)

}

func read(conn net.Conn) {
	stream := bufio.NewReader(conn)
	for true {
		rawHeader, err := stream.ReadString('\n') //read the header ==> series of caracteres
		header := strings.Fields(rawHeader)       // parse the header into sub categories such as type, rows, columns operation
		fmt.Println(header)
		types, err := strconv.Atoi(header[0])
		rowA, err := strconv.Atoi(header[1])
		columnA, err := strconv.Atoi(header[2])
		rowB, err := strconv.Atoi(header[3])
		columnB, err := strconv.Atoi(header[4])
		buffer := make([]byte, rowA*columnA*8)
		read, err := io.ReadFull(stream, buffer)
		fmt.Println(read)
		byteMatrixA := byteSliceToByteMatrix(types, rowA, columnA, buffer)
		floatMatrixA := byteMatrixToFloatMatrix(byteMatrixA, rowA, columnA)
		buffer = make([]byte, rowB*columnB*8)
		read, err = io.ReadFull(stream, buffer)
		byteMatrixB := byteSliceToByteMatrix(types, rowB, columnB, buffer)
		floatMatrixB := byteMatrixToFloatMatrix(byteMatrixB, rowB, columnB)
		//fmt.Println(floatMatrixA)
		//fmt.Println(floatMatrixB)
		start := time.Now()
		floatMatrixC := routinesHandeler(floatMatrixA, floatMatrixB, 1)
		end := time.Since(start)
		fmt.Println("elapsed", end)
		bytesMatrixC := FloatMatrixToBytes(floatMatrixC)
		//fmt.Println(floatMatrixC)
		rowC := len(floatMatrixC)
		columnC := len(floatMatrixC[0])
		ClientHeader := makeHeader("2", rowC, columnC) // MAKE PROPER HEADER FOR CLIENT
		rawClientHeader := []byte(ClientHeader)
		fmt.Println(ClientHeader, rawClientHeader)
		send(conn, rawClientHeader)
		send(conn, bytesMatrixC)
		fmt.Println(err)
	}
}

func makeHeader(types string, row int, column int) string {
	rowstring := strconv.Itoa(row)
	columnstring := strconv.Itoa(column)
	header := types + " " + rowstring + " " + columnstring + "\n"
	return header
}

func send(conn net.Conn, message []byte) {
	conn.Write(message)
}

func BytesToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func Float64ToBytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func FloatMatrixToBytes(matrix [][]float64) []byte {
	var cell []byte
	var bytematrix []byte
	rows := len(matrix)
	columns := len(matrix[0])
	for x := 0; x < rows; x++ {
		for y := 0; y < columns; y++ {
			cell = Float64ToBytes(matrix[x][y])
			bytematrix = append(bytematrix, cell...)
		}
	}
	return bytematrix
}

func byteSliceToByteMatrix(types int, rows int, columns int, slice []byte) [][][]byte {
	var matrix [][][]byte
	var matrixLine [][]byte
	var cell []byte
	var bytes int

	if types == 2 {
		bytes = 8
	} else {
		bytes = 8
	}

	start := 0
	end := bytes
	for x := 0; x < rows; x++ {
		matrixLine = [][]byte{}
		for y := 0; y < columns; y++ {
			cell = []byte{}
			for j := start; j < end; j++ {
				cell = append(cell, slice[j])
			}
			start = start + bytes
			end = end + bytes
			matrixLine = append(matrixLine, cell)
		}
		matrix = append(matrix, matrixLine)
	}
	return matrix
}

func byteMatrixToFloatMatrix(byteMatrix [][][]byte, rows int, columns int) [][]float64 {
	var floatMatrix [][]float64
	var floatLine []float64
	var floatCell float64
	for x := 0; x < rows; x++ {
		floatLine = []float64{}
		for y := 0; y < columns; y++ {
			floatCell = BytesToFloat64(byteMatrix[x][y])
			floatLine = append(floatLine, floatCell)
		}
		floatMatrix = append(floatMatrix, floatLine)
	}
	return floatMatrix
}

func LBL(start int, end int, aColumns int, bColumns int, matrixA [][]float64, matrixB [][]float64, sc *sync.WaitGroup, ch chan [][]float64) {
	var results [][]float64
	for line := start; line < end; line++ {
		var lineRes []float64
		for x := 0; x < bColumns; x++ {
			var result float64
			result = 0
			for y := 0; y < aColumns; y++ {
				result = result + (matrixA[line][y] * matrixB[y][x])
			}
			lineRes = append(lineRes, result)
		}
		results = append(results, lineRes)
	}
	ch <- results
	sc.Done()
}

// this function will synchronize the routines
func routinesHandeler(matrixA [][]float64, matrixB [][]float64, cpuCount int) [][]float64 {
	var sc sync.WaitGroup

	start := 0
	end := len(matrixA)
	step := end / cpuCount
	rest := end % cpuCount
	aColumns := len(matrixA[0])
	bColumns := len(matrixB[0])
	chanPool := make([]chan [][]float64, cpuCount)
	for i := 0; i < len(chanPool); i++ {
		chanPool[i] = make(chan [][]float64, 1)
	}
	for routines := 0; routines < cpuCount-1; routines++ {
		sc.Add(1)
		fmt.Println(start, start+step)
		go LBL(start, start+step, aColumns, bColumns, matrixA, matrixB, &sc, chanPool[routines])
		start = start + step

	}
	fmt.Println(start, start+step+rest)
	go LBL(start, start+step+rest, aColumns, bColumns, matrixA, matrixB, &sc, chanPool[cpuCount])
	sc.Wait()
	matrix := [][]float64{}
	for i := 0; i < len(chanPool); i++ {
		matrix = append(matrix, <-chanPool[i]...)
	}
	return matrix
}

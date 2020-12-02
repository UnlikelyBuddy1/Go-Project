package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"os"
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

func main() {
	// Connexion au serveur
	clientHandler()
}

func clientHandler() {
	var wg sync.WaitGroup
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", IP, PORT))
	fmt.Println(conn)
	fmt.Println(err)
	wg.Add(1)
	go read(conn, &wg)

	matrixA := [][]float64{}
	matrixB := [][]float64{}
	matrixA = readFileWithReadString("matrixA.txt")
	matrixB = readFileWithReadString("matrixB.txt")
	types := "2"
	operation := "0"
	rowA := len(matrixA)
	columnA := len(matrixA[0])
	rowB := len(matrixB)
	columnB := len(matrixB[0])
	header := makeHeader(types, rowA, columnA, rowB, columnB, operation)
	rawHeader := []byte(header)
	send(conn, rawHeader)
	bytes := FloatMatrixToBytes(matrixA)
	fmt.Println(len(bytes))
	send(conn, bytes)
	bytes = FloatMatrixToBytes(matrixB)
	fmt.Println(len(bytes))
	send(conn, bytes)
	wg.Wait()
}
func makeHeader(types string, rowA int, columnA int, rowB int, columnB int, operation string) string {
	rowAstring := strconv.Itoa(rowA)
	rowBstring := strconv.Itoa(rowB)
	columnAstring := strconv.Itoa(columnA)
	columnBstring := strconv.Itoa(columnB)
	header := types + " " + rowAstring + " " + columnAstring + " " + rowBstring + " " + columnBstring + " " + operation + "\n"
	return header
}
func read(conn net.Conn, wg *sync.WaitGroup) {
	stream := bufio.NewReader(conn)
	for true {
		rawHeader, err := stream.ReadString('\n') //read the header ==> series of caracteres
		header := strings.Fields(rawHeader)       // parse the header into sub categories such as type, rows, columns operation
		fmt.Println(err)
		fmt.Println(rawHeader, header)
		types, err := strconv.Atoi(header[0])
		rowC, err := strconv.Atoi(header[1])
		columnC, err := strconv.Atoi(header[2])
		buffer := make([]byte, rowC*columnC*8)
		read, err := io.ReadFull(stream, buffer)
		fmt.Println(read)
		byteMatrixC := byteSliceToByteMatrix(types, rowC, columnC, buffer)
		floatMatrixC := byteMatrixToFloatMatrix(byteMatrixC, rowC, columnC)
		fmt.Println(floatMatrixC)
		fmt.Println(err)
		time.Sleep(2 * time.Second)
	}
	wg.Done()
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

func send(conn net.Conn, message []byte) {
	conn.Write(message)
}

// reading files
func readFileWithReadString(fileAdrr string) [][]float64 {
	matrix := [][]float64{}
	file, err := os.Open(fileAdrr)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var line string
	index := 0
	for {
		line, err = reader.ReadString('\n')

		if err != nil && err != io.EOF {
			break
		}
		// Process the line
		matrixLine := []float64{}
		words := strings.Fields(line)
		for word := 0; word < len(words); word++ {
			if s, err := strconv.ParseFloat(words[word], 64); err == nil {
				matrixLine = append(matrixLine, s)
			}
		}
		matrix = append(matrix, matrixLine)
		index++
		if err != nil {
			break
		}
	}
	if err != io.EOF {
		fmt.Printf(" > Failed with error: %v\n", err)
	}
	return matrix
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

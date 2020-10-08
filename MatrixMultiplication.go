package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	matrixA := [][]float64{}
	matrixB := [][]float64{}
	cpu := 1
	matrixA = readFileWithReadString("matrixA.txt")
	matrixB = readFileWithReadString("matrixB.txt")
	for index := 1; index < 12; index++ {
		start := time.Now()
		routinesHandeler(matrixA, matrixB, cpu)
		elapsed := time.Since(start)
		fmt.Println(cpu, elapsed)
		cpu = cpu * 2
	}

}

//LBL is a function Happy now !!!!
func LBL(start int, end int, aColumns int, bColumns int, matrixA [][]float64, matrixB [][]float64, sc *sync.WaitGroup) {
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
	sc.Done()
}

// this function will synchronize the routines
func routinesHandeler(matrixA [][]float64, matrixB [][]float64, cpuCount int) {
	var sc sync.WaitGroup
	start := 0
	end := len(matrixA)
	step := end / cpuCount
	aColumns := len(matrixA[0])
	bColumns := len(matrixB[0])
	for routines := 0; routines < cpuCount; routines++ {
		sc.Add(1)
		go LBL(start, start+step, aColumns, bColumns, matrixA, matrixB, &sc)
		start = start + step
	}
	sc.Wait()
}

// reading files
func readFileWithReadString(fileAdrr string) [][]float64 {
	matrix := [][]float64{}
	matrixLine := []float64{}
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
	print(matrix)
	if err != io.EOF {
		fmt.Printf(" > Failed with error: %v\n", err)
	}
	return matrix
}

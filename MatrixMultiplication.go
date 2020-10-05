package main

import "fmt"

func main() {

	matrixA := [][]uint{{1, 2}, {3, 4}, {5, 6}}
	matrixB := [][]uint{{1, 2, 3}, {4, 5, 6}}
	LBL(0, len(matrixA), len(matrixA[0]), len(matrixB[0]), matrixA, matrixB)
	var s []int
	s = append(s, 0)
}

//LBL is a function Happy now !!!!
func LBL(start int, end int, aColumns int, bColumns int, matrixA [][]uint, matrixB [][]uint) {
	var results [][]uint
	for line := start; line < end; line++ {
		var lineRes []uint
		for x := 0; x < bColumns; x++ {
			var result uint
			result = 0
			for y := 0; y < aColumns; y++ {
				result = result + (matrixA[line][y] * matrixB[y][x])
			}
			lineRes = append(lineRes, result)
		}
		results = append(results, lineRes)
	}
	fmt.Println(results)
}

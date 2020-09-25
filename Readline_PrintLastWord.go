package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	readFileWithReadString()
}

func readFileWithReadString() {
	file, err := os.Open("/home/adri/Desktop/GoProject/TextForGo.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var line string
	for {
		line, err = reader.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}
		// Process the line
		words := strings.Fields(line)
		lenght := len(words)
		fmt.Printf(words[lenght-1])
		if err != nil {
			break
		}
	}
	if err != io.EOF {
		fmt.Printf(" > Failed with error: %v\n", err)
	}
	return
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

type WorkerSetting struct {
	blockSize int
}

func worker(index int, inputFile string, config WorkerSetting, report chan string, control chan bool) {
	input, err := os.OpenFile(inputFile, os.O_RDONLY, 0600)
	check(err)
	defer func() {
		err := input.Close()
		check(err)
	}()

	scanner := bufio.NewScanner(input)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++

		if (lineNumber % config.blockSize) == 0 {
			report <- fmt.Sprintf("worker %d completed %d lines", index, lineNumber)
		}
	}

	if scanner.Err() != nil {
		report <- scanner.Err().Error()
	}
	control <- true
}

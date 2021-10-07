package main

import (
	"bufio"
	"fmt"
	"os"
)

type WorkerSetting struct {
	blockSize int
}

type WorkerChannels struct {
	control chan string
}

func worker(index int,
	inputFile string,
	config WorkerSetting,
	reportChannel chan string,
	waitChannel chan bool,
	doneChannel chan bool,
	channels WorkerChannels) {
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

		// PROCESS LINE HERE

		// LOAD OUTPUT TO BUFFER HERE

		// finished a block, report and waiting
		if (lineNumber % config.blockSize) == 0 {
			reportChannel <- fmt.Sprintf("worker %d completed %d lines", index, lineNumber)
			waitChannel <- true

			// proceed to write file output
			// wait for write command
			for <-channels.control != "write" {
			}

			// WRITE FILE HERE

			// report writing complete
			channels.control <- "done"
		}
	}

	if scanner.Err() != nil {
		reportChannel <- scanner.Err().Error()
	}
	doneChannel <- true
}

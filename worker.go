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
	done    chan bool
}

func worker(index int,
	inputFile string,
	config WorkerSetting,
	reportChannel chan string,
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

		// finished a block, report and waiting
		if (lineNumber % config.blockSize) == 0 {
			reportChannel <- fmt.Sprintf("worker %d completed %d lines", index, lineNumber)
			//channels.control <- fmt.Sprintf("%d", index)

			// proceed to write file output
			//for <-channels.control != fmt.Sprintf("%d", index) {}

			// WRITE FILE HERE
			//channels.report <- fmt.Sprintf("worker %d write file", index)
			//channels.control <- fmt.Sprintf("%d", index+1)

			// proceed to next block
			//for <-channels.control != "go" {}
		}
	}

	if scanner.Err() != nil {
		reportChannel <- scanner.Err().Error()
	}
	doneChannel <- true
}

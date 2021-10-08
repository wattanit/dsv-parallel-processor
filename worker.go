package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

type WorkerSetting struct {
	blockSize  int
	numProcess int
}

type WorkerChannels struct {
	Report  chan string
	Wait    chan bool
	Done    chan bool
	Control chan string
}

func worker(index int,
	inputFile string,
	spec Spec,
	config WorkerSetting,
	channels WorkerChannels) {
	input, err := os.OpenFile(inputFile, os.O_RDONLY, 0600)
	check(err)
	defer func() {
		err := input.Close()
		check(err)
	}()

	scanner := bufio.NewScanner(input)
	lineNumber := 0
	var outputBuffer []string

	for scanner.Scan() {

		if (lineNumber % config.numProcess) == index {
			// PROCESS LINE HERE
			line := scanner.Text()
			if filter(line, spec) {
				outputBuffer = append(outputBuffer, line)
			}
		}

		lineNumber++
		// finished a block, report and waiting
		if (lineNumber % config.blockSize) == 0 {
			channels.Report <- fmt.Sprintf("worker %d completed %d lines", index, lineNumber)
			//channels.Wait <- true

			// proceed to write file output
			// wait for write command
			//for <-channels.Control != "write" {
			//}

			func(outputBuffer []string, spec Spec) {
				outputFilePath := path.Join(spec.Output.OutputFile, fmt.Sprintf("%d", index))
				outputFile, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_APPEND, 0644)
				check(err)
				defer func() {
					err := outputFile.Close()
					check(err)
				}()

				for _, line := range outputBuffer {
					cells := strings.Split(line, spec.Input.Separator)
					_, err := outputFile.WriteString(strings.Join(cells, spec.Output.Separator) + "\n")
					check(err)
				}
			}(outputBuffer, spec)

			// report writing complete
			channels.Report <- fmt.Sprintf("worker %d wrote %d filtered lines", index, len(outputBuffer))
			outputBuffer = nil
			//channels.Control <- "done"
		}
	}

	if scanner.Err() != nil {
		channels.Report <- scanner.Err().Error()
	}
	channels.Done <- true
}

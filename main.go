package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func main() {
	numProcessFlagPtr := flag.Int("p", 1, "number of worker processes")
	blockSizeFlagPtr := flag.Int("block-size", 100000, "processing block size in number of lines")

	flag.Parse()

	numProcess := *numProcessFlagPtr
	blockSize := *blockSizeFlagPtr
	fmt.Println(numProcess)
	fmt.Println(blockSize)

	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("[Error] missing input path")
		return
	}

	workerConfig := WorkerSetting{
		blockSize: blockSize,
	}

	inputPath := args[0]
	inputFiles, _ := ioutil.ReadDir(inputPath)

	for _, f := range inputFiles {
		filePath := path.Join(inputPath, f.Name())
		fmt.Println(filePath)

		// spawn workers
		reportChannel := make(chan string, 100)
		doneChannel := make(chan bool, numProcess)
		waitChannel := make(chan bool, numProcess)
		controlChannels := []WorkerChannels{}

		for i := 0; i < numProcess; i++ {
			controlChannels = append(controlChannels, WorkerChannels{
				control: make(chan string),
			})
		}

		for i := 0; i < numProcess; i++ {
			fmt.Println(i)
			go worker(i, filePath, workerConfig, reportChannel, waitChannel, doneChannel, controlChannels[i])
		}

		// process monitoring
		doneWorkers := 0
		waitWorkers := 0
		for {
			if len(reportChannel) > 0 {
				r := <-reportChannel
				fmt.Println(r)
			}
			if waitWorkers == numProcess {
				waitWorkers = 0
				for i := 0; i < numProcess; i++ {
					_, err := fmt.Fprintf(os.Stdout, "Writing worker %d\n", i)
					check(err)
					controlChannels[i].control <- "write"
					_, err = fmt.Fprintf(os.Stdout, "Waiting for worker %d to finish writing\n", i)
					check(err)
					<-controlChannels[i].control
					_, err = fmt.Fprintf(os.Stdout, "Worker %d done writing\n", i)
					check(err)
				}
			}
			if doneWorkers == numProcess {
				fmt.Println(doneWorkers)
				break
			}
			select {
			case <-doneChannel:
				doneWorkers++
			case <-waitChannel:
				waitWorkers++
			default:

			}
		}
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

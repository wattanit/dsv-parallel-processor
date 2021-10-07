package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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
		controlChannels := []WorkerChannels{}

		for i := 0; i < numProcess; i++ {
			controlChannels = append(controlChannels, WorkerChannels{
				control: make(chan string),
				done:    make(chan bool),
			})
		}

		for i := 0; i < numProcess; i++ {
			fmt.Println(i)
			go worker(i, filePath, workerConfig, reportChannel, doneChannel, controlChannels[i])
		}

		// process monitoring
		doneWorkers := 0
		//waitWorkers := 0
		for {

			//if doneWorkers == numProcess {
			//	break
			//}
			//if waitWorkers == numProcess{
			//	fmt.Println("Reseting waitWorkers")
			//	waitWorkers = 0
			//	fmt.Println(waitWorkers)
			//	for i := 0; i < numProcess; i++ {
			//		channels.control <- "0"
			//	}
			//
			//	fmt.Println("Shout writing")
			//	channels.control <- "write"
			//	fmt.Println("Waiting response")
			//	//<-channels.control
			//}

			//for i := 0; i < numProcess; i++{
			//	channels.control <- "go"
			//}
			//}
			if len(reportChannel) > 0 {
				r := <-reportChannel
				fmt.Println(r)
			}

			if doneWorkers == numProcess {
				fmt.Println(doneWorkers)
				break
			}
			select {
			case <-doneChannel:
				doneWorkers++
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

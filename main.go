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
		report := make(chan string)
		control := make(chan bool)
		for i := 0; i < numProcess; i++ {
			fmt.Println(i)
			go worker(i, filePath, workerConfig, report, control)
		}

		// process monitoring
		doneWorkers := 0
		for {
			if doneWorkers == numProcess {
				break
			}
			select {
			case r := <-report:
				fmt.Println(r)
			case <-control:
				doneWorkers++
				fmt.Println(doneWorkers)
			}
		}
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

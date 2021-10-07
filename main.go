package main

import (
	"flag"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"io/ioutil"
	"os"
	"path"

	//"io/ioutil"
	"log"
)

type Spec struct {
	Input struct {
		FilePaths []string
		Directory string
		Separator string
	}
	Output struct {
		OutputFile string
		Separator  string
	}
	Filter []struct {
		Column     int
		ColumnType string
		Values     []string
	}
}

func main() {
	// parse CLI flags
	numProcessFlagPtr := flag.Int("p", 1, "number of worker processes")
	blockSizeFlagPtr := flag.Int("block-size", 100000, "processing block size in number of lines")

	flag.Parse()

	numProcess := *numProcessFlagPtr
	blockSize := *blockSizeFlagPtr

	workerConfig := WorkerSetting{
		blockSize:  blockSize,
		numProcess: numProcess,
	}
	fmt.Println(workerConfig)

	// load spec file
	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("[Error] missing spec file")
		return
	}

	specFilePath := args[0]
	specFile, err := os.ReadFile(specFilePath)
	check(err)

	var spec Spec
	err = toml.Unmarshal(specFile, &spec)
	check(err)

	// create input file list
	var inputPaths []string
	if spec.Input.Directory != "" {
		fmt.Println("Use directory input")
		inputFiles, _ := ioutil.ReadDir(spec.Input.Directory)
		for _, f := range inputFiles {
			filePath := path.Join(spec.Input.Directory, f.Name())
			inputPaths = append(inputPaths, filePath)
		}
	} else if len(spec.Input.FilePaths) > 0 {
		fmt.Println("Use path list input")
		inputPaths = spec.Input.FilePaths
	} else {
		log.Fatalf("[Error] No input file specified in %s", specFilePath)
	}

	// init loggers
	debugLog := log.New(os.Stdout, "[DEBUG]", log.Ldate|log.Ltime)
	infoLog := log.New(os.Stdout, "[INFO]", log.Ldate|log.Ltime)

	// loop through files
	for _, f := range inputPaths {
		debugLog.Printf("Processing file %s", f)

		// spawn workers
		reportChannel := make(chan string, 100)
		doneChannel := make(chan bool, numProcess)
		waitChannel := make(chan bool, numProcess)
		var controlChannels []chan string

		for i := 0; i < numProcess; i++ {
			controlChannels = append(controlChannels, make(chan string))
		}

		for i := 0; i < numProcess; i++ {
			debugLog.Printf("Spawning worker %d", i)
			go worker(i, f, workerConfig, WorkerChannels{
				reportChannel,
				waitChannel,
				doneChannel,
				controlChannels[i],
			})
		}

		// process monitoring
		doneWorkers := 0
		waitWorkers := 0

		for {
			// display worker report
			if len(reportChannel) > 0 {
				r := <-reportChannel
				infoLog.Println(r)
			}
			// sync file writing queue
			if waitWorkers == numProcess {
				waitWorkers = 0
				for i := 0; i < numProcess; i++ {
					debugLog.Printf("Writing worker %d\n", i)
					controlChannels[i] <- "write"

					debugLog.Printf("Waiting for worker %d to finish writing\n", i)
					<-controlChannels[i]

					debugLog.Printf("Worker %d done writing\n", i)
				}
			}

			// break monitoring loop when all workers are done
			if doneWorkers == numProcess {
				infoLog.Printf("%s completed", f)
				break
			}

			// channel monitoring
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

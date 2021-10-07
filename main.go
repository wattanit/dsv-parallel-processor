package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func worker(index int, inputFile string, report chan string, control chan bool) {
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
	}

	if scanner.Err() != nil {
		report <- scanner.Err().Error()
	}
	report <- fmt.Sprintf("worker %d completed %d lines", index, lineNumber)

	//scanner := bufio.NewScanner(input)
	//lineNumber := 0
	//fucker := scanner.Scan()
	//for fucker{
	//	lineNumber++
	//}
	//report <- fmt.Sprintf("worker %d completed %d lines",index, lineNumber)

	//for lineNumber:=0; lineAvailable; lineNumber++ {
	//	//line := scanner.Text()
	//	scanner.Text()
	//	//report <- line
	//	lineAvailable = scanner.Scan()
	//}
	control <- true
}

func main() {
	numProcessFlagPtr := flag.Int("p", 1, "number of worker processes")

	flag.Parse()

	numProcess := *numProcessFlagPtr
	fmt.Println(numProcess)

	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("[Error] missing input path")
		return
	}

	inputPath := args[0]
	inputFiles, _ := ioutil.ReadDir(inputPath)

	for _, f := range inputFiles {
		filePath := path.Join(inputPath, f.Name())
		fmt.Println(filePath)

		report := make(chan string)
		control := make(chan bool)
		for i := 0; i < numProcess; i++ {
			fmt.Println(i)
			go worker(i, filePath, report, control)
		}
		for i := 0; i < numProcess; i++ {
			fmt.Println(<-report)
		}

		//input, err := os.OpenFile(filePath, os.O_RDONLY, 0600)
		//check(err)
		//defer input.Close()
		//
		//scanner := bufio.NewScanner(input)
		//lineNumber := 0
		//for scanner.Scan(){
		//	lineNumber++
		//}
		//fmt.Println(lineNumber)

		//doneWorkers := 0
		//for {
		//
		//	//select {
		//	//case <-control:
		//	//	doneWorkers++
		//	//default:
		//	//	if doneWorkers == 4 {
		//	//		break
		//	//	}
		//	//}
		//}
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

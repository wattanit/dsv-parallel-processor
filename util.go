package main

import (
	"bufio"
	"log"
	"os"
)

func checkFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		log.Fatal(err)
	}
	return false
}

func readValueFile(inputFile string) []string {
	input, err := os.OpenFile(inputFile, os.O_RDONLY, 0600)
	check(err)
	defer func() {
		err := input.Close()
		check(err)
	}()

	var outputList []string

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		outputList = append(outputList, scanner.Text())
	}
	return outputList
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func isin(a string, list []string) bool {
	for _, item := range list {
		if item == a {
			return true
		}
	}
	return false
}

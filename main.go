package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
)

var wg sync.WaitGroup

func main() {
	/*
		this is to ensure that the runtime uses all the
		available physical processors on the machine to run the program
	*/
	runtime.GOMAXPROCS(runtime.NumCPU())

	file, err := os.OpenFile("file.txt", os.O_APPEND, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}

}

// countLines counts the number of lines in a file
func countLines(file *os.File) (int, error) {
	var count int
	lineBreak := byte('\n')
	tempBuffer := make([]byte, bufio.MaxScanTokenSize)

	for {
		bufSize, err := file.Read(tempBuffer)
		if err != nil && err != io.EOF {
			return 0, err
		}

		var bufferPosition int

		for {
			i := bytes.IndexByte(tempBuffer[bufferPosition:], lineBreak)
			if i == -1 || bufferPosition == bufSize {
				break
			}
			bufferPosition += i + 1
			count++
		}

		if err == io.EOF {
			break
		}
	}
	return count, nil
}

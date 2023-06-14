/*
	fscram is a program to reshuffle files line by line

	It works by making goroutines that have buffers which split the lines equally to then rewrite them
	into a new file randomly

	The max number of goroutines allowed is 1024

	If the file has n lines and n is less than or equal 1024, n goroutines are created and each goroutine
	processes one line.

	If the number of lines in the file exceeds 1024, a sharing formula is used

	The sharing formula involves determining a maximum size for the buffers of each goroutine

	This maximum size can be determined by the formula: math.Floor(n/1024) + 1 if n mod 1024 != 0 else 1
	or math.Floor(n/1024) if n mod 1024 == 0
			n = number of lines in file
*/

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"sync"
	"time"
)

// FileControl provides a mechanism that allows safe concurrent access to the lines in a file.
// It encapsulates a bufio Scanner and a mutex to control access
type FileControl struct {
	file  *bufio.Scanner
	mutex sync.Mutex
}

const MAXGONUM = 1024

var wg sync.WaitGroup
var goNum int     // goNum is the number of goroutines required to handle the file <= 1024
var bufferLen int // bufferLen is the size of the buffer that each goroutines would require to store lines

var regulatorChannel = make(chan int)

func main() {
	/*
		this is to ensure that the runtime uses all the
		available physical processors on the machine to run the program
	*/
	runtime.GOMAXPROCS(runtime.NumCPU())

	var fileName string = "file.txt"

	file, err := os.Open("file.txt")
	if err != nil {
		log.Fatal(err)
	}

	numLines, err := countLines(fileName) // number of lines in file
	if err != nil {
		log.Fatal(err)
	}

	if numLines < 1024 {
		goNum = numLines
	} else {
		goNum = 1024
	}

	bufferLen = int(math.Floor(float64(numLines) / 1024))
	if numLines%1024 != 0 {
		bufferLen++
	}

	fileReader := bufio.NewScanner(file)

	fileController := &FileControl{fileReader, sync.Mutex{}}
	wg.Add(goNum)

	for i := 0; i < goNum; i++ {
		go randomize(i, regulatorChannel, fileController)
	}

	wg.Wait()
}

// countLines counts the number of lines in a file
func countLines(fileName string) (int, error) {
	var count int
	lineBreak := byte('\n')
	tempBuffer := make([]byte, bufio.MaxScanTokenSize)

	file, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}

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

func randomize(idNum int, supervisor chan int, fileController *FileControl) {
	defer wg.Done()
	lines := []string{}

	for i := 0; i < bufferLen; i++ {
		time.Sleep(750)
		fileController.mutex.Lock()
		fileController.file.Scan()
		newLine := fileController.file.Text()
		fileController.mutex.Unlock()

		lines = append(lines, newLine)
		runtime.Gosched()
	}

	fmt.Println(lines)
}

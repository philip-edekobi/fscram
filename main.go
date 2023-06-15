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
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"

	filePkg "fscram/file"
)

const MAXGONUM = 1024

var wg sync.WaitGroup
var goNum int     // goNum is the number of goroutines required to handle the file <= 1024
var bufferLen int // bufferLen is the size of the buffer that each goroutines would require to store lines

func main() {
	/*
		this is to ensure that the runtime uses all the
		available physical processors on the machine to run the program
	*/
	runtime.GOMAXPROCS(runtime.NumCPU())

	if len(os.Args) <= 2 {
		if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
			// display help
		} else {
			fmt.Println("Usage: fscram <input file> <output file>")
			os.Exit(1)
		}
	}
	var inFileName string = os.Args[1]
	var outFileName string = os.Args[2]

	file, err := os.Open(inFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	outFile, err := os.Create(outFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	numLines, err := filePkg.CountLines(inFileName) // number of lines in file
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

	inputFileController := filePkg.NewReaderSync(fileReader)
	outputFileController := filePkg.NewFileSync(outFile)

	wg.Add(goNum)

	for i := 1; i <= goNum; i++ {
		go randomize(i, inputFileController, outputFileController)
	}

	wg.Wait()
}

func randomize(idNum int, inFileCtrl *filePkg.ReaderSync, outFileCtrl *filePkg.FileSync) {
	defer wg.Done()
	lines := []string{}

	time.Sleep(time.Duration(rand.Float64()*100*float64(goNum/idNum)) * time.Millisecond)

	for i := 0; i < bufferLen; i++ {
		time.Sleep(time.Duration(rand.Float64()*5000*float64(goNum/idNum)) * time.Millisecond)

		inFileCtrl.Mutex.Lock()

		inFileCtrl.Reader.Scan()
		newLine := inFileCtrl.Reader.Text()

		inFileCtrl.Mutex.Unlock()

		lines = append(lines, newLine)
		runtime.Gosched()
	}

	time.Sleep(time.Duration(rand.Float64()*1000*float64(goNum/idNum)) * time.Millisecond)

	for _, line := range lines {
		time.Sleep(time.Duration(rand.Float64()*6000*float64(goNum/idNum)) * time.Millisecond)

		outFileCtrl.Mutex.Lock()
		outFileCtrl.File.WriteString(line + "\n")
		outFileCtrl.Mutex.Unlock()
	}
}

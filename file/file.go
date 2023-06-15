package file

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"sync"
)

// FileSync provides a mechanism that allows safe concurrent access to a file,
// It does this by providing a mutex that controls access to the file.
type FileSync struct {
	File  *os.File
	Mutex sync.Mutex
}

// ReaderSync provides a mechanism that allows safe concurrent access to a reader via a mutex
type ReaderSync struct {
	Reader *bufio.Scanner
	Mutex  sync.Mutex
}

// NewFileSync creates a new FileSync struct and returns a pointer to it
func NewFileSync(file *os.File) *FileSync {
	return &FileSync{file, sync.Mutex{}}
}

// NewReaderSync created a new ReaderSync struct and returns a pointer to it
func NewReaderSync(reader *bufio.Scanner) *ReaderSync {
	return &ReaderSync{reader, sync.Mutex{}}
}

// CountLines counts the number of lines in a file
func CountLines(fileName string) (int, error) {
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

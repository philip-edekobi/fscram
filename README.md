# Fscram
A utility for reshuffling the lines of files

## Synopsis

Fscram is a program to reshuffle files line by line.

It works by making goroutines. In these goroutines, a buffer is present which holds a certan number of lines such that all goroutines share the lines equally.

The max number of goroutines allowed is 1024.

If the file has n lines and n is less than or equal 1024, n goroutines are created and each goroutine processes one line.

If the number of lines in the file exceeds 1024, a sharing formula is used.

The sharing formula involves determining a maximum size for the buffers of each goroutine.

This maximum size can be determined by the formula:

    math.Floor(n/1024) + 1 if n mod 1024 != 0 else 1
    
or:

    math.Floor(n/1024) if n mod 1024 == 0

where: 
    
    n = number of lines in file

## Usage

### help

```bash
$ fscram -h | fscram --help
```

This command displays a similar explanation to this readme and the commands available

### Normal Usage

```bash
$ fscram inputFile outputFile
```

Fscram accepts two arguments. The first is the input file which is the original file you wish to scramble. The second is the destination or output file.

## Installation

To install this package run:

```bash
$ go install github.com/philip-edekobi/fscram@latest
```

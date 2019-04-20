package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	var f *os.File
	if len(os.Args) == 2 {
		var err error
		f, err = os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	} else {
		f = os.Stdin
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	result, err := addFrequency(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d\n", *result)
}

func addFrequency(r io.Reader) (*int, error) {
	scanner := bufio.NewScanner(r)
	n := 0
	for scanner.Scan() {
		i, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, err
		}
		n += i
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &n, nil
}

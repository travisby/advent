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

	// result, err := addFrequency(f)
	result, err := firstDoubleFrequency(f)
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

func firstDoubleFrequency(r io.Reader) (*int, error) {
	// well, we need to be able to repeat the list now so I guess we're putting it in memory........
	changes := []int{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		i, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, err
		}
		changes = append(changes, i)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	n := 0
	// byte is just something easy to ignore
	encounteredNs := map[int]byte{0: 0x00}

	// okay, let's start iterating through the freqs
	i := 0
	for {
		n += changes[i%len(changes)]
		if _, ok := encounteredNs[n]; ok {
			return &n, nil
		}
		encounteredNs[n] = 0x00

		i++
	}

	return nil, fmt.Errorf("Did not encounter a double frequency")
}

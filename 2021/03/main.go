package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

	countBits := [][2]uint{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		for i, c := range scanner.Text() {
			// I wish there was a way to reduce this without splitting the for-loop into two
			for i >= len(countBits) {
				countBits = append(countBits, [2]uint{})
			}
			if c == '0' {
				countBits[i][0]++
			} else if c == '1' {
				countBits[i][1]++
			} else {
				log.Fatalf("Unexpected Nanary digit in binary: %d", c)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// uint so we can play with bitshifting
	var gamma uint64
	// epsilon is really the ^gamma, but for the _custom int size_
	// doing a go native ^ is just giving us the NOT of a uint64 :sweat-smile:
	// since we're only calculating 1, it's fine to just... calculate it
	var epsilon uint64

	for i, b := range countBits {
		if b[1] > b[0] {
			gamma += 1 << uint64(len(countBits)-i-1)
		} else {
			fmt.Printf("0")
			epsilon += 1 << uint64(len(countBits)-i-1)
		}
	}

	log.Printf("Part 1: %d", gamma*epsilon)

}

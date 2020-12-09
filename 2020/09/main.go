package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

const preambleSize = 25

var ErrInvalidNext = errors.New("The next number is not the sum of preceding numbers")

// a greedyDecoderRing remembers all of the content it's actually seen
// but still acts as a ring buffer
// just with a secret old storage ability
type greedyDecoderRing struct {
	ring    []int
	zeroPos int
	sz      int
}

func newDecoderRing(preamble []int) *greedyDecoderRing {
	return &greedyDecoderRing{preamble, 0, len(preamble)}
}

func (d *greedyDecoderRing) next(i int) error {
	km := d.toKeyMap()

	for k := range km {
		// if there exists the key that is i - our current number
		// (aka i-k + k = i), where (i-k) exists
		if _, ok := km[i-k]; ok {
			// valid!

			d.ring = append(d.ring, i)
			d.zeroPos++
			return nil
		}
	}

	return fmt.Errorf("%w, %d", ErrInvalidNext, i)
}

var ErrWrongStuffing = errors.New("Someone stuffed something baaad in here")

func (d *greedyDecoderRing) toKeyMap() map[int]struct{} {
	m := map[int]struct{}{}
	for _, k := range d.ring[len(d.ring)-d.sz:] {
		m[k] = struct{}{}
	}
	return m
}

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

	scanner := bufio.NewScanner(f)

	// first let's fill the preamble
	preamble := make([]int, 0, preambleSize)
	for len(preamble) < preambleSize && scanner.Scan() {
		i, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		preamble = append(preamble, i)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	} else if len(preamble) != preambleSize {
		log.Fatalf("Not enough content to consume: %+v", preamble)
	}

	decoder := newDecoderRing(preamble)
	var i *int
	// now continue scanning!
	for scanner.Scan() || i != nil {
		in, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		if err := decoder.next(in); err != nil {
			i = &in
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if i != nil {
		log.Printf("Part 1: %d", *i)
	} else {
		log.Fatal("Missing Part 1")
	}

	for j := 0; j < len(decoder.ring); j++ {
		smallest := decoder.ring[j]
		largest := decoder.ring[j]
		for k := j; k < len(decoder.ring) && smallest*largest < *i; k++ {
			if decoder.ring[k] < smallest {
				smallest = decoder.ring[k]
			} else if decoder.ring[k] > largest {
				largest = decoder.ring[k]
			}

		}

		if smallest*largest == *i {
			log.Printf("Part 2: %d", smallest*largest)
			break
		}
	}
}

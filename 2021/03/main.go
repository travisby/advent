package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type diagnosticReport struct {
	parsed    [][]bool
	countBits [][2]uint
}

func (d diagnosticReport) power() int64 {
	// uint so we can play with bitshifting
	var gamma uint64
	// epsilon is really the ^gamma, but for the _custom int size_
	// doing a go native ^ is just giving us the NOT of a uint64 :sweat-smile:
	// since we're only calculating 1, it's fine to just... calculate it
	var epsilon uint64

	for i, b := range d.countBits {
		if b[1] > b[0] {
			gamma += 1 << uint64(len(d.countBits)-i-1)
		} else {
			epsilon += 1 << uint64(len(d.countBits)-i-1)
		}
	}
	return int64(gamma * epsilon)
}

func (d diagnosticReport) lifeSupportRating() (*int64, error) {
	o, err := d.oxygenRating()
	if err != nil {
		return nil, err
	}
	co2, err := d.co2Rating()
	if err != nil {
		return nil, err
	}
	temp := *o * *co2
	return &temp, nil
}

func (d diagnosticReport) oxygenRating() (*int64, error) {
	return d.rating(
		func(bits [][]bool, i int) [][]bool {
			var truths, falses int
			for _, b := range bits {
				if b[i] {
					truths++
				} else {
					falses++
				}
			}
			var matcher bool
			// why couldn't I just do truths >= len(bits)/2?
			// would yield wrong results and I don't know why :o
			if truths >= falses {
				matcher = true
			}

			newBits := [][]bool{}
			for _, b := range bits {
				if b[i] == matcher {
					newBits = append(newBits, b)
				}
			}

			return newBits
		},
	)
}

func (d diagnosticReport) co2Rating() (*int64, error) {
	return d.rating(
		func(bits [][]bool, i int) [][]bool {
			var truths, falses int
			for _, b := range bits {
				if b[i] {
					truths++
				} else {
					falses++
				}
			}
			var matcher bool
			// why couldn't I just do truths >= len(bits)/2?
			// would yield wrong results and I don't know why :o
			if truths < falses {
				matcher = true
			}

			newBits := [][]bool{}
			for _, b := range bits {
				if b[i] == matcher {
					newBits = append(newBits, b)
				}
			}

			return newBits
		},
	)
}

func (d diagnosticReport) rating(bitCriteria func([][]bool, int) [][]bool) (*int64, error) {
	var f func([][]bool, int) [][]bool
	f = func(bits [][]bool, i int) [][]bool {
		if len(bits) <= 1 || i > len(bits[0]) {
			return bits
		}

		return f(bitCriteria(bits, i), i+1)
	}

	results := f(d.parsed, 0)
	if len(results) != 1 {
		return nil, fmt.Errorf("Different number of results from rating, got %+v", results)
	}

	temp := bitsToInt(results[0])
	return &temp, nil
}

func bitsToInt(bits []bool) int64 {
	var temp int64
	for i, b := range bits {
		if b {
			temp += 1 << uint64(len(bits)-i-1)
		}
	}
	return temp
}

func parseRecur(scanner *bufio.Scanner, countBits [][2]uint, parsed [][]bool) ([][2]uint, [][]bool, error) {
	if !scanner.Scan() {
		return countBits, parsed, nil
	}

	// ensure countBits is big enough
	if len(scanner.Text()) > len(countBits) {
		countBits = append(countBits, make([][2]uint, len(scanner.Text())-len(countBits))...)
	}

	parse := make([]bool, len(scanner.Text()))

	for i, c := range scanner.Text() {
		if c == '0' {
			countBits[i][0]++
		} else if c == '1' {
			countBits[i][1]++
			parse[i] = true
		} else {
			return countBits, parsed, fmt.Errorf("Unexpected N-anary digit in binary: %d", c)
		}
	}

	parsed = append(parsed, parse)

	return parseRecur(scanner, countBits, parsed)
}
func parse(r io.Reader) (*diagnosticReport, error) {
	countBits := [][2]uint{}
	scanner := bufio.NewScanner(r)

	countBits, parsed, err := parseRecur(scanner, nil, nil)
	if err != nil {
		return nil, err
	} else if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &diagnosticReport{parsed, countBits}, nil
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

	diag, err := parse(f)
	if err != nil {
		log.Fatal(err)
	} else if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Part 1: %d", diag.power())
	res, err := diag.lifeSupportRating()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Part 2: %d", *res)
}

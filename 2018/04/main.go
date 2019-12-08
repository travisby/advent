package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

func readerToSortedNewlines(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	// the output is unsorted so it's going to be just about one million times easier
	// if we store all the input, sort it, then operate on it
	rows := []string{}
	for scanner.Scan() {
		rows = append(rows, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	// simple ascii sort will be fine since
	// the first part of the line is YYYY-MM-DD
	// we'll end up chronologically
	sort.Strings(rows)

	return rows, nil
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

	rows, err := readerToSortedNewlines(f)
	if err != nil {
		log.Fatal(err)
	}

	shifts, err := readerToGuardShifts(strings.NewReader(strings.Join(rows, "\n")))
	if err != nil {
		log.Fatal(err)
	}
	for _ = range shifts {
	}
}

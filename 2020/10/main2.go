package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
)

// since we are only _SLICING_ the array
// not reordering anything
// we can assume we're talking about the same-point-in-array
// if we simply look at the length-until-end
var lookup = map[int]int{0: 0, 1: 1}

func combos(is []int) int {
	if i, ok := lookup[len(is)]; ok {
		return i
	}

	var validNextIndices []int
	for i := 1; i < len(is); i++ {
		if is[i] > is[0] && is[i]-3 <= is[0] {
			validNextIndices = append(validNextIndices, i)
		}
	}

	count := 0
	for _, v := range validNextIndices {
		count += combos(is[v:])
	}

	lookup[len(is)] = count
	return count
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

	var is []int

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		i, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		is = append(is, i)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sort.Ints(is)

	log.Printf("Part 2: %d", combos(append([]int{0}, is...)))
}

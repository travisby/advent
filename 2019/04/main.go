package main

import (
	"log"
	"os"
	"strconv"
)

type criterias struct {
	intCriteria   []func(int) bool
	strCriteria   []func(string) bool
	digitCriteria []func([]int) bool
}

func (cs criterias) valid(s string) bool {
	for _, c := range cs.strCriteria {
		if !c(s) {
			return false
		}
	}

	if len(cs.intCriteria) > 0 {
		i, err := strconv.Atoi(s)
		if err != nil {
			return false
		}

		for _, c := range cs.intCriteria {
			if !c(i) {
				return false
			}
		}
	}

	if len(cs.digitCriteria) > 0 {
		digits := make([]int, 0, len(s))
		for _, r := range s {
			if r > '9' || r < '0' {
				return false
			}
			digits = append(digits, int(r-'0'))
		}

		for _, c := range cs.digitCriteria {
			if !c(digits) {
				return false
			}
		}
	}
	return true
}

func numDigits(i int) func(string) bool {
	return func(s string) bool {
		return len(s) == i
	}
}

func withinRange(x, y int) func(int) bool {
	return func(i int) bool {
		// TODO is the inclusivity correct here?
		return i >= x && i < y
	}
}

func adjacentDigitsSameness(n []int) bool {
	for i := 0; i < len(n)-1; i++ {
		if n[i] == n[i+1] {
			return true
		}
	}
	return false
}

func adjacentTwoDigitsSameness(n []int) bool {
	// this is meant to be similar to how coreutils uniq works
	// each subsequent match of a character in a row gets put together
	// so a,b,a,a -> a,b,a
	type unique struct {
		digit       int
		occurrences int
	}

	ds := []unique{}
	for _, d := range n {
		if len(ds) == 0 || ds[len(ds)-1].digit != d {
			ds = append(ds, unique{d, 1})
		} else {
			ds[len(ds)-1].occurrences++
		}
	}
	for _, u := range ds {
		if u.occurrences == 2 {
			return true
		}
	}
	return false
}

func neverDecrease(n []int) bool {
	for i := 0; i < len(n)-1; i++ {
		if n[i] > n[i+1] {
			return false
		}
	}
	return true
}

func main() {
	var min, max int
	var err error
	if len(os.Args) != 3 {
		log.Fatalf("Expected arg1 to be the min, and arg2 to be the max")
	} else if min, err = strconv.Atoi(os.Args[1]); err != nil {
		log.Fatalf("Expected arg1 to be the int, got %v", os.Args[1])
	} else if max, err = strconv.Atoi(os.Args[2]); err != nil {
		log.Fatalf("Expected arg2 to be the int, got %v", os.Args[2])
	}

	c := criterias{
		strCriteria:   []func(string) bool{numDigits(6)},
		intCriteria:   []func(int) bool{withinRange(100000, 999999)},
		digitCriteria: []func([]int) bool{adjacentDigitsSameness, neverDecrease},
	}

	cCounter := 0

	for i := min; i < max; i++ {
		if c.valid(strconv.Itoa(i)) {
			cCounter++
		}
	}
	// PART 1
	log.Printf("Total password possibilities: %d", cCounter)

	// PART 2
	cCounter = 0
	c.digitCriteria = append(c.digitCriteria, adjacentTwoDigitsSameness)
	for i := min; i < max; i++ {
		if c.valid(strconv.Itoa(i)) {
			cCounter++
		}
	}
	log.Printf("Total password possibilities: %d", cCounter)
}

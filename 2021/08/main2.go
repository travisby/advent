package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

var permutations [][7]uint8
var segmentToNumber map[segment]uint8

func factorial(i int) int {
	result := 1
	for ; i > 0; i-- {
		result *= i
	}
	return result
}

func init() {
	// Heap's Algorithm to produce permutations
	permutations = make([][7]uint8, 5040)

	// The original order
	A := [7]uint8{0, 1, 2, 3, 4, 5, 6}
	// the "stack" pointer
	c := make([]int, len(A))

	// me tracking which point of the slice we've filled
	nextPermutation := 0

	copy(permutations[nextPermutation][:], A[:])
	nextPermutation++

	for i := 0; i < len(A); {
		if c[i] < i {
			if i%2 == 0 {
				A[0], A[i] = A[i], A[0]
			} else {
				A[c[i]], A[i] = A[i], A[c[i]]
			}
			copy(permutations[nextPermutation][:], A[:])
			nextPermutation++
			c[i] += 1
			i = 0
		} else {
			c[i] = 0
			i += 1
		}

	}

	// populate valid segments
	segmentToNumber = map[segment]uint8{
		newSegmentPanic("abcefg"):  0,
		newSegmentPanic("cf"):      1,
		newSegmentPanic("acdeg"):   2,
		newSegmentPanic("acdfg"):   3,
		newSegmentPanic("bcdf"):    4,
		newSegmentPanic("abdfg"):   5,
		newSegmentPanic("abdefg"):  6,
		newSegmentPanic("acf"):     7,
		newSegmentPanic("abcdefg"): 8,
		newSegmentPanic("abcdfg"):  9,
	}
}

type segment [7]bool

func newSegment(s string) (*segment, error) {
	var seg segment
	for _, c := range s {
		idx, ok := map[rune]uint8{
			'a': 0,
			'b': 1,
			'c': 2,
			'd': 3,
			'e': 4,
			'f': 5,
			'g': 6,
		}[c]
		if !ok {
			return nil, fmt.Errorf("Unknown letter: %q", c)
		}
		seg[idx] = true
	}
	return &seg, nil
}

func newSegmentPanic(s string) segment {
	seg, err := newSegment(s)
	if err != nil {
		panic(err)
	}
	return *seg
}

func (s segment) Valid() bool {
	_, ok := segmentToNumber[s]
	return ok
}

func (s segment) Jumble(jumble [7]uint8) segment {
	var result [7]bool

	for old, new := range jumble {
		result[new] = s[old]
	}
	return result
}

func (s segment) String() string {
	var str []string
	for i, v := range s {
		if v {
			str = append(str,
				map[int]string{
					0: "a",
					1: "b",
					2: "c",
					3: "d",
					4: "e",
					5: "f",
					6: "g",
				}[i])
		}

	}
	return strings.Join(str, "")
}

func (s segment) Score() (*uint8, error) {
	// we are reasonably sure that this will not panic
	// because we check errors when creating a segment
	res, ok := segmentToNumber[s]
	if !ok {
		return nil, fmt.Errorf("Invalid segment")
	}

	return &res, nil

}

func (s segment) Unambiguous() bool {
	var truecount uint8
	for _, b := range s {
		if b {
			truecount++
		}
	}
	// 1 4 7 8
	return truecount == 2 || truecount == 4 || truecount == 3 || truecount == 7
}

type pattern []segment

func newPattern(s string) (pattern, error) {
	// the pattern string appers like:
	// be cfbegad cbdgef fgaecd cgeb fdcge agebfd fecdb fabcd edb | fdgacbe cefdb cefbgd gcbe
	// we want to produce each word in their own  segment, ignoring "|"
	var pat pattern

	for _, p := range strings.Split(s, " ") {
		if p == "|" {
			continue
		}

		seg, err := newSegment(p)
		if err != nil {
			return nil, err
		}

		pat = append(pat, *seg)

	}
	return pat, nil
}

func (p pattern) Valid() bool {
	for _, s := range p {
		if !s.Valid() {
			return false
		}
	}
	return true
}

func (p pattern) Jumble(j [7]uint8) pattern {
	temp := make(pattern, len(p))
	for i := range p {
		temp[i] = p[i].Jumble(j)
	}
	return temp
}

func (p pattern) Score() (*int, error) {
	result := 0
	for i, j := range p {
		score, err := j.Score()
		if err != nil {
			return nil, err
		}
		result += int(*score) * int(math.Pow10(len(p)-i-1))
	}

	return &result, nil
}

func (p pattern) UnambiguousCount() int {
	var unambiguous int
	for _, s := range p {
		if s.Unambiguous() {
			unambiguous++
		}
	}
	return unambiguous
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

	var unambiguous, maxScore int
	for scanner.Scan() {
		p, err := newPattern(scanner.Text())
		if err != nil {
			log.Fatal(err)
		} else if len(p) < 4 {
			log.Fatal("Wrong sized pattern, expected at least 4 at the end")
		}

		for i := 0; !p.Valid() && i < len(permutations); i++ {
			temp := p.Jumble(permutations[i])
			if temp.Valid() {
				p = temp
			}
		}
		if !p.Valid() {
			log.Fatal("Not valid")
		}

		// we already know that we're definitely >= size 4
		p = p[len(p)-4:]

		unambiguous += p.UnambiguousCount()

		score, err := p.Score()
		if err != nil {
			log.Fatal(err)
		}
		maxScore += *score
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Part 1: %d", unambiguous)
	log.Printf("Part 2: %d", maxScore)
}

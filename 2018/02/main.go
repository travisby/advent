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

	counts := []map[int]int{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		counts = append(counts, countExactlyRepeatedLetters(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	result := checksumTwosAndThrees(counts)
	fmt.Printf("%d\n", result)
}

// we use map[int]int because there could be
// 0 letters repeated exactly 1 time
// 0 letters repeated exactly 2 times
// ....
// 1 letter repeated exactly 20000 times
// and that would be expensive ;)
func countExactlyRepeatedLetters(a string) map[int]int {
	// first, let's find out how many times a particular letter is duplicated
	lettersEncountered := make(map[rune]int)
	for _, l := range a {
		lettersEncountered[l] = lettersEncountered[l] + 1
	}

	// now, for each of those letters, categorize it into # of repeats
	results := make(map[int]int)
	for _, d := range lettersEncountered {
		results[d]++
	}
	return results
}

// counts is a list of counts of exactly duplicated letters
// (e.g., the result of countExactlyRepeatedLetters in a slice)
func checksumTwosAndThrees(counts []map[int]int) int {
	// aggregate is a count of how many rows we have that are > 1
	// not all the numbers added together
	// we don't care what the total number of 4s is, we care about
	// how many inputs had a 4
	aggregate := make(map[int]int)
	for _, count := range counts {
		for k := range count {
			aggregate[k]++
		}
	}
	return aggregate[2] * aggregate[3]
}

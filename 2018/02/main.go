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

	// p2 requires looking at the input multiple time
	// so store everything instead of scanning through just once
	inputs := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		inputs = append(inputs, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	counts := []map[int]int{}
	for _, input := range inputs {
		counts = append(counts, countExactlyRepeatedLetters(input))
	}
	result := checksumTwosAndThrees(counts)
	fmt.Printf("p1: %d\n", result)

	str1, str2 := findOnlyOneDifferent(inputs)
	resultP2 := findCommonCharacters(str1, str2)
	fmt.Printf("p2: %s\n", resultP2)
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

func numDifferent(str1, str2 string) int {
	diff := 0

	for pos := range str1 {
		// don't panic, we'll account later for the diff in lengths
		if pos >= len(str2) {
			break
		}

		if str1[pos] != str2[pos] {
			diff++
		}
	}

	// non-existant letters automatically count as diffs
	if len(str2) < len(str1) {
		diff += len(str1) - len(str2)
	} else if len(str1) < len(str2) {
		diff += len(str2) - len(str1)
	}

	return diff
}

// find first result where there's only a diff of 1 between two strings in input
func findOnlyOneDifferent(input []string) (string, string) {
	// TODO optimization: we could start b only after pos(a)
	// since those have all already been tested
	for _, a := range input {
		for _, b := range input {
			if numDifferent(a, b) == 1 {
				return a, b
			}
		}
	}
	// TODO
	// but really, error
	return "", ""
}

// given a string, return only the characters in common
func findCommonCharacters(str1, str2 string) string {
	newStr := make([]byte, 0, len(str1))
	for i := range str1 {
		if str1[i] == str2[i] {
			newStr = append(newStr, str1[i])
		}
	}
	return string(newStr)
}

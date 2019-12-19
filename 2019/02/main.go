package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"gitlab.com/travisby/advent/2019/intcodevm"
)

func main() {
	var f *os.File
	if len(os.Args) == 3 {
		var err error
		f, err = os.Open(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
	} else if len(os.Args) == 2 {
		f = os.Stdin
	} else {
		log.Printf("Missing argument!  Expected app <reverseInputToSearchFor> filename")
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	memory := []int{}

	scanner := bufio.NewScanner(f)
	scanner.Split(scanCommas)

	for scanner.Scan() {
		maybeInt := strings.Trim(scanner.Text(), "\n")

		i, err := strconv.Atoi(maybeInt)
		if err != nil {
			log.Fatal(err)
		}
		memory = append(memory, i)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	virtualMachine := intcodevm.New(len(memory))

	// PART 1
	if err := virtualMachine.Load(0, memory); err != nil {
		log.Fatal(err)
	} else if err := virtualMachine.SetNoun(12); err != nil {
	} else if err := virtualMachine.SetVerb(2); err != nil {
	} else if err := virtualMachine.Run(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Output: %d", virtualMachine.Output())
	// END PART 1

	// PART 2
	reverseInputToSearchFor, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// brute force
	for noun := 0; noun < 100; noun++ {
		for verb := 0; verb < 100; verb++ {
			if err := virtualMachine.Reset(); err != nil {
				log.Fatal(err)
			} else if err := virtualMachine.SetNoun(noun); err != nil {
				log.Fatal(err)
			} else if err := virtualMachine.SetVerb(verb); err != nil {
				log.Fatal(err)
			} else if err := virtualMachine.Run(); err != nil {
				log.Fatal(err)
			}

			if virtualMachine.Output() == reverseInputToSearchFor {
				log.Printf("For output=%d, 100*noun+verb=%d", reverseInputToSearchFor, 100*noun+verb)
				return
			}
		}
	}
	log.Fatalf("Exhaustive search yielded no result for output=%d", reverseInputToSearchFor)
}

// shamelessly stolen from ScanWords, but with "," instead of " " as the delimiter
func scanCommas(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if r != ',' {
			break
		}
	}
	// Scan until space, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if r == ',' {
			return i + width, data[start:i], nil
		}
	}
	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	// Request more data.
	return start, nil, nil
}

package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"gitlab.com/travisby/advent/2019/02/vm"
)

// shamelessly stolen from ScanWords, but with "," instead of " " as the delimiter
func ScanCommas(data []byte, atEOF bool) (advance int, token []byte, err error) {
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

	memory := []int{}

	scanner := bufio.NewScanner(f)
	scanner.Split(ScanCommas)

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

	virtualMachine := vm.New(len(memory))

	// PART 1
	// before running the program, replace position 1 with the value 12 and replace position 2 with the value 2
	memory[1] = 12
	memory[2] = 2
	// END PART 1

	if err := virtualMachine.Load(0, memory); err != nil {
		log.Fatal(err)
	} else if err := virtualMachine.Run(); err != nil {
		log.Fatal(err)
	}

	mem := virtualMachine.Memory[0]
	log.Printf("Value left at position 0: %d", mem)
}

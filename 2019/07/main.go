package main

import (
	"bufio"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"gitlab.com/travisby/advent/2019/intcodevm"
)

type amplifier struct {
	*intcodevm.VM
	phase int
}

func newAmplifier(memory []int, phase int, in io.Reader, out io.Writer) (*amplifier, error) {
	a := amplifier{
		VM:    intcodevm.New(len(memory)),
		phase: phase,
	}

	if err := a.Load(0, memory); err != nil {
		return nil, err
	}

	a.SetIn(io.MultiReader(strings.NewReader(fmt.Sprintf("%d\n", phase)), in))
	a.SetOut(out)

	return &a, nil
}

func factorial(i int) int {
	result := 1
	for ; i > 0; i-- {
		result *= i
	}
	return result
}

func toPermutations(A []int) [][]int {
	// Heap's Algorithm to produce permutations
	// The original order
	// the "stack" pointer
	c := make([]int, len(A))

	permutations := make([][]int, 0, factorial(len(A)))
	for i := 0; i < factorial(5); i++ {
		permutations = append(permutations, make([]int, len(A)))
	}

	// me tracking which point of the slice we've filled
	nextPermutation := 0

	copy(permutations[nextPermutation], A)
	nextPermutation++

	for i := 0; i < len(A); {
		if c[i] < i {
			if i%2 == 0 {
				A[0], A[i] = A[i], A[0]
			} else {
				A[c[i]], A[i] = A[i], A[c[i]]
			}
			copy(permutations[nextPermutation], A)
			nextPermutation++
			c[i] += 1
			i = 0
		} else {
			c[i] = 0
			i += 1
		}

	}
	return permutations
}

func runAmplifiersOnPhases(memory []int, phases []int, feedback bool) (*int, error) {
	var oldIn io.Reader
	stdout, lastAmplifiersWriter := io.Pipe()
	oldIn = io.Reader(strings.NewReader("0\n"))

	if feedback {
		oldIn = io.MultiReader(oldIn, stdout)
	}

	safeToRead := make(chan struct{})
	if !feedback {
		close(safeToRead)
	}

	group := new(errgroup.Group)

	for i, phase := range phases {
		newIn, newStdout := io.Pipe()

		if i == len(phases)-1 {
			newStdout = lastAmplifiersWriter
		}
		amp, err := newAmplifier(memory, phase, oldIn, newStdout)
		if err != nil {
			return nil, err
		}

		if i != 0 {
			group.Go(amp.Run)
		} else {
			group.Go(func() error {
				err := amp.Run()
				if feedback {
					close(safeToRead)
				}
				return err
			})
		}

		oldIn = newIn
	}

	<-safeToRead

	// we need to make sure we consume _all_ the input!
	res, _, err := bufio.NewReader(stdout).ReadLine()
	if err != nil {
		return nil, err
	} else if err := group.Wait(); err != nil {
		return nil, err
	}

	i, err := strconv.Atoi(string(res))
	if err != nil {
		return nil, err
	}
	return &i, nil
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
	scanner.Split(scanCommas)

	for scanner.Scan() {
		i, err := strconv.Atoi(strings.Trim(scanner.Text(), "\n"))
		if err != nil {
			log.Fatal(err)
		}
		memory = append(memory, i)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var highest int
	for _, phases := range toPermutations([]int{0, 1, 2, 3, 4}) {
		i, err := runAmplifiersOnPhases(memory, phases, false)
		if err != nil {
			log.Fatal(err)
		}
		if *i > highest {
			highest = *i
		}
	}
	log.Printf("Part 1: %d", highest)

	highest = 0
	for _, phases := range toPermutations([]int{5, 6, 7, 8, 9}) {
		i, err := runAmplifiersOnPhases(memory, phases, true)
		if err != nil {
			log.Fatal(err)
		}
		if *i > highest {
			highest = *i
		}
	}
	log.Printf("Part 2: %d", highest)
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

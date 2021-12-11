package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
)

var ErrUnimplemented = errors.New("Unimplemented")
var ErrUnknownCharacter = errors.New("Unknown character")
var ErrCorrupted = errors.New("Corrupted")
var ErrAlreadyOpened = errors.New("AlreadyOpened")
var ErrClosingLine = errors.New("Trying to close a line")

// The navigation subsystem syntax is made of several lines...
type navigationSubsystem []line

// There are one or more chunks on each line
// this makes a line just a chunk w/o an opener
type line struct {
	// type embedding to get methods for free
	chunk
}

func (l line) closersToComplete() ([]rune, error) {
	var completionRunes []rune

	for chr, err := l.nextValidCloser(); err == nil; chr, err = l.nextValidCloser() {
		if err != nil && errors.Is(err, ErrClosingLine) {
			return nil, err
		}

		completionRunes = append(completionRunes, *chr)
		l.add(*chr)
		chr, err = l.nextValidCloser()
	}

	return completionRunes, nil
}

// and chunks contain zero or more other chunks.
// (but they also contain the character that _started_ their chunk)
type chunk struct {
	opener opener
	chunks []chunk
	closed bool
}

func (c *chunk) add(r rune) error {
	// if a child is open, always focus our attention there
	if c.child().isOpen() {
		return c.child().add(r)
	}

	if isOpener(r) {
		c.chunks = append(c.chunks, chunk{})
		return c.child().open(r)
	} else {
		return c.close(r)
	}
}

func (c *chunk) isOpen() bool {
	return c != nil && !c.closed
}

func (c *chunk) child() *chunk {
	var ch *chunk
	if c != nil && len(c.chunks) > 0 {
		ch = &c.chunks[len(c.chunks)-1]
	}
	return ch
}

func (c *chunk) nextValidCloser() (*rune, error) {
	if c.child().isOpen() {
		return c.child().nextValidCloser()
	} else if c.opener == nil {
		return nil, ErrClosingLine
	}

	temp := c.opener.closer()
	return &temp, nil

}

func (c *chunk) open(r rune) error {
	if c.opener != nil {
		return fmt.Errorf("%w: already am %q was adding %q", ErrAlreadyOpened, c.opener.String(), r)
	}
	var err error
	c.opener, err = toOpener(r)
	return err
}

func (c *chunk) close(r rune) error {
	if c.opener == nil {
		return fmt.Errorf("%w: closing line!", ErrCorrupted)
	} else if c.opener.closer() != r {
		return fmt.Errorf("%w: wrong closer! Expected %c got %c", ErrCorrupted, c.opener.closer(), r)
	}

	c.closed = true
	return nil
}

type opener interface {
	closer() rune
	String() string
}

func isOpener(r rune) bool {
	_, err := toOpener(r)
	return err == nil
}

type parenthesisOpener struct{}

func (p parenthesisOpener) closer() rune   { return ')' }
func (p parenthesisOpener) String() string { return "(" }

type squareBracketOpener struct{}

func (p squareBracketOpener) closer() rune   { return ']' }
func (p squareBracketOpener) String() string { return "[" }

type braceOpener struct{}

func (p braceOpener) closer() rune   { return '}' }
func (p braceOpener) String() string { return "{" }

type angleBracketOpener struct{}

func (p angleBracketOpener) closer() rune   { return '>' }
func (p angleBracketOpener) String() string { return "<" }

func toOpener(r rune) (opener, error) {
	switch r {
	case '(':
		return parenthesisOpener{}, nil
	case '[':
		return squareBracketOpener{}, nil
	case '{':
		return braceOpener{}, nil
	case '<':
		return angleBracketOpener{}, nil
	}

	return nil, fmt.Errorf("%w: %q", ErrUnknownCharacter, r)
}

func isCloser(r rune) bool {
	return r == ']' || r == '>' || r == ')' || r == '}'
}

func closerToPart1Score(r rune) int {
	switch r {
	case ')':
		return 3
	case ']':
		return 57
	case '}':
		return 1197
	case '>':
		return 25137
	}
	return 0
}

func closersToPart2Score(rs []rune) int {
	var completionScore int
	for _, c := range rs {
		completionScore *= 5
		completionScore += map[rune]int{
			')': 1,
			']': 2,
			'}': 3,
			'>': 4,
		}[c]
	}
	return completionScore
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

	var sum int
	var completionScores []int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var l line
		var err error
		for _, c := range scanner.Text() {
			// if err, we're corrupted
			// and want to calculate part 1 things
			if err = l.add(c); err != nil {
				sum += closerToPart1Score(c)
				break
			}
		}

		// if we got all the way through w/o error
		// we're just incomplete, and want to do p2 things
		if err == nil {
			completionRunes, err := l.closersToComplete()
			if err != nil {
				log.Fatal(err)
			}

			completionScores = append(completionScores, closersToPart2Score(completionRunes))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Part 1: %d", sum)

	if len(completionScores)%2 == 0 {
		log.Fatal("Expected only odd numbers of scores")
	}

	sort.Ints(completionScores)
	log.Printf("Part 2: %d", completionScores[len(completionScores)/2])
}

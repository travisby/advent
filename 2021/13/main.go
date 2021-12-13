package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	// "sort"
)

var ErrInvalidInput = errors.New("Invalid Input")
var ErrInvalidInputPoint = fmt.Errorf("%w for point", ErrInvalidInput)
var ErrInvalidInputFold = fmt.Errorf("%w for folding instruction", ErrInvalidInput)

type Point struct {
	x, y int
}

func NewPoint(s string) (*Point, error) {
	var p Point
	if n, err := fmt.Sscanf(s, "%d,%d", &p.x, &p.y); n != 2 || err != nil {
		return nil, fmt.Errorf("%w got %q w/ err %+v", ErrInvalidInputPoint, s, err)
	}
	return &p, nil

}

// a transparentPaper is a set of points marked on the paper
type transparentPaper map[Point]struct{}

func (t transparentPaper) Fold(f FoldInstruction) transparentPaper {
	newT := make(transparentPaper)

	for p, v := range t {
		if (!f.alongXAxis && p.y == f.i) || (f.alongXAxis && p.x == f.i) {
			panic("Folds shouldn't happen where dots exist")
		}

		if !f.alongXAxis && p.y > f.i {
			p.y = 2*f.i - p.y
		} else if f.alongXAxis && p.x > f.i {
			p.x = 2*f.i - p.x
		}

		if p.x >= 0 && p.y >= 0 {
			newT[p] = v
		}
	}

	return newT
}
func (t transparentPaper) Copy() transparentPaper {
	newPaper := make(transparentPaper)
	for k, v := range t {
		newPaper[k] = v
	}
	return newPaper
}
func (t transparentPaper) NumberDotsVisible() int {
	return len(t)
}

// if we wanted to be space-efficient we could make this a Print (take in an io.Writer)
// func instead
// this ends up being space expensive because we go from a sparse map (only existing points)
// to a full slice of every point (even empty space)
func (t transparentPaper) String() string {
	var largestX, largestY int
	points := make([]Point, 0, len(t))
	for p := range t {
		points = append(points, p)

		if p.x > largestX {
			largestX = p.x
		}
		if p.y > largestY {
			largestY = p.y
		}
	}

	str := make([]byte, 0, largestX*largestY+1)
	for y := 0; y <= largestY; y++ {
		for x := 0; x <= largestX; x++ {
			if _, ok := t[Point{x, y}]; ok {
				str = append(str, '#')
			} else {
				str = append(str, '.')
			}
		}

		str = append(str, '\n')
	}

	return string(str)
}

type FoldInstruction struct {
	alongXAxis bool
	i          int
}

func NewFoldInstruction(s string) (*FoldInstruction, error) {
	var f FoldInstruction

	var axis rune

	if n, err := fmt.Sscanf(s, "fold along %c=%d", &axis, &f.i); n != 2 || err != nil {
		return nil, fmt.Errorf("%w got %q w/ err %+v", ErrInvalidInputFold, s, err)
	}

	f.alongXAxis = axis == 'x'

	if axis != 'x' && axis != 'y' {
		return nil, fmt.Errorf("%w : invalid axis, got %c", ErrInvalidInputFold, s, axis)
	}

	return &f, nil
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

	paper := make(transparentPaper)
	var paperAfterFirstFold transparentPaper
	consumingPoints := true
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// text comes into the scanner in two sections
		// first is a series of `<x>,<y>` coordinates
		// followed by a newline
		// then followed by the folding instructions
		if scanner.Text() == "" {
			consumingPoints = false
		} else if consumingPoints {
			// consuming points
			p, err := NewPoint(scanner.Text())
			if err != nil {
				log.Fatal(err)
			}
			paper[*p] = struct{}{}
		} else {
			// consuming fold instructions
			f, err := NewFoldInstruction(scanner.Text())
			if err != nil {
				log.Fatal(err)
			}
			paper = paper.Fold(*f)
			// part 1 stops after the first fold
			if paperAfterFirstFold == nil {
				// we need a copy or else it'll get modified in subsequent runs :o
				paperAfterFirstFold = paper.Copy()
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Part 1: %d", paperAfterFirstFold.NumberDotsVisible())
	log.Printf("Part 2: \n%s", paper)
}

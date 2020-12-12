package main

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"strconv"
)

type direction int

const (
	North direction = 0
	South           = 180
	East            = 90
	West            = 270
)

type plane struct {
	d *direction
	p *image.Point
}

func (p plane) manhattanDistance() int {
	return int(math.Abs(float64(p.p.X)) + math.Abs(float64(p.p.Y)))
}

func newPlane() *plane {
	p := plane{}

	var temp direction = East
	p.d = &temp
	p.p = &image.Point{}

	return &p
}

type Instruction interface {
	Apply(*plane)
}

var ErrInvalidInstruction = errors.New("Invalid instruction")

func toInstruction(s string) (Instruction, error) {
	if len(s) < 2 {
		return nil, fmt.Errorf("%w: %q", ErrInvalidInstruction, s)
	}
	i, err := strconv.Atoi(s[1:])
	if err != nil {
		return nil, err
	}
	switch s[0] {
	case 'N':
		return north(i), nil
	case 'S':
		return south(i), nil
	case 'E':
		return east(i), nil
	case 'W':
		return west(i), nil
	case 'L':
		return left(i), nil
	case 'R':
		return right(i), nil
	case 'F':
		return forward(i), nil
	}
	return nil, fmt.Errorf("%w: %q", ErrInvalidInstruction, s)
}

type north int

func (n north) Apply(p *plane) {
	*p.p = p.p.Add(image.Point{0, int(n)})
}

type south int

func (s south) Apply(p *plane) {
	*p.p = p.p.Add(image.Point{0, -int(s)})
}

type east int

func (e east) Apply(p *plane) {
	*p.p = p.p.Add(image.Point{int(e), 0})
}

type west int

func (w west) Apply(p *plane) {
	*p.p = p.p.Add(image.Point{-int(w), 0})
}

type left int

func (l left) Apply(p *plane) {
	// I think we're assuming only 90s?
	switch (int(*p.d) - int(l)) % 360 {
	case 0:
		*p.d = North
	case 90, -270:
		*p.d = East
	case 180, -180:
		*p.d = South
	case 270, -90:
		*p.d = West
	default:
		log.Fatal("Not one of the four angles we expected")
	}
}

type right int

func (r right) Apply(p *plane) {
	left(-r).Apply(p)
}

type forward int

func (f forward) Apply(p *plane) {
	switch *p.d {
	case North:
		north(f).Apply(p)
	case South:
		south(f).Apply(p)
	case East:
		east(f).Apply(p)
	case West:
		west(f).Apply(p)
	}
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

	p := newPlane()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		i, err := toInstruction(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		i.Apply(p)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("%d", p.manhattanDistance())
}

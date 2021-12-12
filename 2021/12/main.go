package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

var ErrCaveIdentifier = errors.New("Bad Cave Identifier")

type caveSystem map[string]cave

func (c caveSystem) AddPath(left, right string) error {
	var err error
	if _, ok := c[left]; !ok {
		c[left], err = NewCave(left)
		if err != nil {
			return err
		}
	}
	if _, ok := c[right]; !ok {
		c[right], err = NewCave(right)
		if err != nil {
			return err
		}
	}

	c[left].AddAdjacent(c[right])
	c[right].AddAdjacent(c[left])

	return nil
}

func NewCaveSystem() caveSystem {
	start := terminalCave{}
	start.ident = "start"
	end := terminalCave{}
	end.ident = "end"
	return caveSystem{
		"start": &start,
		"end":   &end,
	}
}

func NewCave(s string) (cave, error) {
	if s == strings.ToLower(s) {
		return &smallCave{ident: s}, nil
	} else if s == strings.ToUpper(s) {
		return &largeCave{smallCave: smallCave{ident: s}}, nil
	}
	return nil, fmt.Errorf("%w: %q", ErrCaveIdentifier, s)
}

type cave interface {
	String() string
	Adjacent() []cave
	CanVisitAnyTime() bool
	AddAdjacent(cave)
	EligibleForDoubleVisit() bool
}

type smallCave struct {
	ident    string
	adjacent []cave
}

func (s smallCave) Adjacent() []cave             { return s.adjacent }
func (s smallCave) String() string               { return s.ident }
func (s smallCave) CanVisitAnyTime() bool        { return false }
func (s *smallCave) AddAdjacent(c cave)          { s.adjacent = append(s.adjacent, c) }
func (s smallCave) EligibleForDoubleVisit() bool { return true }

// a largeCave is just a smallCave that can be visited >1
type largeCave struct {
	smallCave
}

func (l largeCave) CanVisitAnyTime() bool { return true }

// a terminalCave is just a smallCave that isn't eligible for double visits
type terminalCave struct {
	smallCave
}

func (t terminalCave) EligibleForDoubleVisit() bool { return false }

type path struct {
	caves []cave
	// this is a glorifed flag for p2
	canVisitTwice bool
}

func (p path) eligible(c cave) (nextPaths []cave) {
	canVisitTwice := p.canVisitTwice

	// build up a map of already-visited caves
	// for easy lookup
	visited := make(map[cave]bool, len(p.caves))
	for _, v := range p.caves {
		if visited[v] && !v.CanVisitAnyTime() {
			canVisitTwice = false
		}

		visited[v] = true
	}

	for _, potentialCave := range c.Adjacent() {
		// if we can visit this cave at any time, go for it
		if potentialCave.CanVisitAnyTime() ||
			// or if we haven't been there before
			!visited[potentialCave] ||
			// or finally, if we haven't used up our visit-two yet, and it's not the start/end
			(canVisitTwice && potentialCave.EligibleForDoubleVisit()) {

			// then include it as an eligible next-path!
			nextPaths = append(nextPaths, potentialCave)
		}
	}

	return nextPaths
}

// allPaths recursively finds every available path from this particular cave point until `end`
func allPaths(c cave, p path) (paths []path) {
	// first, make a new copy of this path
	// leaving room to add ourselves to it!
	newP := p
	newP.caves = make([]cave, len(p.caves), len(p.caves)+1)
	copy(newP.caves, p.caves)
	newP.caves = append(newP.caves, c)

	// terminal case
	// stop when we've hit the "end"
	if c.String() == "end" {
		// XXX: string matching is kinda meh
		// but it works
		return []path{newP}
	}

	// else, include every eligible path
	for _, newC := range newP.eligible(c) {
		paths = append(paths, allPaths(newC, newP)...)
	}

	return paths
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

	cs := NewCaveSystem()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		splits := strings.Split(scanner.Text(), "-")
		if len(splits) != 2 {
			log.Fatalf("Bad path, expected \"a-b\", got %q", splits)
		}

		cs.AddPath(splits[0], splits[1])
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Part 1: %d", len(allPaths(cs["start"], path{})))
	log.Printf("Part 2: %d", len(allPaths(cs["start"], path{canVisitTwice: true})))
}

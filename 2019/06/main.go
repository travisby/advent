package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type orbital struct {
	name       string
	orbits     *orbital
	orbittedBy []*orbital
}

func (o *orbital) addOrbittedBy(n *orbital) {
	if n.orbits == nil {
		n.orbits = o
	}

	o.orbittedBy = append(o.orbittedBy, n)
}

type orbitalMap map[string]*orbital

func (m orbitalMap) checksum() uint {
	// VERY inefficient
	// we calculate the checksum for intermediaries many times over
	var checksum uint
	for _, v := range m {
		for curPos := v; curPos.orbits != nil; curPos = curPos.orbits {
			checksum++
		}
	}

	return checksum
}

var ErrUnknownParent = errors.New("Unknown orbit-ee")

func (m orbitalMap) addOrbitByName(parent string, child string) error {
	p, ok := m[parent]
	if !ok {
		return fmt.Errorf("%w: %q", ErrUnknownParent, parent)
	}

	if _, ok := m[child]; !ok {
		m[child] = &orbital{child, p, []*orbital{}}
	}

	p.addOrbittedBy(m[child])

	return nil
}
func (m orbitalMap) getOrbital(name string) (*orbital, bool) {
	res, ok := m[name]
	return res, ok
}
func (m orbitalMap) addCOM(name string) {
	m[name] = &orbital{name, nil, []*orbital{}}
}

func newOrbitalMap() orbitalMap {
	return orbitalMap(map[string]*orbital{})
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

	om := newOrbitalMap()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		splits := strings.Split(scanner.Text(), ")")
		if len(splits) != 2 {
			log.Fatalf("Failed to parse orbital instruction")
		}

		// create a new COM and retry
		if _, ok := om.getOrbital(splits[0]); !ok {
			om.addCOM(splits[0])
		}

		if err := om.addOrbitByName(splits[0], splits[1]); err != nil {
			log.Fatal(err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Part 1: %d", om.checksum())
}

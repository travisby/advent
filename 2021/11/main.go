package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

var ErrNotImplemented = errors.New("Not implemented")

type octopuses [10][10]*octopus

func (os octopuses) fillInAdjacents() {
	for x := 0; x < len(os); x++ {
		for y := 0; y < len(os[x]); y++ {
			for _, adj := range os.adjacent(x, y) {
				os[x][y].adjacent = append(os[x][y].adjacent, adj)
			}

		}
	}
}

func (os *octopuses) stepAndCountFlashes() int {
	for x := range os {
		for y := range os[x] {
			os[x][y].step()
		}
	}

	return os.countFlashingAndReset()
}

func (os *octopuses) countFlashingAndReset() int {
	sum := 0
	for x := range os {
		for y := range os[x] {
			if os[x][y].resetFlashing() {
				sum++
			}
		}
	}
	return sum
}

func (os *octopuses) adjacent(x, y int) []*octopus {
	var adjacents []*octopus

	// left
	if x > 0 {
		adjacents = append(adjacents, os[x-1][y])
	}

	// right
	if x < len(os)-1 {
		adjacents = append(adjacents, os[x+1][y])
	}

	// up
	if y > 0 {
		adjacents = append(adjacents, os[x][y-1])
	}

	// down
	if y < len(os[x])-1 {
		adjacents = append(adjacents, os[x][y+1])
	}

	// up-left
	adjacents = append(adjacents, os[x][y])
	if y > 0 && x > 0 {
		adjacents = append(adjacents, os[x-1][y-1])
	}
	// up-right
	if y > 0 && x < len(os)-1 {
		adjacents = append(adjacents, os[x+1][y-1])
	}
	// down-left
	if y < len(os[x])-1 && x > 0 {
		adjacents = append(adjacents, os[x-1][y+1])
	}
	// down-right
	if y < len(os[x])-1 && x < len(os)-1 {
		adjacents = append(adjacents, os[x+1][y+1])
	}

	return adjacents
}

type octopus struct {
	energyLevel uint8
	adjacent    []*octopus
	flashing    bool
}

func (o *octopus) step() {
	if o.flashing {
		// no need to track energy level after we've already started flashed this round
		return
	}

	o.energyLevel++
	if o.energyLevel > 9 {
		o.flashing = true
		for _, adj := range o.adjacent {
			adj.step()
		}
	}
}

func (o *octopus) resetFlashing() bool {
	if o.flashing {
		o.flashing = false
		o.energyLevel = 0
		return true
	}
	return false
}

func newOctopus(r rune) (*octopus, error) {
	i, err := strconv.ParseUint(fmt.Sprintf("%c", r), 10, 8)
	if err != nil {
		return nil, err
	}

	return &octopus{energyLevel: uint8(i)}, nil
}

func strToOctopi(s string) ([]*octopus, error) {
	var octopi []*octopus

	for _, c := range s {
		octopus, err := newOctopus(c)
		if err != nil {
			return nil, err
		}
		octopi = append(octopi, octopus)
	}

	return octopi, nil
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

	scanner := bufio.NewScanner(f)

	var octopi octopuses
	var i int

	for ; i < 10 && scanner.Scan(); i++ {
		octs, err := strToOctopi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		} else if len(octs) != 10 {
			log.Fatal("Expected 10 octs, got %d", octs)
		} else if n := copy(octopi[i][:], octs); n != 10 {
			log.Fatal("Could not write ten into octopi")
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	} else if i != 10 {
		log.Fatalf("Input ran out too soon!  Only produced %d rows", i)
	}

	octopi.fillInAdjacents()

	var sum int
	// the index at which they simultaneously flashed
	var simultaneousFlash *int
	for i := 1; i <= 100 || simultaneousFlash == nil; i++ {
		numFlashing := octopi.stepAndCountFlashes()
		if i <= 100 {
			sum += numFlashing
		}

		if numFlashing == 100 {
			temp := i
			simultaneousFlash = &temp
		}
	}

	log.Printf("Part 1: %d", sum)
	log.Printf("Part 2: %d", *simultaneousFlash)
}

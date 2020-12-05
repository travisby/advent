package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
)

var ROW_RE = regexp.MustCompile("((?:B|F){7})((?:L|R){3})")
var ErrInvalidBoardingPass = errors.New("Invalid Boarding Pass")

type boardingPass uint16

func (b boardingPass) row() uint8 {
	return uint8(b >> 3)
}

func (b boardingPass) column() uint8 {
	return uint8(b & (1 + 2 + 4))
}

// XXX: I believe that seatID() is actually the raw uint16
// represented here?
func (b boardingPass) seatID() uint16 {
	return uint16(b.row())*8 + uint16(b.column())
}

// sort.Interface
type boardingPasses []boardingPass

func (b boardingPasses) Len() int           { return len(b) }
func (b boardingPasses) Less(i, j int) bool { return b[i].seatID() < b[j].seatID() }
func (b boardingPasses) Swap(i, j int)      { temp := b[i]; b[i] = b[j]; b[j] = temp }

func boardingPassFromString(s string) (*boardingPass, error) {
	// NOTE: We validate that the ordering is actually F/B then L/R, but the implementation below
	// does not care after validation, and will gladly use F and L or B and R interchangably
	if !ROW_RE.MatchString(s) {
		return nil, fmt.Errorf("%w: failed to match on %q", ErrInvalidBoardingPass, s)
	}

	var bp uint16
	// our boarding pass is 10 bits
	// but note that we're going from an arr where 0 is the lowest bit
	for i, j := range s {
		if j == 'B' || j == 'R' {
			bp = bp | (1 << uint16(9-i))
		}
	}

	res := boardingPass(bp)
	return &res, nil
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

	var bps []boardingPass

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		bp, err := boardingPassFromString(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		bps = append(bps, *bp)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sort.Sort(sort.Reverse(boardingPasses(bps)))

	log.Printf("Part 1: %d", bps[0].seatID())

	var bp *boardingPass
	// since we're comparing i to i+1
	// we want to loop to one-less than the end
	for i := 0; i < len(bps)-1; i++ {
		// if we find a scenario where  the next seat isn't simply - 1
		// we know -1 is our seat!
		if uint16(bps[i]-1) != uint16(bps[i+1]) {
			temp := boardingPass(uint16(bps[i]) - 1)
			bp = &temp
			break
		}
	}

	if bp == nil {
		log.Fatal("Did not find our seat!")
	}

	log.Printf("Part 2: %d", bp.seatID())

}

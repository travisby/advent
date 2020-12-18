package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var ErrUnknownInstruction = errors.New("Unknown instruction")
var MASK_REGEX = regexp.MustCompile(`mask = ([0-9X]+)`)
var MEMSET_REGEX = regexp.MustCompile(`mem\[(\d+)\] = (\d+)`)

type SoftwareDecoder struct {
	mask   mask
	memory map[address]uint64
}

func newSoftwareDecoder() *SoftwareDecoder {
	return &SoftwareDecoder{memory: make(map[address]uint64)}
}

type Instruction interface {
	apply(*SoftwareDecoder)
}

func ToInstruction(s string) (Instruction, error) {
	if match := MASK_REGEX.FindStringSubmatch(s); len(match) == 2 {
		return newMask(match[1])
	} else if match := MEMSET_REGEX.FindStringSubmatch(s); len(match) == 3 {
		m, err := strconv.ParseUint(match[1], 10, 36)
		if err != nil {
			return nil, err
		}

		v, err := strconv.ParseUint(match[2], 10, 36)
		if err != nil {
			return nil, err
		}

		return memset{addr: address(m), val: v}, nil

	}
	return nil, fmt.Errorf("%w, %q", ErrUnknownInstruction, s)
}

// addresses are 36 bits, not 64
// but this is going to be the easiest way to represent it
type address uint64

type mask struct {
	// a mask consists of two addresses
	// one is the mask-values, where Xs are 0s, 0s are 0s, and 1s are 1s
	// the other is 1 where there was either a 0 or a 1, and 0 where there was an X
	// to apply a mask, you will use `b` to determine if a mask is actually set on that bit
	// if Bi is 1, then you will use Ai.  If Bi is 0 then you will not use anything in A
	a address
	b address
}

func newMask(s string) (*mask, error) {
	m := mask{}

	var temp uint64
	if n, err := fmt.Sscanf(strings.ReplaceAll(s, "X", "0"), "%b", &temp); err != nil {
		return nil, err
	} else if n != 1 {
		return nil, fmt.Errorf("Expected 1 items in %q got %d", strings.ReplaceAll(s, "X", "0"), n)
	}
	m.a = address(temp)

	if n, err := fmt.Sscanf(strings.ReplaceAll(strings.ReplaceAll(s, "0", "1"), "X", "0"), "%b", &temp); err != nil {
		return nil, err
	} else if n != 1 {
		return nil, fmt.Errorf("Expected 1 items in %q got %d", s, n)
	}
	m.b = address(temp)

	return &m, nil
}

func (m mask) apply(s *SoftwareDecoder) {
	s.mask = m
}

// take the value v and apply m mask to it
func (m mask) mask(v uint64) uint64 {
	for i := uint(0); i < 36; i++ {
		if m.b>>i&1 == 1 {
			// we want to SET v[i] = 1
			if m.a>>i&1 == 1 {
				v |= 1 << i
			} else {
				// we want to SET v[i] = 0
				// thanks to https://stackoverflow.com/questions/23192262/how-would-you-set-and-clear-a-single-bit-in-go for the idea on how to clear
				v &= ^(1 << i)
			}
		}
	}
	return v
}

type memset struct {
	addr address
	val  uint64
}

func (m memset) apply(s *SoftwareDecoder) {
	s.memory[m.addr] = s.mask.mask(m.val)
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

	sd := newSoftwareDecoder()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		inst, err := ToInstruction(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		inst.apply(sd)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var sum uint64
	for _, v := range sd.memory {
		sum += v
	}
	log.Printf("Part 1: %d", sum)

}

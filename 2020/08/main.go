package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
)

type vm struct {
	pc  int
	acc int
	// not a von-neuman machine
	// our memory is separate from our program instruction
	inst []instruction

	// keys are isntructions we've already executed
	visited map[int]struct{}
}

func (v *vm) reset() {
	v.pc = 0
	v.acc = 0
	v.visited = map[int]struct{}{}
}
func newVM(inst []instruction) *vm {
	return &vm{inst: inst, visited: map[int]struct{}{}}
}
func (v *vm) runOne() {
	v.visited[v.pc] = struct{}{}

	v.inst[v.pc].apply(v)
}
func (v *vm) inInfiniteLoop() bool {
	_, ok := v.visited[v.pc]
	return ok
}

// XXX: No protection if there is no infinite loop
// we'll actually crash on slice access
func (v *vm) runUntilInfiniteLoop() {
	for !v.inInfiniteLoop() {
		v.runOne()
	}
}

type instruction interface {
	apply(*vm)
}
type nop struct{}

func (n nop) apply(v *vm) {
	v.pc++
}

type acc struct {
	arg int
}

func (a acc) apply(v *vm) {
	v.acc += a.arg
	v.pc++
}

type jmp struct {
	arg int
}

func (j jmp) apply(v *vm) {
	v.pc += j.arg
}

var INST = regexp.MustCompile("(nop|acc|jmp) ([\\+-]\\d+)")
var ErrUnknownInstruction = errors.New("Unknown instruction")

func parseInst(s string) (instruction, error) {
	matches := INST.FindStringSubmatch(s)
	if len(matches) != 3 {
		return nil, fmt.Errorf("%w: %q", ErrUnknownInstruction, s)
	}

	var arg int
	if _, err := fmt.Sscan(matches[2], &arg); err != nil {
		return nil, fmt.Errorf("%w: parsing %q", err, matches[2])
	}

	var inst instruction
	switch matches[1] {
	case "jmp":
		inst = jmp{arg}
	case "acc":
		inst = acc{arg}
	case "nop":
		inst = nop{}
	}

	return inst, nil
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

	var insts []instruction

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		inst, err := parseInst(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		insts = append(insts, inst)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	v := newVM(insts)
	v.runUntilInfiniteLoop()

	log.Printf("Part 1: %d", v.acc)
}

package program

import (
	"errors"
	"fmt"
)

// HALT is the error returned when we attempt to Apply a halt instruction.  This is only returned when it's a graceful halt, otherwise there is ErrUnexpectedHalt
// stupid syntax here is to avoid golint "ErrFoo" comment
var HALT = func() error { return errors.New("HALT") }()

// ErrUnexpectedHalt is when we have an unexpected error that leads to a halt
var ErrUnexpectedHalt = errors.New("Unexpected err: HALT")

// Instruction is an instruction in intcode
type Instruction interface {
	// Apply performs the Instruction instruction on the provided piece of memory
	Apply(memory []int) error
	// NumParametrs returns the number of ints that made up the instruction
	numParameters() int
	String() string
}

type opcode int

const addOp opcode = 1
const multiplyOp opcode = 2
const haltOp opcode = 99

type add struct {
	arg1pos int
	arg2pos int
	destpos int
}

func (a add) Apply(memory []int) error {
	if a.arg1pos >= len(memory) || a.arg2pos >= len(memory) || a.destpos >= len(memory) {
		return ErrUnexpectedHalt
	}

	memory[a.destpos] = memory[a.arg1pos] + memory[a.arg2pos]

	return nil
}
func (a add) numParameters() int {
	return 3
}
func (a add) String() string {
	return fmt.Sprintf("Add{$(%d) $(%d)} -> $(%d)", a.arg1pos, a.arg2pos, a.destpos)
}

type multiply struct {
	arg1pos int
	arg2pos int
	destpos int
}

func (m multiply) Apply(memory []int) error {
	if m.arg1pos >= len(memory) || m.arg2pos >= len(memory) || m.destpos >= len(memory) {
		return ErrUnexpectedHalt
	}

	memory[m.destpos] = memory[m.arg1pos] * memory[m.arg2pos]

	return nil
}
func (m multiply) numParameters() int {
	return 3
}
func (m multiply) String() string {
	return fmt.Sprintf("Multiply{$(%d) $(%d)} -> $(%d)", m.arg1pos, m.arg2pos, m.destpos)
}

type halt struct {
}

func (h halt) Apply(memory []int) error {
	return HALT
}
func (h halt) numParameters() int {
	return 0
}
func (h halt) String() string {
	return "Halt"
}

func newInstruction(memory []int) (Instruction, error) {
	switch opcode(memory[0]) {
	case addOp:
		return add{memory[1], memory[2], memory[3]}, nil
	case multiplyOp:
		return multiply{memory[1], memory[2], memory[3]}, nil
	case haltOp:
		// by returning a HALT here we can get things to stop w/o running Apply()
		return halt{}, HALT
	}

	return nil, fmt.Errorf("Unknown opcode: %d", memory[0])
}

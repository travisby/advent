package program

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// HALT is the error returned when we attempt to Apply a halt instruction.  This is only returned when it's a graceful halt, otherwise there is ErrUnexpectedHalt
// stupid syntax here is to avoid golint "ErrFoo" comment
var HALT = func() error { return errors.New("HALT") }()

// ErrUnexpectedHalt is when we have an unexpected error that leads to a halt
var ErrUnexpectedHalt = errors.New("Unexpected err: HALT")

// ErrNoInput is propagated up when an input instruction was reached but there is no input to read from
var ErrNoInput = errors.New("No input readable")

// ErrInvalidInput is propagated up when an input instruction was reached and data successfully read from, but it could not be read as an integer
var ErrInvalidInput = errors.New("Invalid input")

// ErrOutput is propagated up when we cannot write after encountering a write instruction
var ErrOutput = errors.New("Unable to output")

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
const inputOp opcode = 3
const outputOp opcode = 4

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

type input struct {
	parameter1 int
	input      io.Reader
}

func (i input) Apply(memory []int) error {
	var inputer = io.Reader(os.Stdin)
	if i.input != nil {
		inputer = i.input
	}

	if _, err := fmt.Fscanf(inputer, "%d\n", &memory[i.parameter1]); err != nil {
		return ErrInvalidInput
	}

	return nil
}

func (i input) numParameters() int {
	return 1
}
func (i input) String() string {
	return fmt.Sprintf("input{%d}", i.parameter1)
}

type output struct {
	parameter1 int
	writer     io.Writer
}

func (o output) Apply(memory []int) error {
	// XXX: maybe shouldn't directly write to Stderr
	fmt.Fprintf(os.Stderr, "Please enter input: ")

	if _, err := fmt.Fprintf(o.writer, "%d\n", memory[o.parameter1]); err != nil {
		return ErrOutput
	}
	return nil
}

func (o output) numParameters() int {
	return 1
}
func (o output) String() string {
	return fmt.Sprintf("output{%d}", o.parameter1)
}

func newInstruction(memory []int, in io.Reader, out io.Writer) (Instruction, error) {
	switch opcode(memory[0]) {
	case addOp:
		return add{memory[1], memory[2], memory[3]}, nil
	case multiplyOp:
		return multiply{memory[1], memory[2], memory[3]}, nil
	case haltOp:
		// by returning a HALT here we can get things to stop w/o running Apply()
		return halt{}, HALT
	case inputOp:
		return input{memory[1], in}, nil
	case outputOp:
		return output{memory[1], out}, nil
	}

	return nil, fmt.Errorf("Unknown opcode: %d", memory[0])
}

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

type add struct {
	parameter1 parameter
	parameter2 parameter
	dest       position
}

func (a add) Apply(memory []int) error {
	p1, err := a.parameter1.Get(memory)
	if err != nil {
		return err
	}
	p2, err := a.parameter2.Get(memory)
	if err != nil {
		return err
	}

	return a.dest.Set(*p1+*p2, memory)
}
func (a add) numParameters() int {
	return 3
}
func (a add) String() string {
	return fmt.Sprintf("Add{%s, %s} -> %s", a.parameter1, a.parameter2, a.dest)
}

type multiply struct {
	parameter1 parameter
	parameter2 parameter
	dest       position
}

func (m multiply) Apply(memory []int) error {
	p1, err := m.parameter1.Get(memory)
	if err != nil {
		return err
	}
	p2, err := m.parameter2.Get(memory)
	if err != nil {
		return err
	}

	return m.dest.Set(*p1**p2, memory)
}
func (m multiply) numParameters() int {
	return 3
}
func (m multiply) String() string {
	return fmt.Sprintf("Multiply{%s, %s} -> %s", m.parameter1, m.parameter2, m.dest)
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
	parameter1 position
	input      io.Reader
}

func (i input) Apply(memory []int) error {
	var inputer = io.Reader(os.Stdin)
	if i.input != nil {
		inputer = i.input
	}

	// XXX: maybe shouldn't directly write to Stderr
	fmt.Fprintf(os.Stderr, "Please enter input: ")

	var temp int
	if _, err := fmt.Fscanf(inputer, "%d\n", &temp); err != nil {
		return ErrInvalidInput
	}

	return i.parameter1.Set(temp, memory)
}

func (i input) numParameters() int {
	return 1
}
func (i input) String() string {
	return fmt.Sprintf("input{%d}", i.parameter1)
}

type output struct {
	parameter1 parameter
	writer     io.Writer
}

func (o output) Apply(memory []int) error {
	i, err := o.parameter1.Get(memory)
	if err != nil {
		return ErrOutput
	} else if _, err := fmt.Fprintf(o.writer, "%d\n", *i); err != nil {
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
	// instructions are of form ABCDE
	// DE - two-digit opcode
	// C - mode of 1st parameter
	// B - mode of 2nd parameter
	// A - mode of 3rd parameter
	// we assume leading zeros up until the correct number of arguments

	// XXX: this assumes everything but halt is a one-digit opcode preceeded by parameter modes
	// if we ever add in an opcode like "42" we'll have to refactor

	// % 10 will give us just the last digit now that we've gotten all of our two-digit opcodes out of the way
	// for each parameter we're going to / 10, /100, etc. to get the parameter mode
	switch opcode(memory[0] % 10) {
	// TODO we can make all of these type assertions safe, but then the code will look a lot rougher
	// maybe when we refactor parametermodes and opcodes into their own files...
	case addOp:
		return add{
			parameterMode(memory[1], digitAt(memory[0], 100)),
			parameterMode(memory[2], digitAt(memory[0], 1000)),
			parameterMode(memory[3], digitAt(memory[0], 10000)).(position),
		}, nil
	case multiplyOp:
		return multiply{
			parameterMode(memory[1], digitAt(memory[0], 100)),
			parameterMode(memory[2], digitAt(memory[0], 1000)),
			parameterMode(memory[3], digitAt(memory[0], 10000)).(position),
		}, nil
	case haltOp % 10:
		// now do the full comparison since we care about both digits
		if opcode(memory[0]) == haltOp {
			return halt{}, HALT
		}
	case inputOp:
		return input{
			parameterMode(memory[1], digitAt(memory[0], 100)).(position),
			in,
		}, nil
	case outputOp:
		return output{
			parameterMode(memory[1], digitAt(memory[0], 100)).(position),
			out,
		}, nil
	}

	return nil, fmt.Errorf("Unknown opcode: %d", memory[0])
}

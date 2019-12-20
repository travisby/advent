package program

import (
	"errors"
	"io"
)

// ErrUnknownOpcode is when we don't support the presented two-digit opcode
var ErrUnknownOpcode = errors.New("Unexpected opcode")

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
			parameterMode(memory[1], digitAt(memory[0], 100)),
			out,
		}, nil
	}

	return nil, ErrUnknownOpcode
}

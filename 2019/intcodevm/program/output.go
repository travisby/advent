package program

import (
	"errors"
	"fmt"
	"io"
)

// ErrOutput is propagated up when we cannot write after encountering a write instruction
var ErrOutput = errors.New("Unable to output")

// Instruction is an instruction in intcode
type Instruction interface {
	// Apply performs the Instruction instruction on the provided piece of memory
	Apply(memory []int) error
	// NumParametrs returns the number of ints that made up the instruction
	numAdvanceIP() int
	String() string
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

func (o output) numAdvanceIP() int {
	return 2
}
func (o output) String() string {
	return fmt.Sprintf("output{%d}", o.parameter1)
}

package program

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// ErrNoInput is propagated up when an input instruction was reached but there is no input to read from
var ErrNoInput = errors.New("No input readable")

// ErrInvalidInput is propagated up when an input instruction was reached and data successfully read from, but it could not be read as an integer
var ErrInvalidInput = errors.New("Invalid input")

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

func (i input) numAdvanceIP() int {
	return 2
}
func (i input) String() string {
	return fmt.Sprintf("input{%d}", i.parameter1)
}

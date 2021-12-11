package program

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"
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
	// let's remind the program runner we want input, if we haven't gotten any after some delay
	inputGiven := make(chan struct{})
	go func() {
		select {
		case <-inputGiven:
		case <-time.After(500 * time.Millisecond):
			fmt.Fprintf(os.Stderr, "Please enter input: ")
		}
	}()

	var temp int
	if _, err := fmt.Fscan(i.input, &temp); err != nil {
		return ErrInvalidInput
	}
	// we finally got our input, signal that we might not need to ask for input
	close(inputGiven)

	return i.parameter1.Set(temp, memory)
}

func (i input) numAdvanceIP() int {
	return 2
}
func (i input) String() string {
	return fmt.Sprintf("input{%d}", i.parameter1)
}

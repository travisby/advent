package program

import "errors"

// HALT is the error returned when we attempt to Apply a halt instruction.  This is only returned when it's a graceful halt, otherwise there is ErrUnexpectedHalt
// stupid syntax here is to avoid golint "ErrFoo" comment
var HALT = func() error { return errors.New("HALT") }()

type halt struct {
}

func (h halt) Apply(memory []int) error {
	return HALT
}

func (h halt) numAdvanceIP() int {
	return 0
}

func (h halt) String() string {
	return "Halt"
}

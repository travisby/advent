package program

import (
	"fmt"
)

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
func (a add) numAdvanceIP() int {
	return 4
}
func (a add) String() string {
	return fmt.Sprintf("Add{%s, %s} -> %s", a.parameter1, a.parameter2, a.dest)
}

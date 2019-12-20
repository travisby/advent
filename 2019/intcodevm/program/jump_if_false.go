package program

import (
	"fmt"
)

type jumpFalse struct {
	parameter parameter
	goTo      position
	pc        *int
}

func (j jumpFalse) Apply(memory []int) error {
	if i, err := j.parameter.Get(memory); err != nil {
		return err
	} else if *i == 0 {
		*j.pc = j.goTo.address
	}

	return nil
}
func (j jumpFalse) numAdvanceIP() int {
	return 3
}
func (j jumpFalse) String() string {
	return fmt.Sprintf("JumpIfFalse{%s} -> %s", j.parameter, j.goTo)
}

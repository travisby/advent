package program

import (
	"fmt"
)

type jumpTrue struct {
	parameter parameter
	goTo      position
	pc        *int
}

func (j jumpTrue) Apply(memory []int) error {
	if i, err := j.parameter.Get(memory); err != nil {
		return err
	} else if *i != 0 {
		*j.pc = j.goTo.address
	}

	return nil
}
func (j jumpTrue) numAdvanceIP() int {
	return 3
}
func (j jumpTrue) String() string {
	return fmt.Sprintf("JumpIfTrue{%s} -> %s", j.parameter, j.goTo)
}

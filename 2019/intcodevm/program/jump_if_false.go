package program

import (
	"fmt"
)

type jumpFalse struct {
	parameter parameter
	goTo      parameter
	ip        *int
}

func (j jumpFalse) Apply(memory []int) error {
	i, err := j.parameter.Get(memory)
	if err != nil {
		return err
	}

	ip, err := j.goTo.Get(memory)
	if err != nil {
		return err
	}

	if *i == 0 {
		*j.ip = *ip
	}

	return nil
}
func (j jumpFalse) numAdvanceIP() int {
	return 3
}
func (j jumpFalse) String() string {
	return fmt.Sprintf("JumpIfFalse{%s} -> %s", j.parameter, j.goTo)
}

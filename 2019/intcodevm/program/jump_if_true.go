package program

import (
	"fmt"
)

type jumpTrue struct {
	parameter parameter
	goTo      parameter
	ip        *int
}

func (j jumpTrue) Apply(memory []int) error {
	i, err := j.parameter.Get(memory)
	if err != nil {
		return err
	}

	ip, err := j.goTo.Get(memory)
	if err != nil {
		return err
	}

	if *i != 0 {
		*j.ip = *ip
	}

	return nil
}
func (j jumpTrue) numAdvanceIP() int {
	return 3
}
func (j jumpTrue) String() string {
	return fmt.Sprintf("JumpIfTrue{%s} -> %s", j.parameter, j.goTo)
}

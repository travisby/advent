package program

import "fmt"

type multiply struct {
	parameter1 parameter
	parameter2 parameter
	dest       position
}

func (m multiply) Apply(memory []int) error {
	p1, err := m.parameter1.Get(memory)
	if err != nil {
		return err
	}
	p2, err := m.parameter2.Get(memory)
	if err != nil {
		return err
	}

	return m.dest.Set(*p1**p2, memory)
}
func (m multiply) numAdvanceIP() int {
	return 4
}
func (m multiply) String() string {
	return fmt.Sprintf("Multiply{%s, %s} -> %s", m.parameter1, m.parameter2, m.dest)
}

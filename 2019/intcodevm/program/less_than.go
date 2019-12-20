package program

import "fmt"

type lessThan struct {
	parameter1 parameter
	parameter2 parameter
	dest       position
}

func (l lessThan) Apply(memory []int) error {
	p1, err := l.parameter1.Get(memory)
	if err != nil {
		return err
	}
	p2, err := l.parameter2.Get(memory)
	if err != nil {
		return err
	}

	result := 0
	if *p1 < *p2 {
		result = 1
	}

	return l.dest.Set(result, memory)
}
func (l lessThan) numAdvanceIP() int {
	return 4
}
func (l lessThan) String() string {
	return fmt.Sprintf("LessThan{%s, %s} -> %s", l.parameter1, l.parameter2, l.dest)
}

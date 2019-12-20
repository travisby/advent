package program

import "fmt"

type equals struct {
	parameter1 parameter
	parameter2 parameter
	dest       position
}

func (e equals) Apply(memory []int) error {
	p1, err := e.parameter1.Get(memory)
	if err != nil {
		return err
	}
	p2, err := e.parameter2.Get(memory)
	if err != nil {
		return err
	}

	result := 0
	if *p1 == *p2 {
		result = 1
	}

	return e.dest.Set(result, memory)
}
func (e equals) numAdvanceIP() int {
	return 4
}
func (e equals) String() string {
	return fmt.Sprintf("Equal{%s, %s} -> %s", e.parameter1, e.parameter2, e.dest)
}

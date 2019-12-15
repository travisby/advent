package program

import "fmt"

// Scanner is a program with its iterator
type Scanner interface {
	Err() error
	Scan() bool
	Intcode() Intcode
}

type scanner struct {
	error
	memory []int
	pc     int
	token  *Intcode
}

// NewScanner creates a new Program scanner frmo a memory block
func NewScanner(memory []int) (Scanner, error) {
	return &scanner{memory: memory}, fmt.Errorf("Not implemented")
}

func (s *scanner) Err() error {
	return s.error
}

func (s *scanner) Scan() bool {
	s.pc += 4 // each intcode is 4 bytes
	s.token, s.error = newIntcode(s.memory[s.pc], s.memory[s.pc+1], s.memory[s.pc+2], s.memory[s.pc+3])

	// TODO I think this is backwards
	return s.error != nil
}

func (s *scanner) Intcode() Intcode {
	if s.Err() != nil {
		// TODO
		return Intcode([4]int{0, 0, 0, 0})
	}
	// TODO dangerous
	return *s.token
}

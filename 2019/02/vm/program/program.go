package program

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
	token  Intcode
}

// NewScanner creates a new Program scanner from a memory block
func NewScanner(memory []int) Scanner {
	return &scanner{memory: memory}
}

// Err returns the first non-HALT error that was encountered by the Scanner.
func (s *scanner) Err() error {
	if s.error == HALT {
		return nil
	}

	return s.error
}

func (s *scanner) Scan() bool {
	if s.error != nil {
		return false
	}

	s.token, s.error = newIntcode(s.memory[s.pc], s.memory[s.pc+1], s.memory[s.pc+2], s.memory[s.pc+3])

	// advance the program counter
	s.pc += 4 // each intcode is 4 bytes

	// TODO I think this is backwards
	return s.error == nil
}

func (s *scanner) Intcode() Intcode {
	if s.Err() != nil {
		return nil
	}

	return s.token
}

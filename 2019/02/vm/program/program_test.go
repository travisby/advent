package program

import (
	"errors"
	"testing"
)

/*
type scanner struct {
	error
	memory []int
	pc     int
	token  Intcode
}
*/
func TestProgramHaltsIfErrorBeforeScan(t *testing.T) {
	s := scanner{error: errors.New(""), pc: 5}
	if s.Scan() {
		t.Error("Scan() should return false if there's an error condition")
	}
	if s.pc != 5 {
		t.Errorf("Expected pc to stay at 5, got: (%d)", s.pc)
	}
}

func TestNewScanner(t *testing.T) {
	s := NewScanner([]int{})
	if err := s.Err(); err != nil {
		t.Errorf("NewScanner should not set an error, got: (%+v)", err)
	}

	if intcode := s.Intcode(); intcode != nil {
		t.Errorf("NewScanner should not set an intcode, got: (%+v)", intcode)
	}

}

func TestSimpleHaltProgram(t *testing.T) {
	s := scanner{memory: []int{int(haltOp), 0, 0, 0}}

	if s.Scan() {
		t.Error("Scan() of an error should return false")
	}

	if s.pc != 4 {
		t.Errorf("Expected pc to advance to 4, got: (%d)", s.pc)
	}

	if s.error != HALT {
		t.Errorf("Expected internal err to be set to HALT, got (%+v)", s.error)
	}

	if _, ok := s.token.(halt); !ok {
		t.Errorf("Expected token to be halt, got (%+v)", s.token)
	}

	if err := s.Err(); err != nil {
		t.Errorf("Expected Err() to return nil for HALT, got (%+v)", err)
	}
}

func TestProgramGeneratesCorrectOpcodes(t *testing.T) {
	testCases := []struct {
		title            string
		memory           []int
		expectedIntcodes []Intcode
		expectedError    error
	}{
		{
			"Simple use of everything",
			[]int{1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50},
			[]Intcode{add{9, 10, 3}, multiply{3, 11, 0}, halt{}},
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			actualIntcodes := make([]Intcode, 0, len(tc.expectedIntcodes))

			p := scanner{memory: tc.memory}
			for p.Scan() {
				actualIntcodes = append(actualIntcodes, p.Intcode())
			}
			// add trailing halt
			if p.Err() == nil {
				actualIntcodes = append(actualIntcodes, p.Intcode())
			}

			if !intcodesEqual(tc.expectedIntcodes, actualIntcodes) {
				t.Errorf("Expected intcodes (%+v), got (%+v)", tc.expectedIntcodes, actualIntcodes)
			}

			if p.Err() != nil && tc.expectedError != nil && p.Err().Error() != tc.expectedError.Error() {
				t.Errorf("Expect error (%+v), got (%+v)", p.Err(), tc.expectedError)
			}
		})
	}
}

func intcodesEqual(a []Intcode, b []Intcode) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].String() != b[i].String() {
			return false
		}
	}
	return true
}

package program

import (
	"errors"
	"testing"
)

func TestProgramHaltsIfErrorBeforeScan(t *testing.T) {
	s := scanner{error: errors.New(""), instructionPointer: 5}
	if s.Scan() {
		t.Error("Scan() should return false if there's an error condition")
	}
	if s.instructionPointer != 5 {
		t.Errorf("Expected instruction poitner to stay at 5, got: (%d)", s.instructionPointer)
	}
}

func TestNewScanner(t *testing.T) {
	s := NewScanner([]int{}, nil, nil)
	if err := s.Err(); err != nil {
		t.Errorf("NewScanner should not set an error, got: (%+v)", err)
	}

	if intcode := s.Instruction(); intcode != nil {
		t.Errorf("NewScanner should not set an intcode, got: (%+v)", intcode)
	}

}

func TestSimpleHaltProgram(t *testing.T) {
	s := scanner{memory: []int{int(haltOp), 0, 0, 0}}

	if s.Scan() {
		t.Error("Scan() of an error should return false")
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
		title                string
		memory               []int
		expectedInstructions []Instruction
		expectedError        error
	}{
		{
			"Simple use of everything",
			[]int{1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50},
			[]Instruction{add{position{9}, position{10}, position{3}}, multiply{position{3}, position{11}, position{0}}, halt{}},
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			actualInstructions := make([]Instruction, 0, len(tc.expectedInstructions))

			p := scanner{memory: tc.memory}
			for p.Scan() {
				actualInstructions = append(actualInstructions, p.Instruction())
			}
			// add trailing halt
			if p.Err() == nil {
				actualInstructions = append(actualInstructions, p.Instruction())
			}

			if !intcodesEqual(tc.expectedInstructions, actualInstructions) {
				t.Errorf("Expected intcodes (%+v), got (%+v)", tc.expectedInstructions, actualInstructions)
			}

			if p.Err() != nil && tc.expectedError != nil && p.Err().Error() != tc.expectedError.Error() {
				t.Errorf("Expect error (%+v), got (%+v)", p.Err(), tc.expectedError)
			}
		})
	}
}

func intcodesEqual(a []Instruction, b []Instruction) bool {
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

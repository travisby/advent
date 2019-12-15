package program

import (
	"errors"
	"testing"
)

func TestIntcodeParse(t *testing.T) {
	testCases := []struct {
		title           string
		args            [4]int
		expectedIntcode Intcode
		expectedErr     error
	}{
		{"Simple add", [4]int{1, 10, 20, 30}, add{10, 20, 30}, nil},
		{"Simple multiply", [4]int{2, 10, 20, 30}, multiply{10, 20, 30}, nil},
		{"Simple halt", [4]int{99, -1, 0, 8}, halt{}, HALT},
		{"Simple error", [4]int{98, 0, 0, 0}, nil, errors.New("Unknown opcode: 98")},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			intcode, err := newIntcode(tc.args[0], tc.args[1], tc.args[2], tc.args[3])
			if err == nil && err != tc.expectedErr {
				t.Errorf("Got err (%+v) expected (%+v)", err, tc.expectedErr)
			} else if err != nil && tc.expectedErr == nil {
				t.Errorf("Got err (%+v) expected (%+v)", err, tc.expectedErr)
			} else if err != nil && err.Error() != tc.expectedErr.Error() {
				t.Errorf("Got err (%+v) expected (%+v)", err, tc.expectedErr)
			} else if intcode != tc.expectedIntcode {
				t.Errorf("Got intcode (%s) expected (%s)", intcode, tc.expectedIntcode)
			}
		})
	}
}

func TestApply(t *testing.T) {
	testCases := []struct {
		title               string
		intcode             Intcode
		memory              []int
		expectedMemoryAfter []int
		expectedErr         error
	}{
		{
			"Simple Add",
			add{10, 20, 30},
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3},
			nil,
		},
		{
			"Simple Mult",
			multiply{10, 20, 30},
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
			nil,
		},
		{
			"Simple Halt",
			halt{},
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			HALT,
		},
		{
			"Error Halt",
			add{5, 1, 2},
			[]int{0, 0, 0, 0, 0},
			[]int{0, 0, 0, 0, 0},
			ErrUnexpectedHalt,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			if err := tc.intcode.Apply(tc.memory); err != nil {
				if tc.expectedErr == nil {
					t.Fatalf("Got err (%+v) expected (%+v)", err, tc.expectedErr)
				} else if err.Error() != tc.expectedErr.Error() {
					t.Fatalf("Got err (%+v) expected (%+v)", err, tc.expectedErr)
				}
			} else if !memEquals(tc.expectedMemoryAfter, tc.memory) {
				t.Errorf("Got memory: (%+v) expected (%+v)", tc.memory, tc.expectedMemoryAfter)
			}
		})
	}
}

func memEquals(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

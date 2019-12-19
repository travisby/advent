package program

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestInstructionParse(t *testing.T) {
	testCases := []struct {
		title               string
		args                []int
		expectedInstruction Instruction
		expectedErr         error
	}{
		{"Simple add", []int{1, 10, 20, 30}, add{10, 20, 30}, nil},
		{"Simple multiply", []int{2, 10, 20, 30}, multiply{10, 20, 30}, nil},
		{"Simple halt", []int{99, -1, 0, 8}, halt{}, HALT},
		{"Simple error", []int{98, 0, 0, 0}, nil, errors.New("Unknown opcode: 98")},
		{"Simple input", []int{3, 50}, input{parameter1: 50}, nil},
		{"Simple output", []int{4, 50}, output{parameter1: 50}, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			intcode, err := newInstruction(tc.args, nil, nil)
			if err == nil && err != tc.expectedErr {
				t.Errorf("Got err (%+v) expected (%+v)", err, tc.expectedErr)
			} else if err != nil && tc.expectedErr == nil {
				t.Errorf("Got err (%+v) expected (%+v)", err, tc.expectedErr)
			} else if err != nil && err.Error() != tc.expectedErr.Error() {
				t.Errorf("Got err (%+v) expected (%+v)", err, tc.expectedErr)
			} else if intcode != tc.expectedInstruction {
				t.Errorf("Got intcode (%s) expected (%s)", intcode, tc.expectedInstruction)
			}
		})
	}
}

func TestApply(t *testing.T) {
	testCases := []struct {
		title               string
		intcode             Instruction
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
		{
			"Simple Input",
			input{3, strings.NewReader("-45")},
			[]int{0, 0, 0, 0},
			[]int{0, 0, 0, -45},
			nil,
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

	// TODO this is an awful way to test simple output, we should be testing actually with the test cases
	t.Run("Simple output", func(t *testing.T) {
		buffer := bytes.NewBuffer(nil)
		memory := []int{0, 0, 9, 0}
		expectedMemory := []int{0, 0, 9, 0}
		expectedOutput := "9\n"
		if err := (output{2, buffer}.Apply(memory)); err != nil {
			t.Fatal(err)
		}

		actualOutput := buffer.String()
		if expectedOutput != actualOutput {
			t.Errorf("Got output %q, expected %q", actualOutput, expectedOutput)
		}
		if !memEquals(expectedMemory, memory) {
			t.Errorf("Got memory: (%+v) expected (%+v)", memory, expectedMemory)
		}
	})

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

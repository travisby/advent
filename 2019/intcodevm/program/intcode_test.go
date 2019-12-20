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
		{"halt", []int{99, -1, 0, 8}, halt{}, HALT},
		{"error", []int{-1, 0, 0, 0}, nil, errors.New("Unexpected opcode")},

		{"Position add", []int{1, 10, 20, 30}, add{position{10}, position{20}, position{30}}, nil},
		{"Immediate add", []int{1101, 10, 20, 30}, add{immediate{10}, immediate{20}, position{30}}, nil},
		{"Mixed add", []int{1001, 10, 20, 30}, add{position{10}, immediate{20}, position{30}}, nil},

		{"Position multiply", []int{2, 10, 20, 30}, multiply{position{10}, position{20}, position{30}}, nil},
		{"Immediate multiply", []int{1102, 10, 20, 30}, multiply{immediate{10}, immediate{20}, position{30}}, nil},
		{"Mixed multiply", []int{1002, 10, 20, 30}, multiply{position{10}, immediate{20}, position{30}}, nil},

		{"input", []int{3, 50}, input{parameter1: position{50}}, nil},

		{"output", []int{4, 50}, output{parameter1: position{50}}, nil},

		{"Position equals", []int{8, 10, 20, 30}, equals{position{10}, position{20}, position{30}}, nil},
		{"Immediate equals", []int{1108, 10, 20, 30}, equals{immediate{10}, immediate{20}, position{30}}, nil},
		{"Mixed equals", []int{1008, 10, 20, 30}, equals{position{10}, immediate{20}, position{30}}, nil},

		{"Position less than", []int{7, 10, 20, 30}, lessThan{position{10}, position{20}, position{30}}, nil},
		{"Immediate less than", []int{1107, 10, 20, 30}, lessThan{immediate{10}, immediate{20}, position{30}}, nil},
		{"Mixed less than", []int{1007, 10, 20, 30}, lessThan{position{10}, immediate{20}, position{30}}, nil},
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
			"Halt",
			halt{},
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			HALT,
		},
		{
			"Error condition halts",
			add{position{5}, position{1}, position{2}},
			[]int{0, 0, 0, 0, 0},
			[]int{0, 0, 0, 0, 0},
			ErrUnexpectedHalt,
		},
		{
			"Input",
			input{position{3}, strings.NewReader("-45")},
			[]int{0, 0, 0, 0},
			[]int{0, 0, 0, -45},
			nil,
		},

		{
			"Position Add",
			add{position{10}, position{20}, position{30}},
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3},
			nil,
		},
		{
			"Position Mult",
			multiply{position{10}, position{20}, position{30}},
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
			nil,
		},
		{
			"Position Eq",
			equals{position{10}, position{20}, position{30}},
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			nil,
		},
		{
			"Position Eq (!)",
			equals{position{10}, position{20}, position{30}},
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			nil,
		},
		{
			"Position LT",
			lessThan{position{10}, position{20}, position{30}},
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			nil,
		},
		{
			"Position LT (!)",
			lessThan{position{10}, position{20}, position{30}},
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			nil,
		},

		{
			"Immediate Add",
			add{immediate{10}, immediate{20}, position{0}},
			[]int{5},
			[]int{30},
			nil,
		},
		{
			"Immediate Mult",
			multiply{immediate{10}, immediate{20}, position{0}},
			[]int{5},
			[]int{200},
			nil,
		},
		{
			"Immediate Eq",
			equals{immediate{10}, immediate{10}, position{0}},
			[]int{0},
			[]int{1},
			nil,
		},
		{
			"Immediate Eq (!)",
			equals{immediate{10}, immediate{20}, position{0}},
			[]int{1},
			[]int{0},
			nil,
		},
		{
			"Immediate LT",
			lessThan{immediate{9}, immediate{10}, position{0}},
			[]int{0},
			[]int{1},
			nil,
		},
		{
			"Immediate LT (!)",
			lessThan{immediate{20}, immediate{10}, position{0}},
			[]int{1},
			[]int{0},
			nil,
		},
		{
			"Mixed Add",
			add{immediate{10}, position{0}, position{0}},
			[]int{5},
			[]int{15},
			nil,
		},
		{
			"Mixed Mult",
			multiply{immediate{10}, position{0}, position{0}},
			[]int{5},
			[]int{50},
			nil,
		},
		{
			"Mixed Eq",
			equals{immediate{10}, position{0}, position{0}},
			[]int{10},
			[]int{1},
			nil,
		},
		{
			"Mixed LT",
			lessThan{immediate{10}, position{0}, position{0}},
			[]int{10},
			[]int{0},
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
	t.Run("Position output", func(t *testing.T) {
		buffer := bytes.NewBuffer(nil)
		memory := []int{0, 0, 9, 0}
		expectedMemory := []int{0, 0, 9, 0}
		expectedOutput := "9\n"
		if err := (output{position{2}, buffer}.Apply(memory)); err != nil {
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

	t.Run("Immediate output", func(t *testing.T) {
		buffer := bytes.NewBuffer(nil)
		memory := []int{}
		expectedMemory := []int{}
		expectedOutput := "9\n"
		if err := (output{immediate{9}, buffer}.Apply(memory)); err != nil {
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

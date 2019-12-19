package intcodevm

import (
	"fmt"
	"testing"
)

func TestSimplePrograms(t *testing.T) {
	testCases := []struct {
		memory         []int
		expectedMemory []int
	}{
		{
			[]int{1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50},
			[]int{3500, 9, 10, 70, 2, 3, 11, 0, 99, 30, 40, 50},
		},
		{
			[]int{1, 0, 0, 0, 99},
			[]int{2, 0, 0, 0, 99},
		},
		{
			[]int{2, 3, 0, 3, 99},
			[]int{2, 3, 0, 6, 99},
		},
		{
			[]int{2, 4, 4, 5, 99, 0},
			[]int{2, 4, 4, 5, 99, 9801},
		},
		{
			[]int{1, 1, 1, 4, 99, 5, 6, 0, 99},
			[]int{30, 1, 1, 4, 2, 5, 6, 0, 99},
		},
	}

	for _, tc := range testCases {
		vm := VM{memory: tc.memory}

		if err := vm.Run(); err != nil {
			t.Fatal(err)
		}

		if !memEquals(
			tc.expectedMemory,
			vm.memory) {
			t.Errorf("Expected memory (%+v), got (%+v)", tc.expectedMemory, vm.memory)
		}
	}
}

func TestReset(t *testing.T) {
	vm := VM{
		[]int{1, 2, 3},
		[]int{4, 5, 6},
	}

	if err := vm.Reset(); err != nil {
		t.Fatal(err)
	}
	if !memEquals(vm.roMemory, vm.memory) {
		t.Errorf("Expected memory (%+v), got (%+v)", vm.roMemory, vm.memory)
	}
}

func TestSetNoun(t *testing.T) {
	vm := VM{
		[]int{1, 2, 3},
		[]int{1, 2, 3},
	}

	if err := vm.SetNoun(9); err != nil {
		t.Fatal(err)
	}

	if vm.memory[1] != 9 {
		t.Errorf("Expected noun to be (%d), got (%d)", 9, vm.memory[1])
	}

	// roMemory shouldn't be configured by SetNoun
	if vm.roMemory[1] != 2 {
		t.Errorf("Expected ro[1] to be unchanged from (%d), got (%d)", 2, vm.roMemory[1])
	}
}

func TestNounOverflow(t *testing.T) {
	testCases := []struct {
		memory         []int
		noun           int
		expectOverflow bool
	}{
		{[]int{}, 1, true},
		{[]int{1}, 1, true},
		{[]int{1, 2}, -1, true},
		{[]int{1, 2}, 100, true},
		{[]int{1, 2}, 0, false},
		{[]int{1, 2}, 99, false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("VM{%+v}.SetNoun(%d)", tc.memory, tc.noun), func(t *testing.T) {
			vm := VM{tc.memory, tc.memory}
			err := vm.SetNoun(tc.noun)
			if tc.expectOverflow != (err == ErrOverflow) {
				t.Errorf("Expected overflow (%t) and got (%+v)", tc.expectOverflow, (err == ErrOverflow))
			}
		})
	}
}

func TestVerb(t *testing.T) {
	testCases := []struct {
		memory         []int
		verb           int
		expectOverflow bool
	}{
		{[]int{}, 1, true},
		{[]int{1}, 1, true},
		{[]int{1, 1}, 1, true},
		{[]int{1, 2, 3}, -1, true},
		{[]int{1, 2, 3}, 100, true},
		{[]int{1, 2, 3}, 0, false},
		{[]int{1, 2, 3}, 99, false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("VM{%+v}.SetVerb(%d)", tc.memory, tc.verb), func(t *testing.T) {
			vm := VM{tc.memory, tc.memory}
			err := vm.SetVerb(tc.verb)
			if tc.expectOverflow != (err == ErrOverflow) {
				t.Errorf("Expected overflow (%t) and got (%+v)", tc.expectOverflow, (err == ErrOverflow))
			}
		})
	}
}

func TestSetVerb(t *testing.T) {
	vm := VM{
		[]int{1, 2, 3},
		[]int{1, 2, 3},
	}

	if err := vm.SetVerb(9); err != nil {
		t.Fatal(err)
	}

	if vm.memory[2] != 9 {
		t.Errorf("Expected verb to be (%d), got (%d)", 9, vm.memory[2])
	}

	// roMemory shouldn't be configured by SetVerb
	if vm.roMemory[2] != 3 {
		t.Errorf("Expected ro[2] to be unchanged from (%d), got (%d)", 3, vm.roMemory[2])
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

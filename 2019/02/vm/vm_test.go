package vm

import (
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
			[]int{2, 0, 0, 0, 99, 0, 0, 0},
		},
		{
			[]int{2, 3, 0, 3, 99},
			[]int{2, 3, 0, 6, 99, 0, 0, 0},
		},
		{
			[]int{2, 4, 4, 5, 99, 0},
			[]int{2, 4, 4, 5, 99, 9801, 0, 0},
		},
		{
			[]int{1, 1, 1, 4, 99, 5, 6, 0, 99, 0, 0, 0},
			[]int{30, 1, 1, 4, 2, 5, 6, 0, 99, 0, 0, 0},
		},
	}

	for _, tc := range testCases {
		vm := New(len(tc.memory))
		// for these test cases we want the aligned memory feature
		if err := vm.Load(0, tc.memory); err != nil {
			t.Fatal(err)
		}

		if err := vm.Run(); err != nil {
			t.Fatal(err)
		}

		if !memEquals(
			tc.expectedMemory,
			vm.Memory) {
			t.Errorf("Expected memory (%+v), got (%+v)", tc.expectedMemory, vm.Memory)
		}
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

package program

import (
	"bytes"
	"testing"
)

func PositionOutputTest(t *testing.T) {
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
}

func ImmediateOutputTest(t *testing.T) {
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
}

package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"
)

type testCase struct {
	input  []string
	output struct {
		result *int
		err    error
	}
	method func(r io.Reader) (*int, error)
}

// prints *result or nil, instead of the address
func (tc testCase) outputStringer() string {
	resultString := "<nil>"
	if tc.output.result != nil {
		resultString = fmt.Sprintf("%d", *tc.output.result)
	}
	return fmt.Sprintf("(%s, %v)", resultString, tc.output.err)
}

func (tc testCase) Test(t *testing.T) {
	result, err := tc.method(strings.NewReader(strings.Join(tc.input, "\n")))

	// we don't expect both to be nil, nor neither
	if result == nil && err == nil {
		t.Fatal("Did not expect result & err to be nil")
	} else if result != nil && err != nil {
		t.Fatalf("Did not expect result (%+v) & err (%+v) to not be nil", *result, err)
	}

	if tc.output.result != nil && result != nil && *tc.output.result == *result {
		return
	} else if tc.output.err != nil && err != nil && tc.output.err.Error() == err.Error() {
		return
	}

	temp := "<nil>"
	if result != nil {
		temp = strconv.Itoa(*result)
	}
	t.Fatalf("expected: %v, actual: (%+v, %+v)", tc.outputStringer(), temp, err)
}

// this way we can write test cases inline, avoiding trying to do `&0` which won't work
func newTestCase(input []string, outputSuccess int, outputError error) testCase {
	tc := testCase{input: input}
	if outputError != nil {
		tc.output.err = outputError
	} else {
		tc.output.result = &outputSuccess
	}
	tc.method = addFrequency
	return tc
}

func TestAddFrequency(t *testing.T) {
	for i, tc := range []testCase{
		newTestCase([]string{"+1", "-2", "+3", "+1"}, 3, nil),
		newTestCase([]string{"+1", "+1", "+1"}, 3, nil),
		newTestCase([]string{"-1", "-2", "-3"}, -6, nil),
	} {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tc.Test(t)
		})
	}
}

// this way we can write test cases inline, avoiding trying to do `&0` which won't work
func newTestCaseP2(input []string, outputSuccess int, outputError error) testCase {
	tc := testCase{input: input}
	if outputError != nil {
		tc.output.err = outputError
	} else {
		tc.output.result = &outputSuccess
	}
	tc.method = firstDoubleFrequency
	return tc
}

func TestFirstDoubleFrequency(t *testing.T) {
	for i, tc := range []testCase{
		newTestCaseP2([]string{"+1", "-1"}, 0, nil),
		newTestCaseP2([]string{"+3", "+3", "+4", "-2", "-4"}, 10, nil),
		newTestCaseP2([]string{"-6", "+3", "+8", "+5", "-6"}, 5, nil),
		newTestCaseP2([]string{"+7", "+7", "-2", "-7", "-4"}, 14, nil),
	} {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tc.Test(t)
		})
	}
}

package main

import (
	"bufio"
	"strings"
	"testing"
)

func stringSliceEquals(a, b []string) bool {
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

func TestGuardShiftScanner(t *testing.T) {
	tc := struct {
		input  string
		output []string
		err    error
	}{
		`[1518-02-12 23:50] Guard #1789 begins shift
[1518-02-13 00:05] falls asleep
[1518-02-13 00:25] wakes up
[1518-02-13 00:40] falls asleep
[1518-02-13 00:52] wakes up
[1518-02-14 00:01] Guard #2617 begins shift
[1518-02-14 00:22] falls asleep
[1518-02-14 00:29] wakes up
[1518-02-14 00:47] falls asleep
[1518-02-14 00:50] wakes up
`,
		[]string{
			`[1518-02-12 23:50] Guard #1789 begins shift
[1518-02-13 00:05] falls asleep
[1518-02-13 00:25] wakes up
[1518-02-13 00:40] falls asleep
[1518-02-13 00:52] wakes up
`,
			`[1518-02-14 00:01] Guard #2617 begins shift
[1518-02-14 00:22] falls asleep
[1518-02-14 00:29] wakes up
[1518-02-14 00:47] falls asleep
[1518-02-14 00:50] wakes up
`,
		}, nil}
	// testing the actual split func sounds insanely hard
	// let's just wrap it in a scanner and make sure _that_ works ok!
	scanner := bufio.NewScanner(strings.NewReader(tc.input))
	scanner.Split(scanGuardShift)

	// store everything we scanned
	actualOutputs := []string{}
	for scanner.Scan() {
		actualOutputs = append(actualOutputs, scanner.Text())
	}

	// ugg, this is annoying
	actualErr := "<nil>"
	expectedErr := "<nil>"
	if scanner.Err() != nil {
		actualErr = scanner.Err().Error()
	}
	if tc.err != nil {
		expectedErr = tc.err.Error()
	}

	if actualErr != expectedErr {
		t.Fatalf("Expected (%+v) to match (%+v)", actualErr, expectedErr)
	} else if !stringSliceEquals(actualOutputs, tc.output) {
		t.Fatalf("Expected (%+v) to match (%+v)", actualOutputs, tc.output)
	}
}

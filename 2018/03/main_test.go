package main

import (
	"fmt"
	"testing"
)

func TestNewClaim(t *testing.T) {
	type testCase struct {
		input  string
		output *claim
		err    error
	}
	for i, tc := range []testCase{
		testCase{"#123 @ 3,2: 5x4", &claim{123, struct {
			x int
			y int
		}{3, 2}, 5, 4}, nil},
		testCase{"#1 @ 1,3: 4x4", &claim{1, struct {
			x int
			y int
		}{1, 3}, 4, 4}, nil},
		testCase{"#2 @ 3,1: 4x4", &claim{2, struct {
			x int
			y int
		}{3, 1}, 4, 4}, nil},
		testCase{"#3 @ 5,5: 2x2", &claim{3, struct {
			x int
			y int
		}{5, 5}, 2, 2}, nil},
	} {
		// TODO handle errors in test cases...
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			actual, err := newClaim(tc.input)
			if err != nil {
				t.Fatal(err)
			} else if actual == nil {
				t.Fatalf("Did not expect claim to be nil")
			} else if *actual != *tc.output {
				t.Fatalf("Expected (%+v) to be (%+v)", *actual, *tc.output)
			}
		})
	}

}

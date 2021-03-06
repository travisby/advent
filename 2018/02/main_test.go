package main

import (
	"fmt"
	"testing"
)

func mapIntEquality(a, b map[int]int) bool {
	if len(a) != len(b) {
		return false
	}

	for k := range a {
		// avoid nils
		if _, ok := b[k]; !ok {
			return false
		} else if a[k] != b[k] {
			return false
		}
	}
	return true
}

func TestCountExactlyRepeatedLetters(t *testing.T) {
	type testCase struct {
		input  string
		output map[int]int
	}

	for i, tc := range []testCase{
		testCase{"abcdef", map[int]int{1: 6}},
		testCase{"bababc", map[int]int{1: 1, 2: 1, 3: 1}},
		testCase{"abbcde", map[int]int{1: 4, 2: 1}},
		testCase{"abcccd", map[int]int{1: 3, 3: 1}},
		testCase{"aabcdd", map[int]int{1: 2, 2: 2}},
		testCase{"abcdee", map[int]int{1: 4, 2: 1}},
		testCase{"ababab", map[int]int{3: 2}},
	} {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			if count := countExactlyRepeatedLetters(tc.input); !mapIntEquality(tc.output, count) {
				t.Fatalf("Expected %d, got %d", tc.output, count)
			}
		})
	}
}

func TestChecksumsTwosAndThrees(t *testing.T) {
	type testCase struct {
		input  []map[int]int
		output int
	}
	tcs := []testCase{testCase{
		[]map[int]int{
			map[int]int{1: 6},
			map[int]int{1: 1, 2: 1, 3: 1},
			map[int]int{1: 4, 2: 1},
			map[int]int{1: 3, 3: 1},
			map[int]int{1: 2, 2: 2},
			map[int]int{1: 4, 2: 1},
			map[int]int{3: 2},
		},
		12,
	}}
	for _, tc := range tcs {
		if result := checksumTwosAndThrees(tc.input); tc.output != result {
			t.Fatalf("Expected %d, got %d", tc.output, result)
		}
	}

}

func TestNumDifferent(t *testing.T) {
	type testCase struct {
		input  []string
		output int
	}

	for i, tc := range []testCase{
		testCase{[]string{"abcde", "axcye"}, 2},
		testCase{[]string{"fghij", "fguij"}, 1},
	} {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			if result := numDifferent(tc.input[0], tc.input[1]); result != tc.output {
				t.Fatalf("Expected %d, got %d", tc.output, result)
			}
		})
	}
}

func TestFindOnlyOneDifferent(t *testing.T) {
	str1, str2 := findOnlyOneDifferent([]string{"abcde", "axcye", "fghij", "fguij"})
	if str1 == str2 {
		t.Fatalf("Expected %q and %q to differ", str1, str2)
	} else if str1 != "fghij" && str2 != "fghij" {
		t.Fatalf("expected %q or %q to be \"fghij\"", str1, str2)
	} else if str1 != "fguij" && str2 != "fguij" {
		t.Fatalf("expected %q or %q to be \"fguij\"", str1, str2)
	}
}

func TestFindCommonCharacters(t *testing.T) {
	result := findCommonCharacters("fghij", "fguij")
	if result != "fgij" {
		t.Fatalf("Expected %q to be \"fgij\"", result)
	}
}

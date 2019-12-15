package main

import (
	"fmt"
	"log"
	"testing"
)

func TestGetCrossings(t *testing.T) {
	crossings := getCrossings(
		[]instruction{right(8), up(5), left(5), down(3)},
		[]instruction{up(7), right(6), down(4), left(4)},
	)

	if !pointsEquals([]point{point{3, 3}, point{6, 5}}, crossings) {
		t.Fatalf("Expected crossings to be (%+v), got (%+v)", []point{point{3, 3}, point{6, 5}}, crossings)
	}
}

func TestManhattanSort(t *testing.T) {
	testCases := []struct {
		points        []point
		expectedOrder []point
	}{
		{
			[]point{point{3, 3}, point{6, 5}},
			[]point{point{3, 3}, point{6, 5}},
		},
		{
			[]point{point{6, 5}, point{3, 3}},
			[]point{point{3, 3}, point{6, 5}},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%+v", tc.points), func(t *testing.T) {
			manhattanSort(tc.points)
			if !pointsEquals(tc.points, tc.expectedOrder) {
				t.Fatalf("Expected (%+v), got (%+v)", tc.expectedOrder, tc.points)
			}
		})
	}
}

func TestGetClosestCrossingsDistance(t *testing.T) {
	testCases := []struct {
		instructions     [2][]instruction
		expectedDistance int
	}{
		{
			[2][]instruction{
				[]instruction{right(75), down(30), right(83), up(83), left(12), down(49), right(71), up(7), left(72)},
				[]instruction{up(62), right(66), up(55), right(34), down(71), right(55), down(58), right(83)},
			},
			159,
		},
		{
			[2][]instruction{
				[]instruction{right(98), up(47), right(26), down(63), right(33), up(87), left(62), down(20), right(33), up(53), right(51)},
				[]instruction{up(98), right(91), down(20), right(16), down(67), right(40), up(7), right(15), up(6), right(7)},
			},
			135,
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%+v", tc.instructions), func(t *testing.T) {
			closest, err := getClosestCrossingsDistance(tc.instructions[0], tc.instructions[1])
			if err != nil {
				log.Fatal(err)
			}
			if *closest != tc.expectedDistance {
				t.Errorf("expected (%d) got (%d)", tc.expectedDistance, *closest)
			}
		})
	}
}

func pointsEquals(ps1, ps2 []point) bool {
	if len(ps1) != len(ps2) {
		return false
	}
	for i := range ps1 {
		if ps1[i] != ps2[i] {
			return false
		}
	}
	return true
}

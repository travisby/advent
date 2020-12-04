package main

import (
	"errors"
	"fmt"
	"testing"
)

func assertLayersEqual(actual, expected *treeLayer) error {
	if actual == nil && expected != nil {
		return fmt.Errorf("Unexpected nil treeLayer")
	} else if actual != nil && expected == nil {
		return fmt.Errorf("expected nil treeLayer got %+v", *actual)
	} else if actual != nil && expected != nil {
		if len(actual.trees) != len(expected.trees) {
			return fmt.Errorf("Expected equal layers, got actual=%+v expected=%+v", actual.trees, expected.trees)
		}
		for i := range actual.trees {
			if actual.trees[i] != expected.trees[i] {
				return fmt.Errorf("Expected equal layers, got actual=%+v expected=%+v", actual.trees, expected.trees)
			}
		}
	}

	return nil
}

func TestNewTreeLayer(t *testing.T) {
	type test struct {
		name  string
		input string
		err   error
		layer *treeLayer
	}
	tests := []test{
		{
			"Base Case",
			"",
			nil,
			&treeLayer{[]bool{}},
		},
		{
			"Simple Case 1",
			".",
			nil,
			&treeLayer{[]bool{false}},
		},
		{
			"Simple Case 2",
			"#",
			nil,
			&treeLayer{[]bool{true}},
		},
		{
			"Err Case",
			"@",
			ErrUnknownTreeSyntax,
			nil,
		},
	}

	// here's the sample input given to us from AoC
	for i, sample := range []struct {
		s  string
		bs []bool
	}{
		{"..##.......", []bool{false, false, true, true, false, false, false, false, false, false, false}},
		{"#...#...#..", []bool{true, false, false, false, true, false, false, false, true, false, false}},
		{".#....#..#.", []bool{false, true, false, false, false, false, true, false, false, true, false}},
		{"..#.#...#.#", []bool{false, false, true, false, true, false, false, false, true, false, true}},
		{".#...##..#.", []bool{false, true, false, false, false, true, true, false, false, true, false}},
		{"..#.##.....", []bool{false, false, true, false, true, true, false, false, false, false, false}},
		{".#.#.#....#", []bool{false, true, false, true, false, true, false, false, false, false, true}},
		{".#........#", []bool{false, true, false, false, false, false, false, false, false, false, true}},
		{"#.##...#...", []bool{true, false, true, true, false, false, false, true, false, false, false}},
		{"#...##....#", []bool{true, false, false, false, true, true, false, false, false, false, true}},
		{".#..#...#.#", []bool{false, true, false, false, true, false, false, false, true, false, true}},
	} {
		tests = append(
			tests,
			test{
				fmt.Sprintf("AOC Sample %d", i),
				sample.s,
				nil,
				&treeLayer{sample.bs},
			},
		)
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			layer, err := newTreeLayer(test.input)
			if !errors.Is(err, test.err) {
				t.Fatalf("Expected %v but got %v", test.err, err)
			}
			if err := assertLayersEqual(layer, test.layer); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestNewMapIsEmpty(t *testing.T) {
	m := newTreeMap()
	if m == nil {
		t.Fatal("Did not expect tree map to be nil")
	} else if len(m.layers) != 0 {
		t.Fatalf("Expected there to be zero layers. Got %d layers", len(m.layers))
	}
}

func TestAddLayer(t *testing.T) {
	m := treeMap{}
	treeLayer := treeLayer{[]bool{true}}
	m.addLayer(treeLayer)

	if len(m.layers) != 1 {
		t.Fatalf("Expected there to be one and only one layer now. Got %d layers", len(m.layers))
	}

	if err := assertLayersEqual(&m.layers[0], &treeLayer); err != nil {
		t.Fatal(err)
	}
}

func TestTraversal(t *testing.T) {
	type test struct {
		name    string
		treeMap treeMap
		input   struct {
			x int
			y int
		}
		more   bool
		curPos struct {
			x int
			y int
		}
		trees int
	}
	tests := []test{
		{
			"Simple Case",
			treeMap{layers: []treeLayer{treeLayer{[]bool{false}}, treeLayer{[]bool{false}}}},
			struct {
				x int
				y int
			}{0, 0},
			true,
			struct {
				x int
				y int
			}{0, 0},
			0,
		},
		{
			"Down 1 Right 1 Case No Trees",
			treeMap{layers: []treeLayer{treeLayer{[]bool{false, false}}, treeLayer{[]bool{false, false}}}},
			struct {
				x int
				y int
			}{1, 1},
			false,
			struct {
				x int
				y int
			}{1, 1},
			0,
		},
		{
			"Down 1 Right 2 Case Some Trees",
			treeMap{layers: []treeLayer{treeLayer{[]bool{true, false, true}}, treeLayer{[]bool{false, false, true}}}},
			struct {
				x int
				y int
			}{2, 1},
			false,
			struct {
				x int
				y int
			}{2, 1},
			1,
		},
		{
			"Moving Right should be a repeating pattern",
			treeMap{layers: []treeLayer{treeLayer{[]bool{true, false}}}},
			struct {
				x int
				y int
			}{10, 0},
			false,
			struct {
				x int
				y int
			}{10, 0},
			1,
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			more := tst.treeMap.traverse(tst.input.x, tst.input.y)
			if more != tst.more {
				t.Errorf("Expected more %v got %v", tst.more, more)
			}
			if tst.treeMap.curPos != tst.curPos {
				t.Errorf("Expected new curPos to be: %+v got %+v", tst.treeMap.curPos, tst.curPos)
			}
			if tst.trees != tst.treeMap.treesEncountered {
				t.Errorf("Expected trees %d got %d", tst.trees, tst.treeMap.treesEncountered)
			}
		})
	}
}

func TestReset(t *testing.T) {
	m := treeMap{curPos: struct {
		x int
		y int
	}{5, 5},
		treesEncountered: 5,
	}
	m.reset()

	if m.curPos != struct {
		x int
		y int
	}{0, 0} {
		t.Errorf("Expected curPos to be reset, was actually %+v", m.curPos)
	}

	if m.treesEncountered != 0 {
		t.Errorf("Expected treesEncoutnered to be reset, was actually %+v", m.treesEncountered)
	}
}

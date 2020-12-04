package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
)

var ErrUnknownTreeSyntax = errors.New("Unknown tree syntax")
var ErrTraversal = errors.New("Traversing into the unknown")

type treeLayer struct {
	trees []bool
}

func newTreeLayer(s string) (*treeLayer, error) {
	t := treeLayer{make([]bool, len(s))}

	for i, c := range s {
		if c == '.' {
			t.trees[i] = false
		} else if c == '#' {
			t.trees[i] = true
		} else {
			return nil, fmt.Errorf("%w at pos %d, %q", ErrUnknownTreeSyntax, i, c)
		}
	}

	return &t, nil
}

type treeMap struct {
	layers []treeLayer
	curPos struct {
		x int
		y int
	}
	treesEncountered int
}

func (t *treeMap) reset() {
	t.curPos.x = 0
	t.curPos.y = 0
	t.treesEncountered = 0
}

func (t *treeMap) addLayer(ts treeLayer) {
	t.layers = append(t.layers, ts)
}

// traverse returns false if there's no where left to traverse
// (AKA we've reached the bottom)
// traversel L-R is infinite
// traversal will return false and prevent overtravel
func (t *treeMap) traverse(x, y int) (more bool) {
	if t.curPos.y+y >= len(t.layers) {
		return false
	}

	t.curPos.x += x
	t.curPos.y += y

	if t.layers[t.curPos.y].trees[t.curPos.x%len(t.layers[t.curPos.y].trees)] {
		t.treesEncountered++
	}

	if t.curPos.y == len(t.layers)-1 {
		return false
	}

	return true
}

func newTreeMap() *treeMap {
	return &treeMap{}
}

func main() {
	var f *os.File
	if len(os.Args) == 2 {
		var err error
		f, err = os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	} else {
		f = os.Stdin
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	treeMap := newTreeMap()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		layer, err := newTreeLayer(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		treeMap.addLayer(*layer)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for more := true; more; more = treeMap.traverse(3, 1) {
	}

	log.Printf("Part One: %d", treeMap.treesEncountered)

	// Part 2, get all of these plans and multiply their results
	p2 := 1
	treeMap.reset()
	for more := true; more; more = treeMap.traverse(1, 1) {
	}
	p2 *= treeMap.treesEncountered

	treeMap.reset()
	for more := true; more; more = treeMap.traverse(3, 1) {
	}
	p2 *= treeMap.treesEncountered

	treeMap.reset()
	for more := true; more; more = treeMap.traverse(5, 1) {
	}
	p2 *= treeMap.treesEncountered

	treeMap.reset()
	for more := true; more; more = treeMap.traverse(7, 1) {
	}
	p2 *= treeMap.treesEncountered

	treeMap.reset()
	for more := true; more; more = treeMap.traverse(1, 2) {
	}
	p2 *= treeMap.treesEncountered

	log.Printf("Part Two: %d", p2)
}

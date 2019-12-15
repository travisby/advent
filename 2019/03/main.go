package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type point struct {
	x int
	y int
}
type instruction point

var ErrUnknownInstruction = errors.New("Unknown instruction")

func parseInstruction(s string) (*instruction, error) {
	if len(s) < 2 {
		return nil, ErrUnknownInstruction
	}

	i, err := strconv.Atoi(s[1:])
	if err != nil {
		return nil, err
	}

	var inst instruction
	switch s[0] {
	case 'U':
		inst = up(i)
	case 'D':
		inst = down(i)
	case 'L':
		inst = left(i)
	case 'R':
		inst = right(i)
	default:
		return nil, ErrUnknownInstruction
	}
	return &inst, nil
}

func parseInstructions(instructions string) ([]instruction, error) {
	ss := strings.Split(instructions, ",")

	is := make([]instruction, 0, len(ss))

	for _, s := range ss {
		i, err := parseInstruction(s)
		if err != nil {
			return nil, err
		}
		is = append(is, *i)
	}

	return is, nil
}

func (i instruction) step(p point) []point {
	points := []point{}

	newPoint := point{p.x, p.y}
	for _ = 0; i.x > 0; i.x-- {
		newPoint.x++
		points = append(points, point{newPoint.x, newPoint.y})
	}
	for _ = 0; i.x < 0; i.x++ {
		newPoint.x--
		points = append(points, point{newPoint.x, newPoint.y})
	}
	for _ = 0; i.y > 0; i.y-- {
		newPoint.y++
		points = append(points, point{newPoint.x, newPoint.y})
	}
	for _ = 0; i.y < 0; i.y++ {
		newPoint.y--
		points = append(points, point{newPoint.x, newPoint.y})
	}
	return points
}

func up(y int) instruction {
	return instruction(point{0, y})
}
func down(y int) instruction {
	return instruction(point{0, -y})
}
func left(x int) instruction {
	return instruction(point{-x, 0})
}
func right(x int) instruction {
	return instruction(point{x, 0})
}

func getCrossings(a []instruction, b []instruction) []point {
	pos1 := accumulatePoints(a)
	pos2 := accumulatePoints(b)

	// create some maps to search for existence
	pos1Map := map[point]bool{}
	for _, p := range pos1 {
		pos1Map[p] = true
	}
	pos2Map := map[point]bool{}
	for _, p := range pos2 {
		pos2Map[p] = true
	}

	crossings := []point{}
	for k := range pos1Map {
		if _, ok := pos2Map[k]; ok {
			crossings = append(crossings, k)
		}
	}

	manhattanSort(crossings)
	return crossings
}

func getClosestCrossingsDistance(a []instruction, b []instruction) (*int, error) {
	crossings := getCrossings(a, b)
	if len(crossings) < 1 {
		return nil, fmt.Errorf("No crossings")
	}

	distance := manhattanDistance(crossings[0])
	return &distance, nil
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

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		log.Fatalf("Expected instructions for wire 1")
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// parse wire1
	instructions0, err := parseInstructions(scanner.Text())
	if err != nil {
		log.Fatal(err)
	}

	if !scanner.Scan() {
		log.Fatalf("Expected instructions for wire 2")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// parse wire1
	instructions1, err := parseInstructions(scanner.Text())
	if err != nil {
		log.Fatal(err)
	}

	if scanner.Scan() {
		log.Fatalf("Unexpected additional data")
	}

	distance, err := getClosestCrossingsDistance(instructions0, instructions1)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Closest crossing distance: %d", *distance)
}

func manhattanSort(ps []point) {
	sort.Slice(
		ps,
		func(i, j int) bool {
			return manhattanDistance(ps[i]) < manhattanDistance(ps[j])
		},
	)
}

func manhattanDistance(p point) int {
	return int(math.Abs(float64(p.x-0))) + int(math.Abs(float64(p.y-0)))
}

// accumulate in this case is like python accumulate, or Haskell scan
// we want to perform an application like a reduce/fold, but we also want the intermediate results
// since we not only want to "get" to the end of the instructions, but we want to see each place we went to get there
// to account for "R8" counting as only one point movement, we'll actually treat that as 8*R(1) for the sake of internal values
func accumulatePoints(instructions []instruction) []point {
	accum := []point{}

	p := point{}
	for _, inst := range instructions {
		accum = append(accum, inst.step(p)...)
		p = accum[len(accum)-1]
	}

	return accum
}

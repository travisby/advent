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
type visit struct {
	point
	visitsAt int
}
type wire struct {
	visits []visit
}

func (w *wire) up() {
	last := w.visits[len(w.visits)-1]
	next := last
	next.y++
	next.visitsAt++
	w.visits = append(w.visits, next)
}
func (w *wire) down() {
	last := w.visits[len(w.visits)-1]
	next := last
	next.y--
	next.visitsAt++
	w.visits = append(w.visits, next)
}
func (w *wire) left() {
	last := w.visits[len(w.visits)-1]
	next := last
	next.x--
	next.visitsAt++
	w.visits = append(w.visits, next)
}
func (w *wire) right() {
	last := w.visits[len(w.visits)-1]
	next := last
	next.x++
	next.visitsAt++
	w.visits = append(w.visits, next)
}

func (w *wire) travel(i instruction) {
	for _ = 0; i.x > 0; i.x-- {
		w.right()
	}
	for _ = 0; i.x < 0; i.x++ {
		w.left()
	}
	for _ = 0; i.y > 0; i.y-- {
		w.up()
	}
	for _ = 0; i.y < 0; i.y++ {
		w.down()
	}
}

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

func getCrossings(a []instruction, b []instruction) []visit {
	wire1Visits := accumulateVisits(a)
	wire2Visits := accumulateVisits(b)

	// create some maps to search for visits by index
	wire1Map := make(map[point]int, len(wire1Visits))
	wire2Map := make(map[point]int, len(wire2Visits))

	for _, v := range wire1Visits {
		// only save the first visitsAt per wire
		if _, ok := wire1Map[v.point]; !ok {
			wire1Map[v.point] = v.visitsAt
		}
	}
	for _, v := range wire2Visits {
		// only save the first visitsAt per wire
		if _, ok := wire2Map[v.point]; !ok {
			wire2Map[v.point] = v.visitsAt
		}
	}

	crossings := []visit{}
	for k, v1 := range wire1Map {
		if v2, ok := wire2Map[k]; ok {
			crossings = append(crossings, visit{k, v1 + v2})
		}
	}

	manhattanSort(crossings)

	return crossings
}

func getClosestCrossingsDistance(a []instruction, b []instruction) (*int, error) {
	crossings := getCrossings(a, b)
	if len(crossings) < 2 {
		return nil, fmt.Errorf("No crossings")
	}

	// everything technically crosses at 0, so del this
	crossings = crossings[1:]

	distance := manhattanDistance(crossings[0].point)
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

func manhattanSort(vs []visit) {
	sort.Slice(
		vs,
		func(i, j int) bool {
			return manhattanDistance(vs[i].point) < manhattanDistance(vs[j].point)
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
func accumulateVisits(instructions []instruction) []visit {
	wire := wire{[]visit{visit{}}}

	for _, i := range instructions {
		wire.travel(i)
	}

	return wire.visits
}

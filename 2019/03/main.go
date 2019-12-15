package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type point struct {
	x int
	y int
}
type instruction point

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

func GetCrossings([]instruction, []instruction) []point {
	// TODO not implemented
	return nil
}

func GetClosestCrossingsDistance([]instruction, []instruction) (*int, error) {
	return nil, fmt.Errorf("Not implemented")
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
	for scanner.Scan() {
		_ = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func manhattanSort([]point) {
	// TODO
}

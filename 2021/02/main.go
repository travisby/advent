package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type instruction struct {
	horizontal int
	aim        int
}

func forward(i int) instruction {
	return instruction{horizontal: i}
}
func up(i int) instruction {
	return instruction{aim: -i}
}
func down(i int) instruction {
	return instruction{aim: i}
}

type position struct {
	instruction
	depth int
}

func (p *position) apply(i instruction) {
	p.aim += i.aim

	p.horizontal += i.horizontal
	p.depth += p.aim * i.horizontal
}

func newPosition() position {
	return position{}
}

func strToInstruction(s string) (*instruction, error) {
	var i int
	var j string

	if n, err := fmt.Sscanf(s, "%s %d", &j, &i); err != nil {
		return nil, err
	} else if n != 2 {
		return nil, fmt.Errorf("Unexpected")
	}

	var f func(i int) instruction
	switch j {
	case "forward":
		f = forward
	case "down":
		f = down
	case "up":
		f = up
	default:
		return nil, fmt.Errorf("Unexpected")
	}

	temp := f(i)

	return &temp, nil
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

	pos := newPosition()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		i, err := strToInstruction(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		pos.apply(*i)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Part 2: %d", pos.horizontal*pos.depth)
}

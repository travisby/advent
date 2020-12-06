package main

import (
	"bufio"
	"log"
	"os"
)

// a person is identified by all of the Questions they answered yes to
// where the Q is identified by its rune name
type person []byte

func toPerson(s string) person {
	p := make([]byte, len(s))
	for i := range s {
		p[i] = s[i]
	}
	return p
}

// a group is simply a list of people that together answer questions
type group []person

// count for a group is the number of total unique YES's in a group
// e.g. [[ab] [ac]] -> 3
// this satisfies the count for P1
func (g group) count() uint {
	counter := map[byte]struct{}{}
	for _, p := range g {
		for _, r := range p {
			counter[r] = struct{}{}
		}
	}
	return uint(len(counter))
}

// count for a group is the number of unanimous YES'
// e.g. [[ab] [ac]] -> 1
// this satisfies the count for P2
func (g group) countAllAnswered() uint {
	counter := map[byte]uint{}
	for _, p := range g {
		for _, r := range p {
			counter[r] += 1
		}
	}

	var count uint
	for _, v := range counter {
		if v == uint(len(g)) {
			count += 1
		}
	}
	return count
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

	var gs []group
	var g group

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Text() == "" {
			gs = append(gs, g)
			g = group{}
			continue
		}
		g = append(g, toPerson(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// handle last group on ending input!
	gs = append(gs, g)

	var count uint
	for _, g := range gs {
		count += g.count()
	}

	log.Printf("Part 1: %d", count)

	count = 0
	for _, g := range gs {
		count += g.countAllAnswered()
	}

	log.Printf("Part 2: %d", count)
}

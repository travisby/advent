package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type seat rune

var FLOOR seat = '.'
var EMPTY seat = 'L'
var OCCUPIED seat = '#'
var ErrInvalidSeat = errors.New("ErrInvalidSeat")

type seatLayout [][]seat

// String() aids in debugging
// while also being useful as a hashing function
// e.g. let us compare equality of layouts
func (s seatLayout) String() string {
	if len(s) == 0 {
		return ""
	}

	var b strings.Builder
	// since we know each row is the same size
	// we can assume our size is col*row
	// and just account for potential already capacity
	// (which I think is unlikely, but meh)
	b.Grow((len(s) * len(s[0])) - b.Cap())
	for _, i := range s {
		for _, j := range i {
			// ignore err
			b.WriteRune(rune(j))
		}
		// ignore err
		b.WriteRune('\n')
	}
	// remove the extra \n
	return b.String()[:b.Len()-1]

}

func (a seatLayout) DeepCopy() (b seatLayout) {
	b = make([][]seat, len(a))
	for i := range a {
		b[i] = make([]seat, len(a[i]))
		copy(b[i], a[i])
	}
	return b
}

func (a seatLayout) PerformRound() seatLayout {
	/*
		If a seat is empty (L) and there are no occupied seats adjacent to it, the seat becomes occupied.
		If a seat is occupied (#) and four or more seats adjacent to it are also occupied, the seat becomes empty.
		Otherwise, the seat's state does not change.
	*/

	b := a.DeepCopy()

	for i := range b {
		for j := range a[i] {
			if a[i][j] == EMPTY {
				if a.numAdjacentOccupied(i, j) == 0 {
					b[i][j] = OCCUPIED
				}
			} else if a[i][j] == OCCUPIED {
				if a.numAdjacentOccupied(i, j) >= 4 {
					b[i][j] = EMPTY
				}
			} else if a[i][j] == FLOOR {
				// ignore
			}
		}
	}

	return b
}

func (a seatLayout) PerformRound2() seatLayout {
	/*
		If a seat is empty (L) and there are no occupied seats adjacent to it, the seat becomes occupied.
		If a seat is occupied (#) and four or more seats adjacent to it are also occupied, the seat becomes empty.
		Otherwise, the seat's state does not change.
	*/

	b := a.DeepCopy()

	for i := range b {
		for j := range a[i] {
			if a[i][j] == EMPTY {
				if a.numInSightOccupied(i, j) == 0 {
					b[i][j] = OCCUPIED
				}
			} else if a[i][j] == OCCUPIED {
				if a.numInSightOccupied(i, j) >= 5 {
					b[i][j] = EMPTY
				}
			} else if a[i][j] == FLOOR {
				// ignore
			}
		}
	}

	return b
}

func (a seatLayout) numAdjacentOccupied(i, j int) int {
	adjacents := make([]seat, 0, 8)

	// top row
	if i-1 >= 0 {
		if j-1 >= 0 {
			adjacents = append(adjacents, a[i-1][j-1])
		}
		adjacents = append(adjacents, a[i-1][j])
		if j+1 < len(a[i-1]) {
			adjacents = append(adjacents, a[i-1][j+1])
		}
	}
	// my row
	if j-1 >= 0 {
		adjacents = append(adjacents, a[i][j-1])
	}
	if j+1 < len(a[i]) {
		adjacents = append(adjacents, a[i][j+1])
	}
	// bottom row
	if i+1 < len(a) {
		if j-1 >= 0 {
			adjacents = append(adjacents, a[i+1][j-1])
		}
		adjacents = append(adjacents, a[i+1][j])
		if j+1 < len(a[i+1]) {
			adjacents = append(adjacents, a[i+1][j+1])
		}
	}

	countOccupied := 0
	for _, v := range adjacents {
		if v == OCCUPIED {
			countOccupied++
		}
	}
	return countOccupied
}

func (a seatLayout) numInSightOccupied(y, x int) int {
	count := 0

	// north
	for i := y - 1; i >= 0; i-- {
		if a[i][x] == OCCUPIED {
			count++
			break
		} else if a[i][x] == EMPTY {
			break
		}
	}

	// south
	for i := y + 1; i < len(a); i++ {
		if a[i][x] == OCCUPIED {
			count++
			break
		} else if a[i][x] == EMPTY {
			break
		}
	}

	// east
	for i := x + 1; i < len(a[y]); i++ {
		if a[y][i] == OCCUPIED {
			count++
			break
		} else if a[y][i] == EMPTY {
			break
		}
	}

	// west
	for i := x - 1; i >= 0; i-- {
		if a[y][i] == OCCUPIED {
			count++
			break
		} else if a[y][i] == EMPTY {
			break
		}
	}

	var i, j int
	i = y + 1
	j = x + 1
	for i >= 0 && j >= 0 && i < len(a) && j < len(a[i]) {
		if a[i][j] == OCCUPIED {
			count++
			break
		} else if a[i][j] == EMPTY {
			break
		}
		i++
		j++
	}

	i = y + 1
	j = x - 1
	for i >= 0 && j >= 0 && i < len(a) && j < len(a[i]) {
		if a[i][j] == OCCUPIED {
			count++
			break
		} else if a[i][j] == EMPTY {
			break
		}
		i++
		j--
	}

	i = y - 1
	j = x + 1
	for i >= 0 && j >= 0 && i < len(a) && j < len(a[i]) {
		if a[i][j] == OCCUPIED {
			count++
			break
		} else if a[i][j] == EMPTY {
			break
		}
		i--
		j++
	}

	i = y - 1
	j = x - 1
	for i >= 0 && j >= 0 && i < len(a) && j < len(a[i]) {
		if a[i][j] == OCCUPIED {
			count++
			break
		} else if a[i][j] == EMPTY {
			break
		}
		i--
		j--
	}

	return count
}

func (a seatLayout) numOccupied() int {
	numOccupied := 0
	for _, i := range a {
		for _, j := range i {
			if j == OCCUPIED {
				numOccupied++
			}
		}
	}
	return numOccupied
}

func newSeat(c rune) (*seat, error) {
	switch c {
	case rune(FLOOR):
		return &FLOOR, nil
	case rune(EMPTY):
		return &EMPTY, nil
	}
	return nil, fmt.Errorf("%w: %q", ErrInvalidSeat, c)
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

	var layout seatLayout

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		seats := make([]seat, 0, len(scanner.Text()))

		for _, c := range scanner.Text() {
			s, err := newSeat(c)
			if err != nil {
				log.Fatal(err)
			}
			seats = append(seats, *s)
		}

		layout = append(layout, seats)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	layoutCopy := layout.DeepCopy()

	var stable bool
	for !stable {
		temp := layout.PerformRound()
		if layout.String() == temp.String() {
			stable = true
		}
		layout = temp
	}

	log.Printf("Part 1: %d", layout.numOccupied())

	layout = layoutCopy
	stable = false
	for !stable {
		temp := layout.PerformRound2()
		if layout.String() == temp.String() {
			stable = true
		}
		layout = temp
	}

	log.Printf("Part 2: %d", layout.numOccupied())
}

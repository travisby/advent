package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Satellite struct {
	ID int
	ps []Point
}

func (s Satellite) Rotated(r Rotation) Satellite {
	newS := Satellite{s.ID, make([]Point, 0, len(s.ps))}
	for _, p := range s.ps {
		newS.ps = append(newS.ps, p.Apply(r))
	}

	return newS
}
func (s Satellite) String() string {
	strs := make([]string, 0, len(s.ps)+1)
	strs = append(strs, fmt.Sprintf("--- scanner %d ---", s.ID))
	for _, p := range s.ps {
		strs = append(strs, fmt.Sprintf("%d,%d,%d", p.x, p.y, p.z))
	}
	return strings.Join(strs, "\n")
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

	var satellites []Satellite
	for scanner.Scan() {
		var satellite Satellite

		if n, err := fmt.Sscanf(scanner.Text(), "--- scanner %d ---", &satellite.ID); n != 1 || err != nil {
			log.Fatal(err)
		}

		// keep scanning either until EOF or we hit a newline
		// that means next coming is a new satellite
		for scanner.Scan() && scanner.Text() != "" {
			var point Point
			if n, err := fmt.Sscanf(scanner.Text(), "%d,%d,%d", &point.x, &point.y, &point.z); n != 3 || err != nil {
				log.Fatal(err)
			}

			satellite.ps = append(satellite.ps, point)
		}
		satellites = append(satellites, satellite)
	}

	rotations := getRotations()
	for _, s := range satellites {
		strs := make([]string, 0, len(rotations))
		for _, r := range rotations {
			strs = append(strs, s.Rotated(r).String())
		}
		fmt.Println(strings.Join(strs, "\n\n"))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

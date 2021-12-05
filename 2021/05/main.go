package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

type Point struct {
	x, y int
}

type Line [2]Point

func parseLine(str string) (Line, error) {
	var l Line
	if n, err := fmt.Sscanf(str, "%d,%d -> %d,%d", &l[0].x, &l[0].y, &l[1].x, &l[1].y); n != 4 || err != nil {
		return l, err
	}
	return l, nil
}

func (l Line) Diagonal() bool {
	return !l.horizontal() && !l.vertical()
}

func (l Line) horizontal() bool {
	return l[0].x == l[1].x
}

func (l Line) vertical() bool {
	return l[0].y == l[1].y
}

func (l Line) Points() []Point {
	var points []Point
	if l.horizontal() {
		points = make([]Point, 0, int(math.Abs(float64(l[0].y-l[1].y))))

		x := l[0].x

		for y := int(math.Min(float64(l[0].y), float64(l[1].y))); y <= int(math.Max(float64(l[0].y), float64(l[1].y))); y++ {
			points = append(points, Point{x, y})
		}

	} else if l.vertical() {
		points = make([]Point, 0, int(math.Abs(float64(l[0].x-l[1].x))))

		y := l[0].y

		for x := int(math.Min(float64(l[0].x), float64(l[1].x))); x <= int(math.Max(float64(l[0].x), float64(l[1].x))); x++ {
			points = append(points, Point{x, y})
		}
	} else {
		// "or a diagonal line at exactly 45 degrees"
		// so we can assume that any diagonal is going to be defined as:
		// bounded(x+1, y+1) / bounded(x+1, y-1) / bounded(x-1, y+1) / bounded(x-1, y-1)

		// writing this as not 4 separate statements would be great
		// we could reduce it to two by re-ordering x and y
		// having the case where they both increase, or one increases
		first, second := l[0], l[1]
		if first.x > second.x {
			first, second = second, first
		}

		x, y := first.x, first.y
		// okay, so first < second now for the x-axis, so we only now do an if for the y axis
		if first.y < second.y {
			// both x and y are increasing
			for x <= second.x && y <= second.y {
				points = append(points, Point{x, y})

				x, y = x+1, y+1
			}
		} else {
			// x is increasing, but y is decreasing
			for x <= second.x && y >= second.y {
				points = append(points, Point{x, y})

				x, y = x+1, y-1
			}
		}
	}

	return points
}

// A grid would more normally be described
// as a slice of points if we knew the limit for x
// or as a slice of slice of points if we didn't
// but in this case, we only care about points that are _touched_
// so why do all that storage if we don't need to?
// let's just refer to it as grid[x][y] instead
type Grid map[int]map[int]int

func newGrid() Grid {
	return map[int]map[int]int{}
}

func (g Grid) Draw(l Line) {
	for _, p := range l.Points() {
		if g[p.x] == nil {
			g[p.x] = map[int]int{}
		}
		g[p.x][p.y]++
	}
}

func (g Grid) PointsWithWatermark(watermark int) []Point {
	var points []Point
	for x := range g {
		for y := range g[x] {
			if g[x][y] >= watermark {
				points = append(points, Point{x, y})
			}
		}
	}

	return points
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

	// Part 1 grid is supposed to ignore diagonal lines
	// But part 2 does not
	// so we'll just track two separate grids
	// because keeping track of it within one grid would be difficult
	p1grid := newGrid()
	p2grid := newGrid()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line, err := parseLine(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		p2grid.Draw(line)

		// For part 1:
		// For now, only consider horizontal and vertical lines
		// A line that is neither horizontal nor vertical is diagonal
		if line.Diagonal() {
			continue
		}

		p1grid.Draw(line)

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Part 1: %d", len(p1grid.PointsWithWatermark(2)))
	log.Printf("Part 2: %d", len(p2grid.PointsWithWatermark(2)))
}

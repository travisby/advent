package main

import (
	"bufio"
	"container/heap"
	"log"
	"os"
	"strconv"
)

type point struct {
	x, y int
}

// a point does not decide what is a valid point
// only the density map will decide something like 0, -1 DNE
func (p point) adjacent() []point {
	return []point{
		// L
		point{p.x - 1, p.y},
		// R
		point{p.x + 1, p.y},
		// U
		point{p.x, p.y - 1},
		// D
		point{p.x, p.y + 1},
	}
}

type chitonDensityMap interface {
	Adjacent(p point) []point
	Get(p point) int
	Rows() int
	Columns() int
}

type sparseDensityMap struct {
	dm              denseChitonDensityMap
	sparsenessRatio int
}

func (c sparseDensityMap) Get(p point) int {
	if p.x < c.dm.Rows() && p.y < c.dm.Columns() {
		return c.dm.Get(p)
	}

	var result int
	if p.x >= c.dm.Rows() {
		p.x -= c.dm.Rows()
		result++
	}
	if p.y >= c.dm.Columns() {
		p.y -= c.dm.Columns()
		result++
	}

	result += c.Get(p)
	if result > 9 {
		result -= 9
	}
	return result
}

func (c sparseDensityMap) Rows() int {
	return c.sparsenessRatio * c.dm.Rows()
}
func (c sparseDensityMap) Columns() int {
	return c.sparsenessRatio * c.dm.Columns()
}

// womp-womp, I couldn't get type-embedding to cover Adjacent
// while not also not letting me override Rows/Columns/Get
func (c sparseDensityMap) Adjacent(p point) []point {
	points := p.adjacent()

	results := make([]point, 0, len(points))
	for _, poynt := range points {
		if poynt.x >= 0 && poynt.x < c.Rows() && poynt.y >= 0 && poynt.y < c.Columns() {
			results = append(results, poynt)
		}
	}

	return results
}

// a denseChitonDensityMap
// is exactly what you get
// a density map that's just a grid of risk levels
// there's no spare data here
type denseChitonDensityMap [][]int

// the density map can decide what is a valid point
func (c denseChitonDensityMap) Adjacent(p point) []point {
	points := p.adjacent()

	results := make([]point, 0, len(points))
	for _, poynt := range points {
		if poynt.x >= 0 && poynt.x < c.Rows() && poynt.y >= 0 && poynt.y < c.Columns() {
			results = append(results, poynt)
		}
	}

	return results
}

func (c denseChitonDensityMap) Get(p point) int {
	return c[p.x][p.y]
}
func (c denseChitonDensityMap) Rows() int {
	return len(c)
}

// XXX: Assumes uniform rows
func (c denseChitonDensityMap) Columns() int {
	return len(c[0])
}

func lowestTotalRisk(c chitonDensityMap, source point) int {
	// https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm#Using_a_priority_queue
	dist := make(map[point]int)

	dist[source] = 0

	var Q PriorityQueue[point]
	heap.Init(&Q)

	for x := 0; x < c.Rows(); x++ {
		for y := 0; y < c.Columns(); y++ {
			v := point{x, y}
			if v != source {
				dist[v] = 1<<32 - 1
			}
			heap.Push(&Q, &Item[point]{value: point{x, y}, priority: dist[v]})

		}
	}

	for Q.Len() > 0 {
		// XXX panic'able
		u := heap.Pop(&Q).(*Item[point])
		for _, v := range c.Adjacent(u.value) {
			if !Q.contains(v) {
				continue
			}

			alt := dist[u.value] + c.Get(v)
			if alt < dist[v] {
				dist[v] = alt
				// XXX panic'able
				Q.updatePriority(v, alt)
			}
		}
	}

	return dist[point{c.Rows() - 1, c.Columns() - 1}]
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

	chitons := make(denseChitonDensityMap, 0)
	for scanner.Scan() {
		risks := make([]int, 0, len(scanner.Text()))
		for _, c := range scanner.Text() {
			i, err := strconv.Atoi(string(c))
			if err != nil {
				log.Fatal(err)
			}
			risks = append(risks, i)
		}
		chitons = append(chitons, risks)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Part 1: %d", lowestTotalRisk(chitons, point{0, 0}))
	log.Printf("Part 2: %d", lowestTotalRisk(sparseDensityMap{chitons, 5}, point{0, 0}))

}

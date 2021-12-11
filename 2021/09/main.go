package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

type point struct {
	x, y int
}

func deduplicate(ps []point) []point {
	tempM := map[point]bool{}
	for _, p := range ps {
		tempM[p] = true
	}

	temp := make([]point, 0, len(tempM))
	for p := range tempM {
		temp = append(temp, p)
	}
	return temp
}

func (p point) left() point {
	return point{p.x - 1, p.y}
}

func (p point) right() point {
	return point{p.x + 1, p.y}
}

func (p point) up() point {
	return point{p.x, p.y + 1}
}

func (p point) down() point {
	return point{p.x, p.y - 1}
}

type heightmap [][]uint8

func (h heightmap) heightAt(p point) uint8 {
	return h[p.x][p.y]
}

func (h heightmap) rows() int {
	return len(h)
}

// xxx: assumes uniform rows
func (h heightmap) columns() int {
	if len(h) == 0 {
		return 0
	}
	return len(h[0])
}

func (h heightmap) adjacentPoints(p point) []point {
	var results []point

	if p.x != 0 && h.rows() > 0 {
		results = append(results, p.left())
	}
	if p.x < h.rows()-1 && h.rows() > 0 {
		results = append(results, p.right())
	}
	if p.y != 0 && h.columns() > 0 {
		results = append(results, p.down())
	}
	if p.y < h.columns()-1 && h.columns() > 0 {
		results = append(results, p.up())
	}

	return results
}

func (h heightmap) adjacentHeights(p point) []uint8 {
	results := []uint8{}
	for _, adjacent := range h.adjacentPoints(p) {
		results = append(results, h.heightAt(adjacent))
	}

	return results
}

func (h heightmap) isLowPoint(p point) bool {
	return h.riskLevel(p) != 0
}

func (h heightmap) riskLevel(p point) uint8 {
	ourHeight := h.heightAt(p)
	for _, adjacentHeight := range h.adjacentHeights(p) {
		if ourHeight >= adjacentHeight {
			return 0
		}
	}
	return ourHeight + 1
}

func (h heightmap) basin(p point) []point {
	var results []point
	if h.heightAt(p) == 9 {
		return results
	}
	results = append(results, p)

	for _, adjacent := range h.adjacentPoints(p) {
		// strictly less-than, or else we'd have a loop
		if h.heightAt(p) < h.heightAt(adjacent) {
			results = append(results, h.basin(adjacent)...)
		}
	}

	return deduplicate(results)
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

	var hm heightmap
	for scanner.Scan() {
		row := make([]uint8, len(scanner.Text()))
		for i, c := range scanner.Text() {
			num, err := strconv.ParseUint(fmt.Sprintf("%c", c), 10, 8)
			if err != nil {
				log.Fatal(err)
			}
			row[i] = uint8(num)
		}

		hm = append(hm, row)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var sumRiskLevels int
	var basinSizes []int
	for x := 0; x < hm.rows(); x++ {
		for y := 0; y < hm.columns(); y++ {
			p := point{x, y}
			sumRiskLevels += int(hm.riskLevel(p))

			if hm.isLowPoint(p) {
				basinSizes = append(basinSizes, len(hm.basin(p)))
			}
		}
	}

	log.Printf("Part 1: %d", sumRiskLevels)

	if len(basinSizes) < 3 {
		log.Fatal("Expected at least 3 basins")
	}

	sort.Sort(sort.Reverse(sort.IntSlice(basinSizes)))
	log.Printf("Part 2: %d", basinSizes[0]*basinSizes[1]*basinSizes[2])
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

// a grid is a x*y matrix with each cell containing the slice of claimIDs that are using it
type grid [][][]int

func newGrid(x, y int) *grid {
	g := make([][][]int, x)
	for i := range g {
		g[i] = make([][]int, y)
		for j := range g[i] {
			g[i][j] = []int{}
		}
	}

	temp := grid(g)
	return &temp
}

func (g grid) apply(c *claim) error {
	if c == nil {
		return fmt.Errorf("Did not expect claim to be nil")
	}
	if len(g) <= c.point.x+c.width || len(g[0]) <= c.point.y+c.height {
		return fmt.Errorf("Grid is not big enough to house %+v", *c)
	}

	for i := c.point.x; i < (c.point.x + c.width); i++ {
		for j := c.point.y; j < (c.point.y + c.height); j++ {
			g[i][j] = append(g[i][j], c.id)
		}
	}
	return nil
}

func (g grid) numContested() int {
	return len(g.contested())
}

func (g grid) contested() [][]int {
	contested := [][]int{}
	for i := range g {
		for j := range g[i] {
			if len(g[i][j]) > 1 {
				contested = append(contested, g[i][j])
			}
		}
	}
	return contested
}

func (g grid) contestingIDs() []int {
	idsDeduplicated := map[int]byte{}
	for _, cs := range g.contested() {
		for _, c := range cs {
			idsDeduplicated[c] = 0x00
		}
	}

	ids := make([]int, 0, len(idsDeduplicated))
	for k := range idsDeduplicated {
		ids = append(ids, k)
	}
	return ids
}

type claim struct {
	id    int
	point struct {
		x int
		y int
	}
	width  int
	height int
}

func newClaim(str string) (*claim, error) {
	var c claim
	// #123 @ 3,2: 5x4
	if n, err := fmt.Sscanf(str, "#%d @ %d,%d: %dx%d", &c.id, &c.point.x, &c.point.y, &c.width, &c.height); err != nil {
		return nil, err
	} else if n != 5 {
		return nil, fmt.Errorf("Expected to unmarshal 5 things, got %d", n)
	}
	return &c, nil
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

	claims := []claim{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		c, err := newClaim(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		claims = append(claims, *c)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	g := newGrid(1000, 1000)
	for _, c := range claims {
		if err := g.apply(&c); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("p1: %d\n", g.numContested())

	contested := g.contestingIDs()
	// make it easy to search!
	sort.Ints(contested)
	var id *int
	for _, c := range claims {
		// sort.SearchInts is weird, it will return the place the ID _should_ go if it were in the list
		// so check if the item there is actually our ID (or if its' the end of the list, don't try to panic)
		// if it's not in there already, we found our int
		potentialPlace := sort.SearchInts(contested, c.id)
		if len(contested) < potentialPlace {
			id = &c.id
			break
		} else if contested[potentialPlace] != c.id {
			id = &c.id
			break
		}
	}
	if id == nil {
		log.Fatal("Did not find a claim without contention")
	}
	fmt.Printf("p2: %d\n", *id)

}

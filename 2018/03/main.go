package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// a grid is a x*y matrix with each cell containing the slice of claims that are using it
type grid [][][]*claim

func newGrid(x, y int) *grid {
	g := make([][][]*claim, x)
	for i := range g {
		g[i] = make([][]*claim, y)
		for j := range g[i] {
			g[i][j] = []*claim{}
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
			g[i][j] = append(g[i][j], c)
		}
	}
	return nil
}

func (g grid) numContested() int {
	n := 0
	for i := range g {
		for j := range g[i] {
			if len(g[i][j]) > 1 {
				n++
			}
		}
	}
	return n
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

}

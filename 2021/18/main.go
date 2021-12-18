package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

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
	var p Pair
	for scanner.Scan() {
		newPair, err := strToPair(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		p = Add(p, *newPair)

		for p.Reducable() {
			p = p.Reduce()
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Part 1: %d", p.Magnitude())
}

func strToPair(str string) (*Pair, error) {
	var p Pair
	for _, c := range str {
		switch c {
		case '[':
			p = append(p, OpenBrace{})
		case ']':
			p = append(p, CloseBrace{})
		case ',':
			p = append(p, Comma{})
		default:
			// IF the last item in p is also a number
			// then we need to make ourselves into its least-sig digit
			num, err := strconv.Atoi(fmt.Sprintf("%c", c))
			if err != nil {
				return nil, err
			}

			switch i := p[len(p)-1].(type) {
			case OpenBrace, CloseBrace, Comma:
				p = append(p, Number(num))
			case Number:
				p[len(p)-1] = Number(int(i)*10 + num)
			}
		}
	}
	return &p, nil
}

type OpenBrace struct{}
type CloseBrace struct{}
type Number int
type Comma struct{}

type Element interface {
	Elementer()
}

func (OpenBrace) Elementer()      {}
func (OpenBrace) String() string  { return "[" }
func (CloseBrace) Elementer()     {}
func (CloseBrace) String() string { return "]" }
func (Number) Elementer()         {}
func (i Number) String() string   { return fmt.Sprintf("%d", i) }
func (Comma) Elementer()          {}
func (Comma) String() string      { return "," }

// When performing explosions
// a pair doesn't actually care about nesting
// e.g. the "left most element" doesn't have to be in the same pair
// so we represent pairs not as a nested structure
// but basically as a list of strings (turned into funner types)
type Pair []Element

func (p Pair) Reducable() bool {
	return p.indexOfExplodablePair() != nil || p.indexOfSplittable() != nil
}

func (p Pair) Reduce() Pair {
	if p.indexOfExplodablePair() != nil {
		return p.Explode()
	}
	return p.Split()
}

func (p Pair) Explode() Pair {
	// XXX: Although this returns (and it must be used!) a Pair
	// we actually are modifying the underlying Pair as well
	// we just have to return a new Pair to give the new len() that gets modified
	// we could do more index math with append to get this to not modify the underlying pair
	// but too lazy for that
	idx := p.indexOfExplodablePair()
	if idx == nil {
		return p
	}

	// p[*idx:*idx+5] = "[Num,Num]"
	// XXX panic'able
	leftPair := p[*idx+1].(Number)
	// +1 would be `,`
	// XXX panic'able
	rightPair := p[*idx+3].(Number)

	// search for the closest number to the left
	for i := *idx; i >= 0; i-- {
		num, ok := p[i].(Number)
		if ok {
			// and if we found it, increase it by our numbers
			p[i] = Number(num + leftPair)
			break
		}
	}
	// search for the closest number to the right
	for i := *idx + 4; i < len(p); i++ {
		num, ok := p[i].(Number)
		if ok {
			// and if we found it, increase it by our numbers
			p[i] = Number(num + rightPair)
			break
		}
	}

	p = append(
		p[:*idx],
		append(
			[]Element{Number(0)},
			// because "Exploding pairs will always consist of two regular numbers"
			// we can assume that an exploding pair is exactly four elements:
			p[*idx+5:]...,
		)...,
	)

	return p

}

func (p Pair) indexOfExplodablePair() *int {
	depth := 0
	for i, c := range p {
		switch c.(type) {
		case OpenBrace:
			depth++
		case CloseBrace:
			depth--
		case Number:
		}
		// we're looking for a _pair_ that is nested inside four pairs
		// meaning we have a depth of 5
		if depth > 4 {
			return &i
		}
	}
	return nil
}

func (p Pair) indexOfSplittable() *int {
	for i, c := range p {
		switch num := c.(type) {
		case Number:
			if num >= 10 {
				return &i
			}
		}
	}
	return nil
}

func (p Pair) Split() Pair {
	idx := p.indexOfSplittable()
	if idx == nil {
		return p
	}
	// XXX: panic'able
	num := int(p[*idx].(Number))

	return append(
		p[:*idx],
		append(
			[]Element{OpenBrace{}, Number(num / 2), Comma{}, Number(num/2 + num%2), CloseBrace{}},
			p[*idx+1:]...,
		)...,
	)
}

func (p Pair) Magnitude() int {
	// to calculate the Magnitude,
	// we now actually care about nesting
	// so we'll actually take our pair string
	// and let a json parser turn it into a nested struct for us
	var i interface{}

	if err := json.Unmarshal([]byte(p.String()), &i); err != nil {
		// XXX panic'able
		panic(err)
	}

	var magnitude func(i interface{}) int
	magnitude = func(i interface{}) int {
		switch j := i.(type) {
		case float64:
			// even though we're only dealing with whole numbers
			// json.Unmarshal -> interface uses float64
			return int(j)
		case []interface{}:
			if len(j) != 2 {
				panic(fmt.Errorf("Wrong num elements"))
			}
			return 3*magnitude(j[0]) + 2*magnitude(j[1])
		}
		return 0
	}

	return magnitude(i)
}

func (p Pair) String() string {
	strs := make([]string, 0, len(p))
	for _, a := range p {
		strs = append(strs, fmt.Sprintf("%s", a))
	}

	return strings.Join(strs, "")
}

func Add(p1, p2 Pair) Pair {
	// invalid pairs get silently replaced when Add'd to
	// this allows us to do:
	// var pair
	// for text {
	//    pair.Add(newPair)
	// }
	// with the input '[0, 0]'
	// and to get back '[0, 0]'
	// instead of '[[0, 0]]'
	if p1 == nil {
		return p2
	} else if p2 == nil {
		return p1
	}

	return append(
		[]Element{OpenBrace{}},
		append(
			p1,
			append(
				[]Element{Comma{}},
				append(
					p2,
					CloseBrace{},
				)...,
			)...,
		)...,
	)
}

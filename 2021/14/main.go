package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type element byte

type pair struct {
	left, right element
}

// maps each type of pair to the number of that pair we have
type polymer map[pair]int

func NewPolymer(s string) polymer {
	p := make(map[pair]int)

	for i := 0; i < len(s)-1; i++ {
		p[pair{element(s[i]), element(s[i+1])}]++
	}

	return p

}

// we don't want to modify our current polymer
// because we need to apply these rules "atomically"
// so we return a brand new one instead
func (p polymer) Apply(rs []rule) polymer {
	newP := p.copy()

	for _, r := range rs {
		if p[r.matcher] == 0 {
			continue
		}

		newP[r.matcher] -= p[r.matcher]
		for _, pare := range r.newPairs() {
			newP[pare] += p[r.matcher]
		}
	}

	return newP
}

func (p polymer) copy() polymer {
	temp := make(map[pair]int, len(p))
	for k, v := range p {
		temp[k] = v
	}
	return temp
}

func (p polymer) pairCountsToElementcounts() map[element]int {
	result := make(map[element]int)
	for k, v := range p {
		result[k.left] += v
		result[k.right] += v
	}

	// because we are counting pairs, we need to 1/2
	// every character we've seen
	for k, _ := range result {
		// the end-caps (firt and last character)
		// will likely (definitely?) appear an odd-number of times
		// and integer division will screw them out of one of their matches!
		if result[k]%2 != 0 {
			result[k]++
		}
		result[k] /= 2
	}

	return result
}

type rule struct {
	matcher   pair
	inBetween element
}

func NewRule(s string) (*rule, error) {
	var matcher, inBetween string
	if n, err := fmt.Sscanf(s, "%s -> %s", &matcher, &inBetween); n != 2 || err != nil {
		return nil, err
	} else if len(matcher) != 2 && len(inBetween) != 1 {
		return nil, fmt.Errorf("Invalid rule string: %q", s)
	}

	return &rule{pair{element(matcher[0]), element(matcher[1])}, element(inBetween[0])}, nil
}

func (r rule) newPairs() []pair {
	return []pair{pair{r.matcher.left, r.inBetween}, pair{r.inBetween, r.matcher.right}}
}

func elementCountsToAnswer(m map[element]int) int {
	var highest, lowest element
	for k, v := range m {
		if v > m[highest] {
			highest = k
		}
		if v < m[lowest] || m[lowest] == 0 {
			lowest = k
		}
	}
	return m[highest] - m[lowest]
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

	if !scanner.Scan() {
		log.Fatal("Expected polymer template")
	} else if err := scanner.Err(); err != nil {
		log.Fatal(scanner.Err())
	}

	polymerTemplate := scanner.Text()
	polymer := NewPolymer(polymerTemplate)

	if !scanner.Scan() {
		log.Fatal("Expected more data")
	} else if err := scanner.Err(); err != nil {
		log.Fatal(scanner.Err())
	} else if txt := scanner.Text(); txt != "" {
		log.Fatalf("Expected a newline before pair insertion rules, got %q", txt)
	}

	var rules []rule
	for scanner.Scan() {
		r, err := NewRule(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		rules = append(rules, *r)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	i := 0
	for ; i < 10; i++ {
		polymer = polymer.Apply(rules)
	}

	log.Printf("Part 1: %d", elementCountsToAnswer(polymer.pairCountsToElementcounts()))

	for ; i < 40; i++ {
		polymer = polymer.Apply(rules)
	}

	log.Printf("Part 2: %d", elementCountsToAnswer(polymer.pairCountsToElementcounts()))
}

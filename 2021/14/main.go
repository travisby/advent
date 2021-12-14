package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

type pairInsertionRule struct {
	element string
	left    string
	right   string
}

func (r pairInsertionRule) Match(left string, right string) bool {
	return r.left == left && r.right == right
}

func NewPairInsertionRule(s string) (*pairInsertionRule, error) {
	var leftright string
	var r pairInsertionRule
	if n, err := fmt.Sscanf(s, "%s -> %s", &leftright, &r.element); n != 2 || err != nil {
		return nil, err
	} else if len(leftright) != 2 {
		return nil, fmt.Errorf("Expected left side of rule to be two characters, got %q", leftright)
	}

	r.left = string(leftright[0])
	r.right = string(leftright[1])

	return &r, nil
}

type polymer struct {
	firstElement string
	next         *polymer
}

func NewPolymer(s string) *polymer {
	var head polymer
	// the head is going to be empty
	// just to make the iteration easier here
	// we'll just remember to return head.next instead of head!
	p := &head
	for _, c := range s {
		p.next = &polymer{firstElement: fmt.Sprintf("%c", c)}
		p = p.next
	}

	return head.next
}

func (p *polymer) Insert(e string) {
	p.next = &polymer{e, p.next}
}

func (p polymer) String() string {
	// XXX: This is recursive and is going to allocate lots of strings
	// we could improve this by not doing that
	// and doing a string builder like technique
	if p.next == nil {
		return p.firstElement
	}
	return p.firstElement + p.next.String()
}

func (p *polymer) ApplyRules(rules []pairInsertionRule) {
	if p == nil || p.next == nil {
		return
	}

	// by defering now we are also "locking in"
	// p.next
	// if any further statements change what p.next is
	// this defer won't care!
	// and will use the value as it was at this point
	defer p.next.ApplyRules(rules)

	for _, r := range rules {
		if r.Match(p.firstElement, p.next.firstElement) {
			p.Insert(r.element)
			// XXX: Assumes only one rule can match a particular pair
			break
		}
	}
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

	// first we expect the template
	if !scanner.Scan() {
		log.Fatal("No polymer template")
	} else if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	polymerTemplate := scanner.Text()

	// next we expect a newline
	if !scanner.Scan() {
		log.Fatal("No polymer template")
	} else if err := scanner.Err(); err != nil {
		log.Fatal(err)
	} else if scanner.Text() != "" {
		log.Fatal("Unexpected input!  Expected newline, got %q", scanner.Text())
	}

	// finally we get many pair insertion rules
	var pairInsertionRules []pairInsertionRule
	for scanner.Scan() {
		r, err := NewPairInsertionRule(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		pairInsertionRules = append(pairInsertionRules, *r)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	polymer := NewPolymer(polymerTemplate)

	i := 0
	for ; i < 10; i++ {
		polymer.ApplyRules(pairInsertionRules)
	}

	elementsToOccurrences := map[string]uint{}
	for p := polymer; p != nil; p = p.next {
		elementsToOccurrences[p.firstElement]++
	}

	var leastOccurring, mostOccurring string

	for e, c := range elementsToOccurrences {
		if c > elementsToOccurrences[mostOccurring] {
			mostOccurring = e
		}
		if c < elementsToOccurrences[leastOccurring] || elementsToOccurrences[leastOccurring] == 0 {
			leastOccurring = e
		}
	}

	log.Printf("Part 1: %d", elementsToOccurrences[mostOccurring]-elementsToOccurrences[leastOccurring])

	/*
		Stack overflows :cry:

		for ; i < 40; i++ {
			polymer.ApplyRules(pairInsertionRules)
		}

		elementsToOccurrences = map[string]uint{}
		for p := polymer; p != nil; p = p.next {
			elementsToOccurrences[p.firstElement]++
		}

		leastOccurring = ""
		mostOccurring = ""

		for e, c := range elementsToOccurrences {
			if c > elementsToOccurrences[mostOccurring] {
				mostOccurring = e
			}
			if c < elementsToOccurrences[leastOccurring] || elementsToOccurrences[leastOccurring] == 0 {
				leastOccurring = e
			}
		}

		log.Printf("Part 2: %d", elementsToOccurrences[mostOccurring]-elementsToOccurrences[leastOccurring])
	*/
}

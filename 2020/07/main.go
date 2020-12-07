package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
)

var BAG_LINE = regexp.MustCompile("^(.*) bags contain (.*)\\.$")
var BAG_CONTAINS = regexp.MustCompile("([0-9]+) ([a-z ]+) bags?")
var BAG_NO_CONTAINS = regexp.MustCompile("no other bags")

type bagType struct {
	color       string
	contains    map[*bagType]int
	containedBy []*bagType
}

func (b bagType) recursiveContainedBy() []*bagType {
	contains := b.containedBy
	for _, v := range b.containedBy {
		contains = append(contains, v.recursiveContainedBy()...)
	}

	// now uniq it
	contains2 := map[*bagType]struct{}{}
	for _, v := range contains {
		contains2[v] = struct{}{}
	}

	contains = nil
	for k := range contains2 {
		contains = append(contains, k)
	}

	return contains
}

func (b bagType) recursiveContains() []*bagType {
	contains := make([]*bagType, 0, len(b.contains))

	for v, k := range b.contains {
		vContains := v.recursiveContains()
		for i := 0; i < k; i++ {
			contains = append(contains, v)
			contains = append(contains, vContains...)
		}
	}

	return contains
}

type bagLookup map[string]*bagType

func (b bagLookup) get(color string) *bagType {
	return b[color]
}

func newBagLookup() bagLookup {
	return make(map[string]*bagType)
}

// contains is not a lookup, but adds a rule that parent contains child
func (b bagLookup) contains(parent string, count int, child string) {
	if _, ok := b[parent]; !ok {
		b[parent] = &bagType{parent, map[*bagType]int{}, nil}
	}

	if _, ok := b[child]; !ok {
		b[child] = &bagType{child, map[*bagType]int{}, nil}
	}

	b[parent].contains[b[child]] = count
	b[child].containedBy = append(b[child].containedBy, b[parent])
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

	bl := newBagLookup()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		bagLine := BAG_LINE.FindStringSubmatch(scanner.Text())
		if len(bagLine) != 3 {
			log.Fatalf("Bad parse %q", scanner.Text())
		}

		for _, containsLine := range BAG_CONTAINS.FindAllStringSubmatch(bagLine[2], -1) {
			if len(containsLine) != 3 {
				log.Fatalf("Bad parse in contains %q", scanner.Text())
			}

			c, err := strconv.Atoi(containsLine[1])
			if err != nil {
				log.Fatalf("Bad int parse in contains %q", scanner.Text())
			}

			bl.contains(bagLine[1], c, containsLine[2])
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Part 1: %d", len(bl.get("shiny gold").recursiveContainedBy()))
	log.Printf("Part 2: %d", len(bl.get("shiny gold").recursiveContains()))
}

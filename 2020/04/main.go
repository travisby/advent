package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type tv string

func (t tv) tagValue() (tag string, value *string) {
	split := strings.SplitN(string(t), ":", 2)
	if len(split) > 1 {
		return split[0], &split[1]
	}
	return split[0], nil
}

type passport []tv

func (p passport) hasTag(s string) bool {
	for _, t := range p {
		if k, _ := t.tagValue(); k == s {
			return true
		}
	}
	return false
}

func (p passport) valid() bool {
	return p.hasTag("byr") && p.hasTag("iyr") && p.hasTag("eyr") && p.hasTag("hgt") && p.hasTag("hcl") && p.hasTag("ecl") && p.hasTag("pid")
}

func (p *passport) addTV(t tv) {
	*p = append(*p, t)
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

	passports := []passport{}
	var p passport

	for scanner.Scan() {
		if scanner.Text() == "" {
			passports = append(passports, p)
			p = passport{}
		}

		for _, t := range strings.Split(scanner.Text(), " ") {
			p.addTV(tv(t))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var validCount int
	for _, p := range passports {
		if p.valid() {
			validCount += 1
		}
	}

	log.Printf("Part 1: %d", validCount)
}

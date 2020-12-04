package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var pidRe = regexp.MustCompile("^\\d{9}$")
var hclRe = regexp.MustCompile("^#[0-9a-f]{6}$")

type tv string

func (t tv) tagValue() (tag string, value *string) {
	split := strings.SplitN(string(t), ":", 2)
	if len(split) > 1 {
		return split[0], &split[1]
	}
	return split[0], nil
}

type passport []tv

func (p passport) hasTag(tag string) bool {
	_, ok := p.value(tag)
	return ok
}

func (p passport) value(tag string) (value *string, oko bool) {
	for _, t := range p {
		if k, v := t.tagValue(); k == tag {
			return v, true
		}
	}
	return nil, false

}

func (p passport) valid() bool {
	return p.hasTag("byr") && p.hasTag("iyr") && p.hasTag("eyr") && p.hasTag("hgt") && p.hasTag("hcl") && p.hasTag("ecl") && p.hasTag("pid")
}

func (p passport) valid2() bool {
	if byr, ok := p.value("byr"); !ok || byr == nil {
		return false
	} else if i, err := strconv.Atoi(*byr); err != nil || i < 1920 || i > 2002 {
		return false
	}

	if iyr, ok := p.value("iyr"); !ok || iyr == nil {
		return false
	} else if i, err := strconv.Atoi(*iyr); err != nil || i < 2010 || i > 2020 {
		return false
	}

	if eyr, ok := p.value("eyr"); !ok || eyr == nil {
		return false
	} else if i, err := strconv.Atoi(*eyr); err != nil || i < 2020 || i > 2030 {
		return false
	}

	if hgt, ok := p.value("hgt"); !ok || hgt == nil {
		return false
	} else if h, err := toHeight(*hgt); err != nil || (h.unit == "cm" && (h.n < 150 || h.n > 193)) || (h.unit == "in" && (h.n < 59 || h.n > 76)) {
		return false
	}

	if hcl, ok := p.value("hcl"); !ok || hcl == nil {
		return false
	} else if !hclRe.MatchString(*hcl) {
		return false
	}

	if ecl, ok := p.value("ecl"); !ok || ecl == nil {
		return false
	} else if *ecl != "amb" && *ecl != "blu" && *ecl != "brn" && *ecl != "gry" && *ecl != "grn" && *ecl != "hzl" && *ecl != "oth" {
		return false
	}

	if pid, ok := p.value("pid"); !ok || pid == nil {
		return false
	} else if !pidRe.MatchString(*pid) {
		return false
	}

	return true
}

func (p *passport) addTV(t tv) {
	*p = append(*p, t)
}

type height struct {
	unit string
	n    int
}

func toHeight(s string) (*height, error) {
	if len(s) < 3 {
		return nil, fmt.Errorf("invalid height")
	}

	unit := s[len(s)-2 : len(s)]
	if unit != "cm" && unit != "in" {
		return nil, fmt.Errorf("invalid height")
	}

	n, err := strconv.Atoi(s[:len(s)-2])
	if err != nil {
		return nil, err
	}
	return &height{unit, n}, nil
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

	validCount = 0
	for _, p := range passports {
		if p.valid2() {
			validCount += 1
		}
	}

	log.Printf("Part 2: %d", validCount)
}

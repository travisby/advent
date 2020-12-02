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

var PASSWORD_LINE_RE = regexp.MustCompile("^([0-9]+)-([0-9]+) ([a-zA-Z0-9]): ([a-zA-Z0-9]+)$")

type passwordLine struct {
	policy struct {
		r       string
		indices []int
	}
	password string
}

func (p passwordLine) isValid() bool {
	substr := make([]byte, len(p.policy.indices))
	for i, j := range p.policy.indices {
		if j > len(p.password) {
			return false
		}
		substr[i] = p.password[j-1]
	}

	return strings.Count(string(substr), p.policy.r) == 1
}

func newPasswordLine(s string) (*passwordLine, error) {
	pass := PASSWORD_LINE_RE.FindStringSubmatch(s)
	if len(pass) < 5 {
		return nil, fmt.Errorf("Unknown password format %q", s)
	}
	p := passwordLine{}

	// indices carrying password indices
	for _, i := range []int{1, 2} {
		j, err := strconv.Atoi(pass[i])
		if err != nil {
			return nil, fmt.Errorf("%w: cannot convert pass index: %q", err, pass[i])
		}
		p.policy.indices = append(p.policy.indices, j)
	}

	p.password = pass[len(pass)-1]
	p.policy.r = pass[len(pass)-2]

	return &p, nil
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

	var count int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if p, err := newPasswordLine(scanner.Text()); err != nil {
			log.Fatal(err)
		} else if p.isValid() {
			count++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d\n", count)
}

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
		r             string
		lowWaterMark  int
		highWaterMark int
	}
	password string
}

func (p passwordLine) isValid() bool {
	c := strings.Count(p.password, p.policy.r)
	return c >= p.policy.lowWaterMark && c <= p.policy.highWaterMark
}

func newPasswordLine(s string) (*passwordLine, error) {
	pass := PASSWORD_LINE_RE.FindStringSubmatch(s)
	if len(pass) != 5 {
		return nil, fmt.Errorf("Unknown password format %q", s)
	}
	p := passwordLine{password: pass[4]}
	p.policy.r = pass[3]

	var err error

	p.policy.lowWaterMark, err = strconv.Atoi(pass[1])
	if err != nil {
		return nil, fmt.Errorf("%w: while decoding low watermark %q", err, pass[0])
	}

	p.policy.highWaterMark, err = strconv.Atoi(pass[2])
	if err != nil {
		return nil, fmt.Errorf("%w: while decoding high watermark %q", err, pass[1])
	}

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

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

type period struct {
	start    time.Time
	duration time.Duration
}

type guardShift struct {
	guardID      int
	shift        period
	sleepPeriods []period
}

var errGuardShiftScan = errors.New("This scan does not seem to begin with a shift beginning.  Is the input ordered?")

func scanGuardShift(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// NOTE: this will break if len(rune) > 1b
	// we assume that str and []byte indices are interchangeable here
	// don't @ me

	// first few characters of data should look like:
	// [1518-02-12 23:50] Guard #1789 begins shift

	// okay, let's verify this is in fact a guard shift
	// but we can't do that if we don't have enough data
	// see if we have one whole line of input yet:
	// (and once we have our first line of input, we just look for "begins shift"
	// to verify if we're really a shift)
	firstNewline := bytes.Index(data, []byte{'\n'})
	if firstNewline == -1 {
		// nope, get more data until we at least do!
		return 0, nil, nil
	} else if !strings.Contains(string(data[:firstNewline]), "begins shift") {
		return 0, nil, errGuardShiftScan
	}

	// now, we basically keep requesting more data until either we're out of data (atEOF)
	// or we see the next "begins shift" (aka, any shift after firstNewLine)
	// tally-ho
	if atEOF {
		return len(data), data, nil
	} else if idx := strings.Index(string(data[firstNewline:]), "begins shift"); idx != -1 {
		// okay, we have a second shift somewhere around `idx`.  Walk back to the last newline
		for ; data[firstNewline+idx] != '\n' && idx > 0; idx-- {
		}
		return firstNewline + idx + 1, data[:firstNewline+idx+1], nil
	}

	// we are not done with this particular split because we don't have a second "begins shift"
	// yet, and we have determined there's more data to get (!atEOF)
	return 0, nil, nil
}

func readerToGuardShifts(r io.Reader) ([]guardShift, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(scanGuardShift)
	for scanner.Scan() {
		log.Printf("------")
		log.Printf(scanner.Text())
		log.Printf("------")
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return nil, fmt.Errorf("Not implemented")
}

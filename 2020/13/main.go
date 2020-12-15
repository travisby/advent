package main

import (
	"bufio"
	"log"
	"math"
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

	if !scanner.Scan() {
		log.Fatalf("Expected offset, got no first line input")
	} else if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	offset, err := strconv.Atoi(scanner.Text())
	if err != nil {
		log.Fatal(err)
	}

	if !scanner.Scan() {
		log.Fatalf("Expected buses, got no first line input")
	} else if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var buses []int
	for _, v := range strings.Split(scanner.Text(), ",") {
		if v == "x" {
			buses = append(buses, -1)
			continue
		}

		bus, err := strconv.Atoi(v)
		if err != nil {
			log.Fatal(err)
		}

		buses = append(buses, bus)
	}

	if scanner.Scan() {
		log.Fatalf("Expected no more input, y u do this?")
	}

	if len(buses) == 0 {
		log.Fatalf("No buses")
	}

	// XXX: assumes buses[0] is not -1!
	nextBus := buses[0]
	// https://math.stackexchange.com/questions/973057/find-smallest-number-bigger-than-y-that-is-multiple-of-x
	nextBusEarliestTime := int(math.Ceil(float64(offset)/float64(buses[0])) * float64(buses[0]))

	for _, bus := range buses[1:] {
		if bus == -1 {
			continue
		}

		busEarliestTime := int(math.Ceil(float64(offset)/float64(bus)) * float64(bus))
		if busEarliestTime < nextBusEarliestTime {
			nextBus = bus
			nextBusEarliestTime = busEarliestTime
		}
	}

	log.Printf("Part 1: %d", nextBus*(nextBusEarliestTime-offset))

	var bringEverythingToZero uint64 = 1
	for _, v := range buses {
		if v == -1 {
			continue
		}

		bringEverythingToZero *= uint64(v)
	}

	var totalSum uint64

	for i, v := range buses {
		if v == -1 {
			continue
		}

		// if we have the scenario where bus[37] == 23
		// then we might not want the FIRST bus to come at minute 37
		// but maybe the second bus comes at 37, e.g. the first bus might come
		// at minute 14
		// and then it would come again at minute 37
		for i > v {
			v += buses[i]
		}

		y := bringEverythingToZero / uint64(v)

		var count uint64 = 1
		for ; (y*count)%uint64(v) != uint64((v-i)%v); count++ {
		}

		totalSum += (y * count) % bringEverythingToZero
		totalSum = totalSum % bringEverythingToZero
	}

	log.Printf("Part 2: %d", totalSum%bringEverythingToZero)
}

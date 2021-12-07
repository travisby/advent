package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
)

func ScanWordsCommaSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		data = nil
	}
	if atEOF {
		return len(data), data, nil
	}

	var i int
	for i = range string(data) {
		if data[i] == ',' {
			break
		}
	}

	return i + 1, data[:i], nil
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
	scanner.Split(ScanWordsCommaSplit)

	positions := map[int64]int64{}
	var sum int64
	for scanner.Scan() {
		i, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		positions[i]++
		sum += i
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// int64 max
	lowestFuel := 9223372036854775807

	// xxx: we could probably do just a few movements
	// if we binary searched starting at the avg
	// but let's just be simple for now
	for i := 0; i <= 10000; i++ {
		fuel := fuelToMoveP1(positions, i)
		if fuel < lowestFuel {
			lowestFuel = fuel
		}
	}
	log.Printf("Part 1: %d", lowestFuel)

	// int64 max
	lowestFuel = 9223372036854775807

	// xxx: we could probably do just a few movements
	// if we binary searched starting at the avg
	// but let's just be simple for now
	for i := 0; i <= 10000; i++ {
		fuel := fuelToMoveP2(positions, i)
		if fuel < lowestFuel {
			lowestFuel = fuel
		}
	}
	log.Printf("Part 2: %d", lowestFuel)
}

func fuelToMoveP1(positions map[int64]int64, i int) (fuel int) {
	// The sum of the integers from 1 to n is n(n+1)/2
	for p, c := range positions {
		n := int(math.Abs(float64(int64(i) - p)))
		fuel += int(c) * n
	}
	return
}

func fuelToMoveP2(positions map[int64]int64, i int) (fuel int) {
	// The sum of the integers from 1 to n is n(n+1)/2
	for p, c := range positions {
		n := int(math.Abs(float64(int64(i) - p)))
		fuel += int(c) * (n * (n + 1) / 2)
	}
	return
}

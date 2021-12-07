package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

type school [9]uint

// xx: not really a string function
func (s school) String() string {
	var i uint
	for _, f := range s {
		i += f
	}
	return strconv.Itoa(int(i))
}

func (s *school) advance() {
	var temp school
	copy(temp[:], s[:])

	s[8] = temp[0]
	s[7] = temp[8]
	s[6] = temp[7] + temp[0]
	s[5] = temp[6]
	s[4] = temp[5]
	s[3] = temp[4]
	s[2] = temp[3]
	s[1] = temp[2]
	s[0] = temp[1]

}

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

	// fishAtDay stores how many fish are at $day in their lifecycle
	// this number is from 0-8, 8 being reserved for new feesh
	// 0 being the spawn cycle, and every number being reduced by 1
	// at the BEGINNING of each day
	// so the cycle is: decrease, spawn, do inventory (print), repeat
	fishAtDay := school{}

	for scanner.Scan() {
		u, err := strconv.ParseUint(scanner.Text(), 10, 64)
		if err != nil {
			log.Fatal(err)
		} else if u > 8 {
			log.Fatalf("Unknown fish lifecycle: %d", u)
		}
		fishAtDay[u]++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for i := 1; i <= 80; i++ {
		fishAtDay.advance()
	}
	log.Printf("Part 1: %s", fishAtDay)
	for i := 81; i <= 256; i++ {
		fishAtDay.advance()
	}
	log.Printf("Part 2: %s", fishAtDay)
}

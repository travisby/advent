package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type point struct{ x, y int }

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
		log.Fatal("")
	} else if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var begin, end point
	if n, err := fmt.Sscanf(scanner.Text(), "target area: x=%d..%d, y=%d..%d", &begin.x, &end.x, &begin.y, &end.y); n != 4 || err != nil {
		log.Fatal(err)
	}

	if scanner.Scan() {
		log.Fatal("")
	}

	log.Printf("target area: x=%d..%d, y=%d..%d", begin.x, end.x, begin.y, end.y)
}

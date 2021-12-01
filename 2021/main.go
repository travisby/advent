package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
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

	var oldI *int
	var increase int

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		i, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		if oldI != nil && i > *oldI {
			increase++
		}
		oldI = &i
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("%d", increase)
}

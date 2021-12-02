package main

import (
	"bufio"
	"container/ring"
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

	var increase int
	var slidingIncrease int

	// to have two separate 3-measurement windows
	// we only need 4 elements
	// A B C D -> [A B C] [B C D]
	intRing := ring.New(4)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		i, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		intRing.Value = &i

		// part a maths
		if intRing.Prev().Value != nil && *intRing.Value.(*int) > *intRing.Prev().Value.(*int) {
			increase++
		}

		// I got pretty stuck on doing this code here
		// I really would have liked Move() and Unlink() to have helped me
		// get "the next 3" or "the previous 3" and throw them into an summation function
		// unfortunately it would fail for some reason every time
		// the best I can think of is it does not like partially empty rings; the comments
		// for Unlink do say that r must not be nil, and assuming their implementation
		// at some point it is nil

		// part b maths, super panicy code
		// only do if the ring buffer is full!
		var notFull bool
		intRing.Do(func(i interface{}) {
			if i == nil {
				notFull = true
			}
		})

		if !notFull {
			prev := *intRing.Next().Value.(*int) + *intRing.Next().Next().Value.(*int) + *intRing.Next().Next().Next().Value.(*int)
			next := *intRing.Value.(*int) + *intRing.Prev().Value.(*int) + *intRing.Prev().Prev().Value.(*int)

			if next > prev {
				slidingIncrease++
			}
		}

		intRing = intRing.Next()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("part 1: %d", increase)
	log.Printf("part 2: %d", slidingIncrease)
}

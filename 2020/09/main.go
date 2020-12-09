package main

import (
	"bufio"
	"container/ring"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

const preambleSize = 25

var ErrInvalidNext = errors.New("The next number is not the sum of preceding numbers")

type decoderRing struct {
	*ring.Ring
}

func newDecoderRing(preamble []int) *decoderRing {
	r := ring.New(len(preamble))
	for _, v := range preamble {
		r = r.Next()
		r.Value = v
	}

	return &decoderRing{r}
}

func (d *decoderRing) next(i int) error {
	km, err := d.toKeyMap()
	if err != nil {
		return err
	}

	for k := range km {
		// if there exists the key that is i - our current number
		// (aka i-k + k = i), where (i-k) exists
		if _, ok := km[i-k]; ok {
			// valid!
			d.Ring = d.Next()
			d.Value = i
			return nil
		}
	}

	return fmt.Errorf("%w, %d", ErrInvalidNext, i)
}

var ErrWrongStuffing = errors.New("Someone stuffed something baaad in here")

func (d *decoderRing) toKeyMap() (map[int]struct{}, error) {
	s := make(map[int]struct{}, d.Len())
	var err error

	d.Do(func(i interface{}) {
		if err != nil {
			return
		}

		in, ok := i.(int)
		if !ok {
			err = fmt.Errorf("%w: %+v", ErrWrongStuffing, i)
			return
		}

		s[in] = struct{}{}
	})
	return s, err
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

	// first let's fill the preamble
	preamble := make([]int, 0, preambleSize)
	for len(preamble) < preambleSize && scanner.Scan() {
		i, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		preamble = append(preamble, i)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	} else if len(preamble) != preambleSize {
		log.Fatalf("Not enough content to consume: %+v", preamble)
	}

	decoder := newDecoderRing(preamble)
	var i *int
	// now continue scanning!
	for scanner.Scan() || i != nil {
		in, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		if err := decoder.next(in); err != nil {
			i = &in
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if i != nil {
		log.Printf("Part 1: %d", *i)
	} else {
		log.Fatal("Missing Part 1")
	}
}

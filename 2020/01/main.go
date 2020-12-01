package main

import (
	"bufio"
	"fmt"
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

	var is []int

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		i, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatalf("%w: Expected number, got %q", err, scanner.Text())
		}
		is = append(is, i)

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for i := range is {
		for j := range is {
			for k := range is {
				if is[i]+is[j]+is[k] == 2020 {
					fmt.Printf("%d\n", is[i]*is[j]*is[k])
					return
				}
			}
		}
	}
}

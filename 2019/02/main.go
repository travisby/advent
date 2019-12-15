package main

import (
	"bufio"
	"log"
	"os"
	"strconv"

	"gitlab.com/travisby/advent/2019/02/vm"
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

	memory := []int{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		maybeInt := scanner.Text()
		i, err := strconv.Atoi(maybeInt)
		if err != nil {
			log.Fatal(err)
		}
		memory = append(memory, i)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	virtualMachine := vm.New()
	if err := virtualMachine.Load(memory); err != nil {
		log.Fatal(err)
	} else if err := virtualMachine.Run(); err != nil {
		log.Fatal(err)
	}

	mem := virtualMachine.Memory[0]
	log.Printf("Value left at position 0: %d", mem)
}

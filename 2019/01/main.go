package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

type Module int // the mass of the Module

func MassFromString(s string) (*Module, error) {
	i, err := strconv.Atoi(s)

	m := Module(i)

	return &m, err
}

func (m Module) Fuel() int {
	fuel := int(m)/3 - 2
	if fuel <= 0 {
		return 0
	}
	return fuel
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

	totalFuel := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		module, err := MassFromString(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		totalFuel += module.Fuel()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Fuel required: %d", totalFuel)
}

package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"sort"
	"strconv"
)

var ErrIncompatibleInput = errors.New("The joltage difference is too great between the producer and consumer")

type joltageProducer interface {
	OutputJoltage() int
	Output() joltageConsumer
	SetOutput(joltageConsumer)
}
type joltageConsumer interface {
	Input() joltageProducer
	SetInput(joltageProducer) error
}

type adapter struct {
	outputJ int
	o       joltageConsumer
	i       joltageProducer
}

func (a *adapter) String() string {
	return strconv.Itoa(a.outputJ)
}

func (a *adapter) OutputJoltage() int          { return a.outputJ }
func (a *adapter) Output() joltageConsumer     { return a.o }
func (a *adapter) SetOutput(c joltageConsumer) { a.o = c }
func (a *adapter) Input() joltageProducer      { return a.i }
func (a *adapter) SetInput(p joltageProducer) error {
	if a.outputJ < p.OutputJoltage() || a.outputJ-3 > p.OutputJoltage() {
		return ErrIncompatibleInput
	}
	a.i = p
	return nil
}
func (a *adapter) Difference() int {
	return a.OutputJoltage() - a.Input().OutputJoltage()
}
func newAdapter(outputJoltage int) *adapter {
	return &adapter{outputJ: outputJoltage}
}

type seatOutlet struct {
	o joltageConsumer
}

// Treat the charging outlet near your seatOutlet as having an effective joltage rating of 0.
func (s *seatOutlet) OutputJoltage() int          { return 0 }
func (s *seatOutlet) Output() joltageConsumer     { return s.o }
func (s *seatOutlet) SetOutput(c joltageConsumer) { s.o = c }

type device struct {
	// this should be 3 higher than the highest adapter in our kit
	inputJ        int
	inputProducer joltageProducer
}

func (d *device) Input() joltageProducer { return d.inputProducer }
func (d *device) SetInput(c joltageProducer) error {
	// Your device has a built-in joltage adapter rated for 3 jolts higher than the highest-rated adapter in your bag.
	if c.OutputJoltage() != d.inputJ-3 {
		return ErrIncompatibleInput
	}

	d.inputProducer = c
	return nil
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

	var adapters []*adapter

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		i, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		adapters = append(adapters, newAdapter(i))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sort.SliceStable(adapters, func(i, j int) bool {
		return adapters[i].OutputJoltage() < adapters[j].OutputJoltage()
	})

	seat := &seatOutlet{}
	seat.SetOutput(adapters[0])
	if err := adapters[0].SetInput(seat); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(adapters)-1; i++ {
		adapters[i].SetOutput(adapters[i+1])
		if err := adapters[i+1].SetInput(adapters[i]); err != nil {
			log.Fatal(err)
		}
	}

	d := &device{inputJ: adapters[len(adapters)-1].OutputJoltage() + 3}
	d.SetInput(adapters[len(adapters)-1])
	adapters[len(adapters)-1].SetOutput(d)

	// 1 for the device not counted here
	joltageDifferences := map[int]int{1: 0, 3: 1}
	for _, a := range adapters {
		joltageDifferences[a.Difference()] = joltageDifferences[a.Difference()] + 1
	}
	log.Printf("%d", joltageDifferences[1]*joltageDifferences[3])

}

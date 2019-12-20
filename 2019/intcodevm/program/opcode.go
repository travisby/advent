package program

import (
	"errors"
	"fmt"
)

// ErrUnexpectedHalt is when we have an unexpected error that leads to a halt
var ErrUnexpectedHalt = errors.New("Unexpected err: HALT")

// ErrUnknownParameterMode is propagated up when we don't recognize the parameter digit
var ErrUnknownParameterMode = errors.New("Unknown parameter mode")

type parameter interface {
	Get(memory []int) (*int, error)
	Set(value int, memory []int) error
	String() string
}
type position struct {
	address int
}

func (p position) Get(memory []int) (*int, error) {
	if len(memory) > p.address {
		return &memory[p.address], nil
	}
	return nil, ErrUnexpectedHalt
}
func (p position) Set(value int, memory []int) error {
	if len(memory) > p.address {
		memory[p.address] = value
		return nil
	}
	return ErrUnexpectedHalt
}
func (p position) String() string {
	return fmt.Sprintf("$%d", p.address)
}

type immediate struct {
	value int
}

func (i immediate) Get(_ []int) (*int, error) {
	return &i.value, nil
}
func (i immediate) Set(value int, memory []int) error {
	return ErrUnexpectedHalt
}
func (i immediate) String() string {
	return fmt.Sprintf("%d", i.value)
}

type unknownParameterMode struct{}

func (u unknownParameterMode) Get(_ []int) (*int, error) {
	return nil, ErrUnknownParameterMode
}
func (u unknownParameterMode) Set(_ int, _ []int) error {
	return ErrUnknownParameterMode
}
func (u unknownParameterMode) String() string {
	return "{Unknown}"
}

func parameterMode(parameter int, mode int) parameter {
	switch mode {
	case 0:
		return position{parameter}
	case 1:
		return immediate{parameter}
	}
	return unknownParameterMode{}
}

type opcode int

const addOp opcode = 1
const multiplyOp opcode = 2
const haltOp opcode = 99
const inputOp opcode = 3
const outputOp opcode = 4

const lessThanOp opcode = 7
const equalsOp opcode = 8

func digitAt(n int, place int) int {
	return (n / place) % 10
}

package program

import "fmt"

// Opcode is the first integer in IntCode
// saying what action we will take
type Opcode int

// Add Opcode
const Add Opcode = 1

// Multiply Opcode
const Multiply Opcode = 2

// Halt Opcode
const Halt Opcode = 99

// Invalid is for broken / erronous cases
const Invalid Opcode = -1

// Intcode is an instruction in the Intcode language
type Intcode [4]int

// type Intcode interface {
// }
//
// type add [4]int
// type multiply [4]int
// type halt [4]int
//
// func Parse(instruction [4]int) (*Intcode, error) {
// 	switch Opcode(instruction[0]) {
// 	case Add:
// 		return &add(instruction), nil
// 	case Multiply:
// 		return &multiply(instruction), nil
// 	case Halt:
// 		return &halt(instruction), nil
// 	}
// 	return nil, fmt.Errorf("")
// }

func newIntcode(op, arg1, arg2, dest int) (*Intcode, error) {
	intcode := Intcode([4]int{op, arg1, arg2, dest})
	switch Opcode(op) {
	case Add, Multiply, Halt:
		return &intcode, nil
	}
	return nil, fmt.Errorf("")
}

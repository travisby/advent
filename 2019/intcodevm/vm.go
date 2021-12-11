package intcodevm

import (
	"errors"
	"io"
	"os"

	"gitlab.com/travisby/advent/2019/intcodevm/program"
)

var ErrOverflow = errors.New("The memory has overflowed")

// VM is our VirtualMachine that runs IntCode
type VM struct {
	memory   []int // the state of memory in the VM
	roMemory []int // the state of Load()'d data, ignoring what might happen after a Run().  This is a good copy of Programs
	in       io.Reader
	out      io.Writer
}

// New creates a new Virtual Machine
func New(memorySize int) *VM {
	return &VM{memory: make([]int, memorySize), roMemory: make([]int, memorySize), in: os.Stdin, out: os.Stdout}
}

// Load an intcode program into memory
func (v *VM) Load(offset int, ints []int) error {
	if len(ints)+offset > len(v.memory) || len(ints)+offset > len(v.roMemory) {
		return ErrOverflow
	}

	for i := range ints {
		v.memory[offset+i] = ints[i]
		v.roMemory[offset+i] = ints[i]
	}

	return nil
}

// Set the Noun for the loaded program
func (v *VM) SetNoun(noun int) error {
	if len(v.memory) < 2 || noun > 99 || noun < 0 {
		return ErrOverflow
	}
	// doesn't affect the ro memory
	v.memory[1] = noun
	return nil
}

// Set the Verb for the loaded program
func (v *VM) SetVerb(verb int) error {
	if len(v.memory) < 3 || verb > 99 || verb < 0 {
		return ErrOverflow
	}
	// doesn't affect the ro memory
	v.memory[2] = verb
	return nil
}

func (v *VM) SetIn(r io.Reader) {
	v.in = r
}

func (v *VM) SetOut(w io.Writer) {
	v.out = w
}

// Run the loaded program
func (v *VM) Run() error {
	p := program.NewScanner(v.memory, v.in, v.out)
	for p.Scan() {
		if err := p.Instruction().Apply(v.memory); err != nil {
			return err
		}

	}
	return p.Err()
}

// Loads the program back to its initial state
func (v *VM) Reset() error {
	if size := copy(v.memory, v.roMemory); size != len(v.roMemory) {
		return ErrOverflow
	}
	return nil
}

// Output contains the program's output
func (v *VM) Output() int {
	// output is really just address=0
	return v.memory[0]
}

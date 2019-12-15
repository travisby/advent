package vm

import (
	"errors"

	"gitlab.com/travisby/advent/2019/02/vm/program"
)

var ErrOverflow = errors.New("The memory has overflowed")

// VM is our VirtualMachine that runs IntCode
type VM struct {
	memory []int // the state of memory in the VM
}

// New creates a new Virtual Machine
func New(memorySize int) *VM {
	// align to 4 ints
	if memorySize%4 != 0 {
		memorySize = (memorySize/4 + 1) * 4
	}

	return &VM{memory: make([]int, memorySize)}
}

// Load an intcode program into memory
func (v *VM) Load(offset int, ints []int) error {
	if len(ints)+offset > len(v.memory) {
		return ErrOverflow
	}

	for i := range ints {
		v.memory[offset+i] = ints[i]
	}

	return nil
}

// Run the loaded program
func (v *VM) Run() error {
	p := program.NewScanner(v.memory)
	for p.Scan() {
		p.Intcode().Apply(v.memory)

	}
	return p.Err()
}

package vm

import (
	"errors"

	"gitlab.com/travisby/advent/2019/02/vm/program"
)

var ErrOverflow = errors.New("The memory has overflowed")

// VM is our VirtualMachine that runs IntCode
type VM struct {
	Memory []int // the state of memory in the VM
}

// New creates a new Virtual Machine
func New(memorySize int) *VM {
	// align to 4 ints
	if memorySize%4 != 0 {
		memorySize = (memorySize/4 + 1) * 4
	}

	return &VM{make([]int, memorySize)}
}

// Load an intcode program into memory
func (v *VM) Load(offset int, ints []int) error {
	if len(ints)+offset > len(v.Memory) {
		return ErrOverflow
	}

	for i := range ints {
		v.Memory[offset+i] = ints[i]
	}

	return nil
}

// Run the loaded program
func (v *VM) Run() error {
	p := program.NewScanner(v.Memory)
	for p.Scan() {
		p.Intcode().Apply(v.Memory)

	}
	return p.Err()
}

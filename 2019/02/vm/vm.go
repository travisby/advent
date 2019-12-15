package vm

import (
	"fmt"

	"gitlab.com/travisby/advent/2019/02/program"
)

// VM is our VirtualMachine that runs IntCode
type VM struct {
	Memory  []int // the state of memory in the VM
	Program program.Program
}

// New creates a new Virtual Machine
func New() *VM {
	return &VM{}
}

// Load an intcode program into memory
func (v *VM) Load([]int) error {
	return fmt.Errorf("Not implemented")
}

// Run the loaded program
func (v *VM) Run() error {
	return fmt.Errorf("Not implemented")
}

package program

import (
	"log"
	"testing"
)

func TestPositionJumpIfFalse(t *testing.T) {
	var ip int
	expectedIp := 99

	j := jumpFalse{position{0}, position{expectedIp}, &ip}

	if err := j.Apply([]int{0}); err != nil {
		log.Fatal(err)
	}

	if ip != expectedIp {
		log.Fatalf("Got instruction pointer (%d) expected %d", ip, expectedIp)
	}
}

func TestPositionNoJumpJumpIfFalse(t *testing.T) {
	var ip int
	expectedIp := ip

	j := jumpFalse{position{0}, position{99}, &ip}

	if err := j.Apply([]int{1}); err != nil {
		log.Fatal(err)
	}

	if ip != expectedIp {
		log.Fatalf("Got instruction pointer (%d) expected %d", ip, expectedIp)
	}
}

func TestImmediateJumpIfFalse(t *testing.T) {
	var ip int
	expectedIp := 99

	j := jumpFalse{immediate{0}, position{expectedIp}, &ip}

	if err := j.Apply([]int{}); err != nil {
		log.Fatal(err)
	}

	if ip != expectedIp {
		log.Fatalf("Got instruction pointer (%d) expected %d", ip, expectedIp)
	}
}

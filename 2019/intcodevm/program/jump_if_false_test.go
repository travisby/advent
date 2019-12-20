package program

import (
	"log"
	"testing"
)

func TestPositionJumpIfFalse(t *testing.T) {
	var ip int
	expectedIp := 0

	j := jumpFalse{position{0}, position{1}, &ip}

	if err := j.Apply([]int{0, expectedIp}); err != nil {
		log.Fatal(err)
	}

	if ip != expectedIp {
		log.Fatalf("Got instruction pointer (%d) expected %d", ip, expectedIp)
	}
}

func TestPositionNoJumpJumpIfFalse(t *testing.T) {
	var ip int
	expectedIp := ip

	j := jumpFalse{position{0}, position{1}, &ip}

	if err := j.Apply([]int{1, 99}); err != nil {
		log.Fatal(err)
	}

	if ip != expectedIp {
		log.Fatalf("Got instruction pointer (%d) expected %d", ip, expectedIp)
	}
}

func TestImmediateJumpIfFalse(t *testing.T) {
	var ip int
	expectedIp := 99

	j := jumpFalse{immediate{0}, immediate{expectedIp}, &ip}

	if err := j.Apply([]int{}); err != nil {
		log.Fatal(err)
	}

	if ip != expectedIp {
		log.Fatalf("Got instruction pointer (%d) expected %d", ip, expectedIp)
	}
}

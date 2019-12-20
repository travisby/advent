package program

import (
	"log"
	"testing"
)

func TestPositionJumpIfTrue(t *testing.T) {
	var ip int
	expectedIp := 99

	j := jumpTrue{position{0}, position{expectedIp}, &ip}

	if err := j.Apply([]int{1}); err != nil {
		log.Fatal(err)
	}

	if ip != expectedIp {
		log.Fatalf("Got instruction pointer (%d) expected %d", ip, expectedIp)
	}
}

func TestPositionNoJumpJumpIfTrue(t *testing.T) {
	var ip int
	expectedIp := ip

	j := jumpTrue{position{0}, position{99}, &ip}

	if err := j.Apply([]int{0}); err != nil {
		log.Fatal(err)
	}

	if ip != expectedIp {
		log.Fatalf("Got instruction pointer (%d) expected %d", ip, expectedIp)
	}
}

func TestImmediateJumpIfTrue(t *testing.T) {
	var ip int
	expectedIp := 99

	j := jumpTrue{immediate{10}, position{expectedIp}, &ip}

	if err := j.Apply([]int{1}); err != nil {
		log.Fatal(err)
	}

	if ip != expectedIp {
		log.Fatalf("Got instruction pointer (%d) expected %d", ip, expectedIp)
	}
}

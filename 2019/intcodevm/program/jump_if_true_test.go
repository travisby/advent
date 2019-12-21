package program

import (
	"log"
	"testing"
)

func TestPositionJumpIfTrue(t *testing.T) {
	var ip int
	expectedIp := 0

	j := jumpTrue{position{0}, position{1}, &ip}

	if err := j.Apply([]int{1, expectedIp}); err != nil {
		log.Fatal(err)
	}

	if ip != expectedIp {
		log.Fatalf("Got instruction pointer (%d) expected %d", ip, expectedIp)
	}
}

func TestPositionNoJumpJumpIfTrue(t *testing.T) {
	var ip int
	expectedIp := ip

	j := jumpTrue{position{0}, position{1}, &ip}

	if err := j.Apply([]int{0, 99}); err != nil {
		log.Fatal(err)
	}

	if ip != expectedIp {
		log.Fatalf("Got instruction pointer (%d) expected %d", ip, expectedIp)
	}
}

func TestImmediateJumpIfTrue(t *testing.T) {
	var ip int
	expectedIp := 99

	j := jumpTrue{immediate{1}, immediate{expectedIp}, &ip}

	if err := j.Apply([]int{}); err != nil {
		log.Fatal(err)
	}

	if ip != expectedIp {
		log.Fatalf("Got instruction pointer (%d) expected %d", ip, expectedIp)
	}
}

func TestMixedJumpIfTrue(t *testing.T) {
	var ip int
	expectedIp := 99

	j := jumpTrue{immediate{1}, position{0}, &ip}

	if err := j.Apply([]int{expectedIp}); err != nil {
		log.Fatal(err)
	}

	if ip != expectedIp {
		log.Fatalf("Got instruction pointer (%d) expected %d", ip, expectedIp)
	}
}

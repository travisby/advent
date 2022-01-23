package main

type threeDRotMatrix [3][3]int

func (m threeDRotMatrix) Apply(p Point) Point {
	return Point{
		p.x*m[0][0] + p.y*m[0][1] + p.z*m[0][2],
		p.x*m[1][0] + p.y*m[1][1] + p.z*m[1][2],
		p.x*m[2][0] + p.y*m[2][1] + p.z*m[2][2],
	}
}

func Rx(theta int) threeDRotMatrix {
	return [3][3]int{
		{1, 0, 0},
		{0, cos(theta), -sin(theta)},
		{0, sin(theta), cos(theta)},
	}
}

func Ry(theta int) threeDRotMatrix {
	return [3][3]int{
		{cos(theta), 0, sin(theta)},
		{0, 1, 0},
		{-sin(theta), 0, cos(theta)},
	}
}

func Rz(theta int) threeDRotMatrix {
	return [3][3]int{
		{cos(theta), -sin(theta), 0},
		{sin(theta), cos(theta), 0},
		{0, 0, 1},
	}
}

func cos(theta int) int {
	// we only return whole numbers
	// realistically really just 1 or 0
	// so panic if we're not a 90deg rot
	switch theta {
	case 0:
		return 1
	case 90, 270:
		return 0
	case 180:
		return -1
	}
	panic("Unknown cos")
}
func sin(theta int) int {
	// we only return whole numbers
	// realistically really just 1 or 0
	// so panic if we're not a 90deg rot
	switch theta {
	case 0, 180:
		return 0
	case 90:
		return 1
	case 270:
		return -1

	}
	panic("Unknown sin")
}

type Point struct {
	x, y, z int
}

func (p Point) Apply(r Rotation) Point {
	return Rz(r.z).Apply(Ry(r.y).Apply(Rx(r.x).Apply(p)))
}

type Rotation struct {
	x, y, z int
}

// returns the valid list of 3d rotations that make unique points
func getRotations() []Rotation {
	// first, we declare what are valid rotations PER axis
	// every 90deg turn is valid
	validRots := []int{0, 90, 180, 270}

	// next, we're going to -- for each axis, produce a permutation of these rotations
	permutations := make([]Rotation, 0, 4*4*4)
	for _, rotX := range validRots {
		for _, rotY := range validRots {
			for _, rotZ := range validRots {
				permutations = append(permutations, Rotation{rotX, rotY, rotZ})
			}
		}
	}

	// that's nice, that produced 64 results.  However only 24 are valid!
	// that's because some of those permutations end up producing identical points after rotation
	// so, to get the UNIQUE ones
	// we're going to apply all of these rotations to some arbitrary point
	// dedupe the points, and then return the transforms that got us to those points

	// some arbitrary point
	p := Point{10, 20, 30}

	// we can dedupe points by using them as the key in the map
	// and then we don't care WHICH of the 2-3 rotations that led to that point
	// as long as we have ONE of them
	points := make(map[Point]Rotation)
	for _, xyz := range permutations {
		points[p.Apply(xyz)] = xyz
	}

	// now just collect all of those transforms
	unDuplicatedRotations := make([]Rotation, 0, len(points))
	for _, v := range points {
		unDuplicatedRotations = append(unDuplicatedRotations, v)
	}

	return unDuplicatedRotations

}

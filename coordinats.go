package main

import "math"

type Coords struct {
	X, Y float64
}


func (c Coords) Add(rhs Coords) Coords {
	return Coords{
		X: c.X + rhs.X,
		Y: c.Y + rhs.Y,
	}
}


func (c Coords) Neg() Coords {
	return c.Mul(-1)
}


func (c Coords) Mul(a float64) Coords {
	return Coords{
		X: c.X * a,
		Y: c.Y * a,
	}
}


func (c Coords) Distance(to Coords) float64 {
	return math.Sqrt(math.Pow(c.X - to.X, 2) + math.Pow(c.Y - to.Y, 2))
}

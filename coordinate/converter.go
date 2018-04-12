package coordinate

import (
	"math"
)

const (
	BEPOCH_WEIGHT float64 = (15019.81352 + (1950.-1900.)*365.242198781 - 51544.5) / 365.25 / (100. * 60. * 60. * 360. / 2.)
)

func J2000ToGal(c Coordinate) Coordinate {
	RMAT := [][]float64{
		{-0.054875539726, -0.873437108010, -0.483834985808},
		{+0.494109453312, -0.444829589425, +0.746982251810},
		{-0.867666135858, -0.198076386122, +0.455983795705},
	}
	c1 := c.ToCartesian()
	c2 := &Cartesian{
		X: RMAT[0][0]*c1.X + RMAT[0][1]*c1.Y + RMAT[0][2]*c1.Z,
		Y: RMAT[1][0]*c1.X + RMAT[1][1]*c1.Y + RMAT[1][2]*c1.Z,
		Z: RMAT[2][0]*c1.X + RMAT[2][1]*c1.Y + RMAT[2][2]*c1.Z,
	}
	s := c2.ToSpherical().ToGal()
	return NewCoordinateFromSphere(`Gal`, s)
}

func GalToJ2000(c Coordinate) Coordinate {
	RMAT := [][]float64{
		{-0.054875539726, -0.873437108010, -0.483834985808},
		{+0.494109453312, -0.444829589425, +0.746982251810},
		{-0.867666135858, -0.198076386122, +0.455983795705},
	}
	c1 := c.ToCartesian()
	c2 := &Cartesian{
		X: RMAT[0][0]*c1.X + RMAT[1][0]*c1.Y + RMAT[2][0]*c1.Z,
		Y: RMAT[0][1]*c1.X + RMAT[1][1]*c1.Y + RMAT[2][1]*c1.Z,
		Z: RMAT[0][2]*c1.X + RMAT[1][2]*c1.Y + RMAT[2][2]*c1.Z,
	}
	s := c2.ToSpherical().ToEq()
	return NewCoordinateFromSphere(`J2000`, s)
}

func J2000ToB1950(c Coordinate) Coordinate {
	A := []float64{-1.62557e-6, -0.31919e-6, -0.13843e-6, 1.245e-3, -1.580e-3, -0.659e-3}
	EMI := [][]float64{
		{+0.9999256795, -0.0111814828, -0.0048590040, -0.000551, -0.238560, +0.435730},
		{+0.0111814828, +0.9999374849, -0.0000271557, +0.238509, -0.002667, -0.008541},
		{+0.0048590039, -0.0000271771, +0.9999881946, -0.435614, +0.012254, +0.002117},
	}

	c1 := c.ToCartesian()
	x := c1.X
	y := c1.Y
	z := c1.Z

	v1 := []float64{x, y, z}
	v2 := []float64{0, 0, 0, 0, 0, 0}

	for i := 0; i < 6; i++ {
		for j := 0; j < 3; j++ {
			v2[i] += EMI[j][i] * v1[j]
		}
	}

	rxyz := math.Sqrt(v2[0]*v2[0] + v2[1]*v2[1] + v2[2]*v2[2])

	w := v2[0]*A[0] + v2[1]*A[1] + v2[2]*A[2]
	x = (1.-w)*v2[0] + A[0]*rxyz
	y = (1.-w)*v2[1] + A[1]*rxyz
	z = (1.-w)*v2[2] + A[2]*rxyz

	rxyz = math.Sqrt(x*x + y*y + z*z)

	w = v2[0]*A[0] + v2[1]*A[1] + v2[2]*A[2]
	x = (1.-w)*v2[0] + A[0]*rxyz
	y = (1.-w)*v2[1] + A[1]*rxyz
	z = (1.-w)*v2[2] + A[2]*rxyz

	rxy := math.Sqrt(x*x + y*y)

	var ra, dec float64
	if x == 0 && y == 0 {
		ra = 0.
	} else {
		ra = math.Atan2(y, x)
		if ra < 0 {
			ra += math.Pi * 2.
		}
	}
	dec = math.Atan2(z, rxy)

	return NewCoordinate(`B1950`, RadToDeg(ra), RadToDeg(dec))
}

func B1950ToJ2000(c Coordinate) Coordinate {
	A := []float64{-1.62557e-6, -0.31919e-6, -0.13843e-6}
	EM := [][]float64{
		{+0.9999256782, +0.0111820610, +0.0048579479, -0.000551, +0.238514, -0.435623},
		{-0.0111820611, +0.9999374784, -0.0000271474, -0.238565, -0.002667, +0.012254},
		{-0.0048579477, -0.0000271765, +0.9999881997, +0.435739, -0.008541, +0.002117},
	}

	c1 := c.ToCartesian()
	r0 := []float64{c1.X, c1.Y, c1.Z}
	w := r0[0]*A[0] + r0[1]*A[1] + r0[2]*A[2]

	for i := 0; i < 3; i++ {
		r0[i] = (1.+w)*r0[i] - A[i]
	}

	v2 := []float64{0., 0., 0., 0., 0., 0.}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			v2[i] += EM[j][i] * r0[j]
			v2[i+3] += EM[j][i+3] * r0[j]
		}
		v2[i] += math.Pi * BEPOCH_WEIGHT * v2[i+3]
	}

	c2 := &Cartesian{X: v2[0], Y: v2[1], Z: v2[2]}
	s := c2.ToSpherical().ToEq()
	return NewCoordinateFromSphere(`J2000`, s)
}

func B1950ToGal(c Coordinate) Coordinate {
	return J2000ToGal(B1950ToJ2000(c))
}

func GalToB1950(c Coordinate) Coordinate {
	return J2000ToB1950(GalToJ2000(c))
}

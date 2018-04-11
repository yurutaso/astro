package coordinate

import (
	"math"
)

const (
	BEPOCH_WEIGHT float64 = (15019.81352 + (1950.-1900.)*365.242198781 - 51544.5) / 365.25 / (100. * 60. * 60. * 360. / 2.)
)

func DCS2C(h, v float64) []float64 {
	return []float64{
		math.Cos(h) * math.Cos(v),
		math.Sin(h) * math.Cos(v),
		math.Sin(v),
	}
}

func DCC2S(V []float64) (float64, float64) {
	x := V[0]
	y := V[1]
	z := V[2]
	r := math.Sqrt(x*x + y*y)
	var h float64
	var v float64
	if r == 0 {
		h = 0.0
	} else {
		h = math.Atan2(y, x)
	}

	if z == 0 {
		v = 0.0
	} else {
		v = math.Atan2(z, r)
	}
	return h, v
}

func J2000ToGal(c *J2000) *Gal {
	RMAT := [][]float64{
		{-0.054875539726, -0.873437108010, -0.483834985808},
		{+0.494109453312, -0.444829589425, +0.746982251810},
		{-0.867666135858, -0.198076386122, +0.455983795705},
	}
	v1 := DCS2C(c.Ra.Radian(), c.Dec.Radian())
	v2 := []float64{0, 0, 0}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			v2[i] += RMAT[i][j] * v1[j]
		}
	}
	h, v := DCC2S(v2)
	l := math.Mod(h, 2.*math.Pi)
	if l < 0 {
		l += 2. * math.Pi
	}
	b := math.Mod(v, 2.*math.Pi)
	if math.Abs(b) > math.Pi {
		if v > 0 {
			b -= math.Pi
		}
		if v < 0 {
			b += math.Pi
		}
	}
	return &Gal{
		Lon: &Degree{deg: l * 180. / math.Pi},
		Lat: &Degree{deg: b * 180. / math.Pi},
	}
}

func GalToJ2000(c *Gal) *J2000 {
	RMAT := [][]float64{
		{-0.054875539726, -0.873437108010, -0.483834985808},
		{+0.494109453312, -0.444829589425, +0.746982251810},
		{-0.867666135858, -0.198076386122, +0.455983795705},
	}
	v1 := DCS2C(c.Lon.Radian(), c.Lat.Radian())
	v2 := []float64{0, 0, 0}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			v2[i] += RMAT[j][i] * v1[j]
		}
	}
	h, v := DCC2S(v2)

	ra := math.Mod(h, 2.*math.Pi)
	if ra < 0 {
		ra += 2. * math.Pi
	}

	dec := math.Mod(v, 2.*math.Pi)
	if math.Abs(dec) > math.Pi {
		if v > 0 {
			dec -= math.Pi
		}
		if v < 0 {
			dec += math.Pi
		}
	}
	return &J2000{
		Ra:  &Degree{deg: ra * 180. / math.Pi},
		Dec: &Degree{deg: dec * 180. / math.Pi},
	}
}

func J2000ToB1950(c *J2000) *B1950 {
	A := []float64{-1.62557e-6, -0.31919e-6, -0.13843e-6, 1.245e-3, -1.580e-3, -0.659e-3}
	EMI := [][]float64{
		{+0.9999256795, -0.0111814828, -0.0048590040, -0.000551, -0.238560, +0.435730},
		{+0.0111814828, +0.9999374849, -0.0000271557, +0.238509, -0.002667, -0.008541},
		{+0.0048590039, -0.0000271771, +0.9999881946, -0.435614, +0.012254, +0.002117},
	}

	ra := c.Ra.Radian()
	dec := c.Dec.Radian()

	x := math.Cos(ra) * math.Cos(dec)
	y := math.Sin(ra) * math.Cos(dec)
	z := math.Sin(dec)

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

	if x == 0 && y == 0 {
		ra = 0.
	} else {
		ra = math.Atan2(y, x)
		if ra < 0 {
			ra += math.Pi * 2.
		}
	}
	dec = math.Atan2(z, rxy)

	return &B1950{
		Ra:  &Degree{deg: ra * 180. / math.Pi},
		Dec: &Degree{deg: dec * 180. / math.Pi},
	}
}

func B1950ToJ2000(c *B1950) *J2000 {
	A := []float64{-1.62557e-6, -0.31919e-6, -0.13843e-6}
	EM := [][]float64{
		{+0.9999256782, +0.0111820610, +0.0048579479, -0.000551, +0.238514, -0.435623},
		{-0.0111820611, +0.9999374784, -0.0000271474, -0.238565, -0.002667, +0.012254},
		{-0.0048579477, -0.0000271765, +0.9999881997, +0.435739, -0.008541, +0.002117},
	}

	r0 := DCS2C(c.Ra.Radian(), c.Dec.Radian())
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

	h, v := DCC2S(v2)
	ra := math.Mod(h, 2.*math.Pi)
	if ra < 0 {
		ra += 2. * math.Pi
	}
	dec := v

	return &J2000{
		Ra:  &Degree{deg: ra * 180. / math.Pi},
		Dec: &Degree{deg: dec * 180. / math.Pi},
	}
}

func B1950ToGal(c *B1950) *Gal {
	return J2000ToGal(B1950ToJ2000(c))
}

func GalToB1950(c *Gal) *B1950 {
	return J2000ToB1950(GalToJ2000(c))
}

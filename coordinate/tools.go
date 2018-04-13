package coordinate

import (
	"math"
)

func DegToRad(deg float64) float64 {
	return deg * math.Pi / 180.
}

func RadToDeg(rad float64) float64 {
	return rad * 180. / math.Pi
}

func Elevation(lst *Angle, c Coordinate, o Observatory) float64 {
	/* Get elevation of the source c from observatory o at LST time lst
	The returned value ranges from -90 deg to +90 deg.
	*/
	x := lst.Radian()
	l := o.Latitude().Radian()
	ra := c.GetX().Radian()
	dec := c.GetY().Radian()

	A := math.Sin(l) * math.Sin(dec)
	B := math.Cos(l) * math.Cos(dec)

	return 90. - 180./math.Pi*math.Acos(A+B*math.Cos(x-ra))
}

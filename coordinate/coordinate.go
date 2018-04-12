package coordinate

import (
	"fmt"
	"log"
	"math"
)

/* Convinient functions */
func DegToRad(deg float64) float64 {
	return deg * math.Pi / 180.
}

func RadToDeg(rad float64) float64 {
	return rad * 180. / math.Pi
}

/* Struct Angle */
type Angle struct {
	deg float64
}

func NewAngle(deg float64) *Angle {
	return &Angle{deg: deg}
}

func NewAngleFromDMS(deg, min, sec float64) *Angle {
	if deg < 0 {
		return &Angle{deg: -(-deg + min/60. + sec/3600.)}
	}
	return &Angle{deg: deg + min/60. + sec/3600.}
}

func NewAngleFromHMS(hour, min, sec float64) *Angle {
	ang := NewAngleFromDMS(hour, min, sec)
	ang.deg *= 15.
	return ang
}

func (ang *Angle) DMS() (float64, float64, float64) {
	sign := 1.
	if ang.deg < 0 {
		sign = -1
	}
	s := math.Abs(ang.ArcSeconds())
	sec := math.Mod(s, 60.)
	m := (s - sec) / 60.
	min := math.Mod(m, 60.)
	deg := (m - min) / 60.
	return sign * deg, min, sec
}

func (ang *Angle) HMS() (float64, float64, float64) {
	ang.deg /= 15.
	h, m, s := ang.DMS()
	ang.deg *= 15.
	return h, m, s
}

func (ang *Angle) Hour() float64 {
	return ang.deg / 15.
}

func (ang *Angle) Minutes() float64 {
	return ang.deg / 15. * 60.
}

func (ang *Angle) Seconds() float64 {
	return ang.deg / 15. * 60. * 60.
}

func (ang *Angle) Degree() float64 {
	return ang.deg
}

func (ang *Angle) ArcMinutes() float64 {
	return ang.deg * 60.
}

func (ang *Angle) ArcSeconds() float64 {
	return ang.deg * 60. * 60.
}

func (ang *Angle) Radian() float64 {
	return ang.deg * math.Pi / 180.
}

func (ang *Angle) String(format string) string {
	switch format {
	case `deg`, `degree`:
		return fmt.Sprintf("%.8fd", ang.Degree())
	case `arcmin`:
		return fmt.Sprintf("%.8f'", ang.ArcMinutes())
	case `arcsec`:
		return fmt.Sprintf("%.8f\"", ang.ArcSeconds())
	case `hour`:
		return fmt.Sprintf("%.8fh", ang.Hour())
	case `min`:
		return fmt.Sprintf("%.8fm", ang.Minutes())
	case `sec`:
		return fmt.Sprintf("%.8fs", ang.Seconds())
	case `rad`, `radian`:
		return fmt.Sprintf("%.8frad", ang.Radian())
	case `dms`:
		d, m, s := ang.DMS()
		return fmt.Sprintf("%02.0fd%02.0f'%.8f\"", d, m, s)
	case `hms`:
		h, m, s := ang.HMS()
		return fmt.Sprintf("%02.0fh%02.0fm%.8fs", h, m, s)
	default:
		log.Printf("Unknown format %s. Return degree.")
		return fmt.Sprintf("%.8fd", ang.Degree())
	}
}

/* Coordinates */
type Spherical struct {
	X *Angle
	Y *Angle
}

func (s *Spherical) ToEq() *Spherical {
	x := s.X.Radian()
	y := s.Y.Radian()

	ra := math.Mod(x, 2.*math.Pi)
	if ra < 0 {
		ra += 2. * math.Pi
	}

	dec := math.Mod(y, 2.*math.Pi)
	if math.Abs(dec) > math.Pi {
		if y > 0 {
			dec -= math.Pi
		}
		if y < 0 {
			dec += math.Pi
		}
	}
	return &Spherical{X: NewAngle(RadToDeg(ra)), Y: NewAngle(RadToDeg(dec))}
}

func (s *Spherical) ToGal() *Spherical {
	x := s.X.Radian()
	y := s.Y.Radian()

	l := math.Mod(x, 2.*math.Pi)
	if l < 0 {
		l += 2. * math.Pi
	}
	b := math.Mod(y, 2.*math.Pi)
	if math.Abs(b) > math.Pi {
		if y > 0 {
			b -= math.Pi
		}
		if y < 0 {
			b += math.Pi
		}
	}
	return &Spherical{X: NewAngle(RadToDeg(l)), Y: NewAngle(RadToDeg(b))}
}

type Cartesian struct {
	X float64
	Y float64
	Z float64
}

func (s *Spherical) ToCartesian() *Cartesian {
	x := s.X.Radian()
	y := s.Y.Radian()
	return &Cartesian{
		X: math.Cos(x) * math.Cos(y),
		Y: math.Sin(x) * math.Cos(y),
		Z: math.Sin(y),
	}
}

func (c *Cartesian) ToSpherical() *Spherical {
	x := c.X
	y := c.Y
	z := c.Z
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
	return &Spherical{
		X: NewAngle(RadToDeg(h)),
		Y: NewAngle(RadToDeg(v)),
	}
}

type Coordinate interface {
	ConvertTo(string) Coordinate
	ToCartesian() *Cartesian
	GetX() *Angle
	GetY() *Angle
}

type coordinate struct {
	*Spherical
	system string
}

func (c coordinate) GetX() *Angle {
	return c.X
}

func (c coordinate) GetY() *Angle {
	return c.Y
}

type B1950 struct {
	*coordinate
}

type J2000 struct {
	*coordinate
}

type Gal struct {
	*coordinate
}

func (c coordinate) ConvertTo(newsystem string) Coordinate {
	switch c.system {
	case `J2000`:
		switch newsystem {
		case `J2000`:
			return c
		case `B1950`:
			return J2000ToB1950(&c)
		case `Gal`:
			return J2000ToGal(&c)
		default:
			log.Fatal(fmt.Errorf("Unknown newsystem %s", newsystem))
		}
	case `B1950`:
		switch newsystem {
		case `J2000`:
			return B1950ToJ2000(&c)
		case `B1950`:
			return c
		case `Gal`:
			return B1950ToGal(&c)
		default:
			log.Fatal(fmt.Errorf("Unknown newsystem %s", newsystem))
		}
	case `Gal`:
		switch newsystem {
		case `J2000`:
			return GalToJ2000(&c)
		case `B1950`:
			return GalToB1950(&c)
		case `Gal`:
			return c
		default:
			log.Fatal(fmt.Errorf("Unknown newsystem %s", newsystem))
		}
	default:
		log.Fatal(fmt.Errorf("Unknown input system %s", c.system))
	}
	return nil
}

func NewCoordinate(system string, x, y float64) Coordinate {
	s := &Spherical{X: NewAngle(x), Y: NewAngle(y)}
	return NewCoordinateFromSphere(system, s)
}

func NewCoordinateFromSphere(system string, s *Spherical) Coordinate {
	c := coordinate{Spherical: s, system: system}
	switch system {
	case `J2000`:
		return &J2000{&c}
	case `B1950`:
		return &B1950{&c}
	case `Gal`:
		return &Gal{&c}
	default:
		log.Fatal(fmt.Errorf("Unknwon system %s", system))
	}
	return nil
}

func (coord *B1950) String() string {
	return fmt.Sprintf(`B1950, RA: %s, DEC: %s`, coord.X.String(`hms`), coord.Y.String(`dms`))
}

func (coord *B1950) Ra() *Angle {
	return coord.X
}

func (coord *B1950) Dec() *Angle {
	return coord.Y
}

func (coord *J2000) String() string {
	return fmt.Sprintf(`J2000, RA: %s, DEC: %s`, coord.X.String(`hms`), coord.Y.String(`dms`))
}

func (coord *J2000) Ra() *Angle {
	return coord.X
}

func (coord *J2000) Dec() *Angle {
	return coord.Y
}

func (coord *Gal) String() string {
	return fmt.Sprintf(`Galac, Lat: %s, Lon: %s`, coord.X.String(`deg`), coord.Y.String(`deg`))
}

func (coord *Gal) Longitude() *Angle {
	return coord.X
}

func (coord *Gal) Latitude() *Angle {
	return coord.Y
}

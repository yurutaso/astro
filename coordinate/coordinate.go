package coordinate

import (
	"fmt"
	"log"
	"math"
)

type Degree struct {
	deg float64
}

func NewDegree(deg float64) *Degree {
	return &Degree{deg: deg}
}

func (deg *Degree) DMS() (float64, float64, float64) {
	sign := 1.
	if deg.deg < 0 {
		sign = -1
	}
	d := math.Abs(deg.deg)
	hour := math.Floor(d)
	d = (d - hour) * 60.
	min := math.Floor(d)
	sec := (d - min) * 60.
	return sign * hour, min, sec
}

func (deg *Degree) HMS() (float64, float64, float64) {
	deg.deg /= 15.
	h, m, s := deg.DMS()
	deg.deg *= 15.
	return h, m, s
}

func (deg *Degree) Hour() float64 {
	return deg.deg / 15.
}

func (deg *Degree) Minutes() float64 {
	return deg.deg / 15. * 60.
}

func (deg *Degree) Seconds() float64 {
	return deg.deg / 15. * 60. * 60.
}

func (deg *Degree) Degree() float64 {
	return deg.deg
}

func (deg *Degree) ArcMinutes() float64 {
	return deg.deg * 60.
}

func (deg *Degree) ArcSeconds() float64 {
	return deg.deg * 60. * 60.
}

func (deg *Degree) Radian() float64 {
	return deg.deg * math.Pi / 180.
}

func (deg *Degree) String(format string) string {
	switch format {
	case `deg`, `degree`:
		return fmt.Sprintf("%.8fd", deg.Degree())
	case `arcmin`:
		return fmt.Sprintf("%.8f'", deg.ArcMinutes())
	case `arcsec`:
		return fmt.Sprintf("%.8f\"", deg.ArcSeconds())
	case `hour`:
		return fmt.Sprintf("%.8fh", deg.Hour())
	case `min`:
		return fmt.Sprintf("%.8fm", deg.Minutes())
	case `sec`:
		return fmt.Sprintf("%.8fs", deg.Seconds())
	case `rad`, `radian`:
		return fmt.Sprintf("%.8frad", deg.Radian())
	case `dms`:
		d, m, s := deg.DMS()
		return fmt.Sprintf("%02.0fd%02.0f'%.8f\"", d, m, s)
	case `hms`:
		h, m, s := deg.HMS()
		return fmt.Sprintf("%02.0fh%02.0fm%.8fs", h, m, s)
	default:
		log.Printf("Unknown format %s. Return degree.")
		return fmt.Sprintf("%.8fd", deg.Degree())
	}
}

func DMSToDeg(deg, min, sec float64) *Degree {
	if deg < 0 {
		return &Degree{deg: -(-deg + min/60. + sec/3600.)}
	}
	return &Degree{deg: deg + min/60. + sec/3600.}
}

func HMSToDeg(hour, min, sec float64) *Degree {
	deg := DMSToDeg(hour, min, sec)
	deg.deg *= 15.
	return deg
}

type Coordinate interface {
	String() string
	ConvertTo(string) Coordinate
}

type B1950 struct {
	Ra  *Degree
	Dec *Degree
}

type J2000 struct {
	Ra  *Degree
	Dec *Degree
}

type Gal struct {
	Lon *Degree
	Lat *Degree
}

func (coord B1950) String() string {
	return fmt.Sprintf(`RA: %s, DEC: %s`, coord.Ra.String(`hms`), coord.Dec.String(`dms`))
}

func (coord B1950) ConvertTo(newcoord string) Coordinate {
	switch newcoord {
	case `J2000`:
		return coord
	case `B1950`:
		return coord
	case `Gal`:
		return coord
	default:
		return coord
	}
}

func (coord J2000) String() string {
	return fmt.Sprintf(`RA: %s, DEC: %s`, coord.Ra.String(`hms`), coord.Dec.String(`dms`))
}

func (coord J2000) ConvertTo(newcoord string) Coordinate {
	switch newcoord {
	case `J2000`:
		return coord
	case `B1950`:
		return coord
	case `Gal`:
		return coord
	default:
		return coord
	}
}

func (coord Gal) String() string {
	return fmt.Sprintf(`l: %s, b: %s`, coord.Lon.String(`deg`), coord.Lat.String(`deg`))
}

func (coord Gal) ConvertTo(newcoord string) Coordinate {
	switch newcoord {
	case `J2000`:
		return coord
	case `B1950`:
		return coord
	case `Gal`:
		return coord
	default:
		return coord
	}
}

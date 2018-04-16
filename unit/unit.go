package unit

import (
	"fmt"
	"math"
)

/* Interface Unit is defined as BaseUnit with dimension */
/* In the usual case, use "Units" rather than "Unit" */
type Unit interface {
	BaseUnit
	String() string
	Copy() Unit
	GetBaseUnit() BaseUnit
	Dimension() float64
	SetDimension(float64)
	AddDimension(float64)
	Inverse() Unit
	Equal(Unit) bool
	RatioTo(Unit) (float64, error)
}

type _unit struct {
	BaseUnit
	dim float64
}

func (u *_unit) String() string {
	switch u.dim {
	case 0:
		return ""
	case 1:
		return u.Name()
	default:
		return fmt.Sprintf("%s^%.0f", u.Name(), u.dim)
	}
}
func (u *_unit) Dimension() float64 {
	return u.dim
}
func (u *_unit) SetDimension(dim float64) {
	u.dim = dim
}
func (u *_unit) AddDimension(dim float64) {
	u.dim += dim
}
func (u *_unit) Inverse() Unit {
	return &_unit{BaseUnit: u.BaseUnit, dim: -u.dim}
}
func (u *_unit) Equal(that Unit) bool {
	return u.Type() == that.Type() && u.dim == that.Dimension()
}
func (u *_unit) RatioTo(that Unit) (float64, error) {
	if that.Equal(u) {
		return math.Pow(u.Prefix()/that.Prefix(), u.dim), nil
	}
	return 1., fmt.Errorf(`Cannot compute ratio of units with different dimensions`)
}

func (u *_unit) GetBaseUnit() BaseUnit {
	return u.BaseUnit
}

func (u *_unit) Copy() Unit {
	return &_unit{BaseUnit: u.BaseUnit, dim: u.dim}
}

/* Composit Unit */
type Units interface {
	Inverse() Units
	Get(string) Unit
	GetAll() map[string]Unit
	Set(Unit)
	Copy() Units
	Equal(Units) bool
	Has(string) bool
	Multiply(...Units) (Units, float64)
	Times(float64) Units
}

type _units struct {
	units map[string]Unit
}

func (units *_units) Times(val float64) Units {
	timed := units.Copy()
	for _, unit := range timed.GetAll() {
		unit.SetDimension(val * unit.Dimension())
	}
	return timed
}

func (units *_units) Inverse() Units {
	inversed := units.Copy()
	for _, unit := range inversed.GetAll() {
		unit.SetDimension(-unit.Dimension())
	}
	return inversed
}

func (units *_units) Set(unit Unit) {
	units.units[unit.Type()] = unit.Copy()
}

func (units *_units) GetAll() map[string]Unit {
	return units.units
}

func (units *_units) Get(utype string) Unit {
	for t, unit := range units.units {
		if t == utype {
			return unit
		}
	}
	return nil
}

func (units *_units) Has(utype string) bool {
	for t, _ := range units.units {
		if t == utype {
			return true
		}
	}
	return false
}

func (units *_units) Copy() Units {
	copied := map[string]Unit{}
	for utype, unit := range units.units {
		copied[utype] = unit.Copy()
	}
	return &_units{units: copied}
}

func (units *_units) Equal(that Units) bool {
	// True if both units have same DIMENSION (prefix is NOT considered)
	if len(units.units) != len(that.GetAll()) {
		return false
	}
	for utype, unit1 := range units.units {
		if that.Has(utype) && unit1.Equal(that.Get(utype)) {
			continue
		}
		return false
	}
	return true
}

func (units *_units) Multiply(those ...Units) (Units, float64) {
	// Returns multiplied units, and a factor which should be multiplied to the value
	// when using this units.
	// If receiver and input units contain the same unit type, use that of the receiver.
	// e.g. m/s * (kg*km) returns (m^2*kg/s, 1000.)
	newunits := units.Copy()
	f := 1.
	for _, that := range those {
		for utype, unit := range that.GetAll() {
			if newunits.Has(utype) {
				_f, _ := unit.RatioTo(newunits.Get(utype)) // Error never happens if utype is the same.
				f *= _f
				newunits.Get(utype).AddDimension(unit.Dimension())
			} else {
				newunits.Set(unit)
			}
		}
	}
	return newunits, f
}

/* IO */
func Empty() Units {
	return &_units{units: map[string]Unit{}}
}

func Multiply(units ...Units) (Units, float64) {
	switch len(units) {
	case 0:
		return nil, 0.
	case 1:
		return units[0], 1.
	default:
		return units[0].Multiply(units[1:]...)
	}
}

func NewSingleUnit(unit BaseUnit, dim float64) Units {
	u := &_unit{BaseUnit: unit, dim: dim}
	units := Empty()
	units.Set(u)
	return units
}

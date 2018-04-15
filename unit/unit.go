package unit

import (
	"fmt"
	"math"
)

/* unit with dimension */
type Unit interface {
	BaseUnit
	Dimension() float64
	setDimension(float64)
	addDimension(float64)
}

type unit struct {
	BaseUnit
	dim float64
}

func NewUnit(u BaseUnit, dim float64) Unit {
	return &unit{BaseUnit: u, dim: dim}
}

func (u *unit) Dimension() float64 {
	return u.dim
}
func (u *unit) setDimension(dim float64) {
	u.dim = dim
}
func (u *unit) addDimension(dim float64) {
	u.dim += dim
}

/* Composit unit */
type Units interface {
	Get(string) Unit
	Set(Unit)
	GetAll() map[string]Unit
	Copy() Units
	Equal(Units) bool
	Has(string) bool
	Multiply(Units) (Units, float64)
}

type units struct {
	units map[string]Unit
}

func UnitsFromSlice(u ...Unit) Units {
	newunits := map[string]Unit{}
	for _, unit := range u {
		newunits[unit.Type()] = unit
	}
	return &units{units: newunits}
}

func (u *units) Set(unit Unit) {
	u.units[unit.Type()] = unit
}
func (u *units) GetAll() map[string]Unit {
	return u.units
}
func (u *units) Get(utype string) Unit {
	if u.Has(utype) {
		return u.units[utype]
	}
	return nil
}

func (u *units) Copy() Units {
	copied := map[string]Unit{}
	for utype, unit := range u.units {
		b := BaseUnitOf(utype, unit.Name(), unit.Prefix())
		copied[utype] = NewUnit(b, unit.Dimension())
	}
	return &units{units: copied}
}

func (u *units) Equal(that Units) bool {
	if len(u.units) != len(that.GetAll()) {
		return false
	}
	for utype, unit1 := range u.units {
		if that.Has(utype) {
			if unit1.Dimension() != that.Get(utype).Dimension() {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func (u *units) Has(utype string) bool {
	// True if Units has type utype
	for t, _ := range u.units {
		if t == utype {
			return true
		}
	}
	return false
}

func (u *units) Multiply(that Units) (Units, float64) {
	// Multiply Units
	// Output Unit is based on that of the receiver.
	// Use the Unit of the input, if the receiver does not have that utype.
	// e.g. m/s * (kg*km) returns (m^2*kg/s, 1000.)
	newunits := u.Copy()
	f := 1.
	for utype, unit := range that.GetAll() {
		if u.Has(utype) {
			f *= math.Pow(unit.Prefix()/u.Get(utype).Prefix(), unit.Dimension())
			newunits.Get(utype).addDimension(unit.Dimension())
		} else {
			newunits.Set(unit)
		}
	}
	return newunits, f
}

/* value with unit */
type UnitValue interface {
	Units() Units
	Value() float64
	Equal(UnitValue) bool
	As(Units) (UnitValue, error)
	Multiply(UnitValue) UnitValue
}

type unitValue struct {
	units Units
	value float64
}

func NewUnitValue(value float64, u Units) UnitValue {
	return &unitValue{units: u, value: value}
}
func NewUnitValueFromSlice(value float64, u ...Unit) UnitValue {
	return &unitValue{units: UnitsFromSlice(u...), value: value}
}

func (uv *unitValue) Value() float64 {
	return uv.value
}

func (uv *unitValue) Units() Units {
	return uv.units
}

func (uv *unitValue) Equal(that UnitValue) bool {
	if converted, err := that.As(uv.units); err != nil {
		return false
	} else {
		return uv.value == converted.Value()
	}
}

func (uv *unitValue) As(units Units) (UnitValue, error) {
	if !uv.units.Equal(units) {
		return nil, fmt.Errorf(`Cannot convert to unit with different dimensions.`)
	}
	f := 1.
	for utype, before := range uv.units.GetAll() {
		after := units.Get(utype)
		f *= math.Pow(before.Prefix()/after.Prefix(), before.Dimension())
	}
	return NewUnitValue(f*uv.value, units.Copy()), nil
}

func (uv *unitValue) Multiply(that UnitValue) UnitValue {
	newunit, f := uv.units.Multiply(that.Units())
	return NewUnitValue(f*uv.value*that.Value(), newunit)
}

/* physical constants */
func C() UnitValue {
	value := 299792458.
	return NewUnitValueFromSlice(value, NewUnit(Meter(), 1.), NewUnit(Second(), -1.))
}
func Kb() UnitValue {
	value := 1.38064852e-23
	return NewUnitValueFromSlice(value,
		NewUnit(Meter(), 2.),
		NewUnit(KiloGram(), 1.),
		NewUnit(Second(), -2.),
		NewUnit(Kelvin(), -1.),
	)
}

package unit

import (
	"fmt"
	"math"
)

/* Value with unit */
type UnitValue interface {
	String() string
	Units() Units
	Value() float64
	SetValue(float64)
	Equal(UnitValue) bool
	As(Units) (UnitValue, error)
	MultiplyValue(float64) UnitValue
	Multiply(...UnitValue) UnitValue
	Divide(...UnitValue) UnitValue
	Add(...UnitValue) (UnitValue, error)
	Inverse() UnitValue
}

type unitValue struct {
	units Units
	value float64
}

func (uv *unitValue) String() string {
	s := ""
	for _, unit := range uv.units.GetAll() {
		if _s := unit.String(); len(_s) != 0 {
			s += _s + " "
		}
	}
	if len(s) != 0 {
		return fmt.Sprintf("%e (%s)", uv.value, s[:len(s)-1])
	}
	return fmt.Sprintf("%e", uv.value)
}

func (uv *unitValue) Inverse() UnitValue {
	return NewUnitValue(1./uv.value, uv.Units().Inverse())
}
func (uv *unitValue) Value() float64 {
	return uv.value
}
func (uv *unitValue) SetValue(value float64) {
	uv.value = value
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

func (uv *unitValue) MultiplyValue(value float64) UnitValue {
	return NewUnitValue(value*uv.value, uv.units.Copy())
}

func (uv *unitValue) Multiply(those ...UnitValue) UnitValue {
	newunit := uv.units.Copy()
	f := 1.
	_f := 1.
	v := uv.value
	for _, that := range those {
		newunit, _f = that.Units().Multiply(newunit)
		f *= _f
		v *= that.Value()
	}
	return NewUnitValue(f*v, newunit)
}

func (uv *unitValue) Divide(those ...UnitValue) UnitValue {
	newunit := uv.units.Copy()
	f := 1.
	_f := 1.
	v := uv.value
	for _, that := range those {
		newunit, _f = newunit.Multiply(that.Units().Inverse())
		f *= _f
		v /= that.Value()
	}
	return NewUnitValue(f*v, newunit)
}

func (uv *unitValue) Add(those ...UnitValue) (UnitValue, error) {
	var converted UnitValue
	var err error
	v := uv.value
	for _, that := range those {
		if converted, err = that.As(uv.units); err != nil {
			return nil, fmt.Errorf("Cannot add UnitValue with different units.")
		}
		v += converted.Value()
	}
	return NewUnitValue(v, uv.units.Copy()), nil
}

/* IO */
func NewUnitValue(value float64, units ...Units) UnitValue {
	u, f := Multiply(units...)
	return &unitValue{units: u, value: value * f}
}

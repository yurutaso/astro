package constants

import (
	"github.com/yurutaso/astro/unit"
)

var (
	m  unit.Units = unit.Meter().AsUnits(1.)    // meter
	s  unit.Units = unit.Second().AsUnits(1.)   // second
	k  unit.Units = unit.Kelvin().AsUnits(1.)   // kelvin
	kg unit.Units = unit.KiloGram().AsUnits(1.) // kilogram
)

/* physical constants */
func C() unit.UnitValue {
	value := 299792458.
	return unit.NewUnitValue(value, m, s.Inverse())
}

func Kb() unit.UnitValue {
	value := 1.38064852e-23
	return unit.NewUnitValue(value, m.Times(2), kg, s.Times(-2), k.Inverse())
}

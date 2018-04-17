package unit

const (
	UNITTYPE_LENGTH      string = `length`
	UNITTYPE_MASS        string = `mass`
	UNITTYPE_TIME        string = `time`
	UNITTYPE_TEMPERATURE string = `temperature`

	PREFIX_FEMT  float64 = 1.e-15
	PREFIX_PICO  float64 = 1.e-12
	PREFIX_NANO  float64 = 1.e-9
	PREFIX_MICRO float64 = 1.e-6
	PREFIX_MILI  float64 = 1.e-3
	PREFIX_CENTI float64 = 1.e-2
	PREFIX_KILO  float64 = 1.e3
	PREFIX_MEGA  float64 = 1.e6
	PREFIX_GIGA  float64 = 1.e9
	PREFIX_TERA  float64 = 1.e12

	PREFIX_PARSEC float64 = 3.085677581e16
	PREFIX_AU     float64 = 149597870700.
	PREFIX_LY     float64 = 9460730472580800.
)

/* Base unit */
type BaseUnit interface {
	Type() string // e.g. length, time, mass etc.
	Name() string // e.g. meter, second, gram
	SetName(string)
	Prefix() float64 // e.g Kilo, Giga
	SetPrefix(float64)
	AsUnits(float64) Units
}

/* Base unit of a specified type */
type baseUnit struct {
	utype  string
	name   string
	prefix float64
}

func (u *baseUnit) AsUnits(dim float64) Units {
	return newSingleUnit(u, dim)
}
func (u *baseUnit) Type() string {
	return u.utype
}
func (u *baseUnit) Name() string {
	return u.name
}
func (u *baseUnit) SetName(name string) {
	u.name = name
}
func (u *baseUnit) Prefix() float64 {
	return u.prefix
}
func (u *baseUnit) SetPrefix(prefix float64) {
	u.prefix = prefix
}

/* Generator */
func BaseUnitOf(utype string, name string, prefix float64) BaseUnit {
	return &baseUnit{utype: utype, name: name, prefix: prefix}
}
func BaseUnitOfLength(name string, prefix float64) BaseUnit {
	return &baseUnit{utype: UNITTYPE_LENGTH, name: name, prefix: prefix}
}
func BaseUnitOfMass(name string, prefix float64) BaseUnit {
	return &baseUnit{utype: UNITTYPE_MASS, name: name, prefix: prefix}
}
func BaseUnitOfTime(name string, prefix float64) BaseUnit {
	return &baseUnit{utype: UNITTYPE_TIME, name: name, prefix: prefix}
}
func BaseUnitOfTemperature(name string, prefix float64) BaseUnit {
	return &baseUnit{utype: UNITTYPE_TEMPERATURE, name: name, prefix: prefix}
}

/* Base Units */
func meter() BaseUnit {
	return BaseUnitOfLength(`m`, 1.)
}

func kiloGram() BaseUnit {
	return BaseUnitOfMass(`kg`, 1.e3)
}

func second() BaseUnit {
	return BaseUnitOfTime(`s`, 1.)
}

func kelvin() BaseUnit {
	return BaseUnitOfTemperature(`K`, 1.)
}

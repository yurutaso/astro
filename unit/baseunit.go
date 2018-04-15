package unit

const (
	UNITTYPE_LENGTH      string = `length`
	UNITTYPE_MASS        string = `mass`
	UNITTYPE_TIME        string = `time`
	UNITTYPE_TEMPERATURE string = `temperature`
)

/* Base unit */
type BaseUnit interface {
	Type() string    // e.g. length, time, mass etc.
	Name() string    // e.g. meter, second, gram
	Prefix() float64 // e.g Kilo, Giga
}

/* Base unit of a specified type */
type baseUnit struct {
	utype  string
	name   string
	prefix float64
}

func (u *baseUnit) Type() string {
	return u.utype
}
func (u *baseUnit) Name() string {
	return u.name
}
func (u *baseUnit) Prefix() float64 {
	return u.prefix
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
func Meter() BaseUnit {
	return BaseUnitOfLength(`m`, 1.)
}
func CentiMeter() BaseUnit {
	return BaseUnitOfLength(`cm`, 1.e-2)
}

func Gram() BaseUnit {
	return BaseUnitOfMass(`g`, 1.)
}
func KiloGram() BaseUnit {
	return BaseUnitOfMass(`kg`, 1.e3)
}

func Second() BaseUnit {
	return BaseUnitOfTime(`s`, 1.)
}

func Kelvin() BaseUnit {
	return BaseUnitOfTemperature(`K`, 1.)
}

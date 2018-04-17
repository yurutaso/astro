package unit

var (
	currentSystem UnitSystem = SIUnit()
)

/* IO */
func GetCurrentSystem() UnitSystem {
	return currentSystem
}

func SetSystem(u UnitSystem) {
	currentSystem = u
}

func SetSystemOf(utype string, u BaseUnit) {
	currentSystem.BaseUnits()[utype] = u
}

/* Unit system */
type UnitSystem interface {
	BaseUnits() map[string]BaseUnit
	Copy() UnitSystem
}

type unitSystem struct {
	units map[string]BaseUnit
}

func SIUnit() UnitSystem {
	units := map[string]BaseUnit{
		UNITTYPE_LENGTH:      meter(),
		UNITTYPE_MASS:        kiloGram(),
		UNITTYPE_TIME:        second(),
		UNITTYPE_TEMPERATURE: kelvin(),
	}
	return &unitSystem{units: units}
}

func (s *unitSystem) BaseUnits() map[string]BaseUnit {
	return s.units
}

func (s *unitSystem) Copy() UnitSystem {
	copied := map[string]BaseUnit{}
	for utype, u := range s.units {
		copied[utype] = u
	}
	return &unitSystem{units: copied}
}

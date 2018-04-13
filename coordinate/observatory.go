package coordinate

const (
	/* ToDo: set collect values for following observatory*/
	ASTE_LATITUDE        float64 = 0.1
	ASTE_LONGITUDE       float64 = 0.1
	EFFELSBERG_LATITUDE  float64 = 0.1
	EFFELSBERG_LONGITUDE float64 = 0.1
	GBT_LATITUDE         float64 = 0.1
	GBT_LONGITUDE        float64 = 0.1
	NRO_LATITUDE         float64 = 0.1
	NRO_LONGITUDE        float64 = 0.1
	VLA_LATITUDE         float64 = 0.1
	VLA_LONGITUDE        float64 = 0.1
)

type Observatory interface {
	Latitude() *Angle
	Longitude() *Angle
	Name() string
}

type observatory struct {
	latitude  *Angle
	longitude *Angle
	name      string
}

/* Implement interface Observatory for observatory */
func (o *observatory) Latitude() *Angle {
	return o.latitude
}

func (o *observatory) Longitude() *Angle {
	return o.longitude
}

func (o *observatory) Name() string {
	return o.name
}

/* IO */
func NewObservatoryFromAngles(lat, lon *Angle, name string) Observatory {
	return &observatory{latitude: lat, longitude: lon, name: name}
}

func NewObservatory(lat, lon float64, name string) Observatory {
	return &observatory{latitude: NewAngle(lat), longitude: NewAngle(lon), name: name}
}

/* Actual observatories */
func ASTE() Observatory {
	return NewObservatory(ASTE_LATITUDE, ASTE_LONGITUDE, `ASTE`)
}

func Effelsberg() Observatory {
	return NewObservatory(EFFELSBERG_LATITUDE, EFFELSBERG_LONGITUDE, `EFFELSBERG`)
}

func GBT() Observatory {
	return NewObservatory(GBT_LATITUDE, GBT_LONGITUDE, `GBT`)
}

func NRO() Observatory {
	return NewObservatory(NRO_LATITUDE, NRO_LONGITUDE, `NRO`)
}

func VLA() Observatory {
	return NewObservatory(VLA_LATITUDE, VLA_LONGITUDE, `VLA`)
}

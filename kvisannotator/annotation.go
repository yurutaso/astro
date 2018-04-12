package kvisannotator

import (
	"fmt"
	"github.com/yurutaso/astro/coordinate"
	"os"
)

// Annotation file consists of header string and Annotations
type AnnotationFile struct {
	Annotations []Annotation
}

func NewAnnotationFile() *AnnotationFile {
	return &AnnotationFile{
		Annotations: make([]Annotation, 0, 0),
	}
}

func (af *AnnotationFile) Add(ann Annotation) {
	af.Annotations = append(af.Annotations, ann)
}

func (af *AnnotationFile) WriteTo(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, ann := range af.Annotations {
		file.Write([]byte(ann.String()))
	}
	return nil
}

type AnnotationOpt struct {
	Color         string
	PositionAngle string
	Coord         string
	Font          string
}

func NewAnnotationOpt() *AnnotationOpt {
	return &AnnotationOpt{Color: `WHITE`, PositionAngle: `STANDARD`, Coord: `W`}
}

func (opt *AnnotationOpt) String() string {
	s := ""
	if len(opt.Color) != 0 {
		s += fmt.Sprintf("COLOR %s\n", opt.Color)
	}
	if len(opt.PositionAngle) != 0 {
		s += fmt.Sprintf("PA %s\n", opt.PositionAngle)
	}
	if len(opt.Coord) != 0 {
		s += fmt.Sprintf("COORD %s\n", opt.Coord)
	}
	if len(opt.Font) != 0 {
		s += fmt.Sprintf("Font %s\n", opt.Font)
	}
	return s
}

// Annotation
type Annotation interface {
	String() string
}

/* Annotation Line */
type Line struct {
	From   coordinate.Coordinate
	To     coordinate.Coordinate
	Option *AnnotationOpt
}

func NewLine(from, to coordinate.Coordinate) *Line {
	return &Line{From: from, To: to, Option: NewAnnotationOpt()}
}

func (ann Line) String() string {
	from := ann.From.ConvertTo(`J2000`)
	to := ann.To.ConvertTo(`J2000`)
	return fmt.Sprintf("%sLINE W %f %f W %f %f\n", ann.Option, from.GetX().Degree(), from.GetY().Degree(), to.GetX().Degree(), to.GetY().Degree())
}

/* Annotation Circle */
type Circle struct {
	Center coordinate.Coordinate
	Width  float64
	Option *AnnotationOpt
}

func NewCircle(center coordinate.Coordinate, width float64) *Circle {
	return &Circle{Center: center, Width: width, Option: NewAnnotationOpt()}
}

func (ann Circle) String() string {
	center := ann.Center.ConvertTo(`J2000`)
	return fmt.Sprintf("%sCIRCLE W %f %f %f\n", ann.Option, center.GetX().Degree(), center.GetY().Degree(), ann.Width)
}

/* Annotation Point */
type Dot struct {
	Center coordinate.Coordinate
	Option *AnnotationOpt
}

func NewDot(center coordinate.Coordinate) *Dot {
	return &Dot{Center: center, Option: NewAnnotationOpt()}
}

func (ann Dot) String() string {
	center := ann.Center.ConvertTo(`J2000`).(*coordinate.J2000)
	return fmt.Sprintf("%sDot W %f %f\n", ann.Option, center.GetX().Degree(), center.GetY().Degree())
}

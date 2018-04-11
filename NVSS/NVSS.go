package NVSS

import (
	"bufio"
	"fmt"
	"github.com/yurutaso/astro/coordinate"
	"math"
	"os"
	"strconv"
	"strings"
)

type Source struct {
	Coord coordinate.J2000
	Flux  float64
}

func (source *Source) String() string {
	return fmt.Sprintf("coord: %s, flux: %f", source.Coord, source.Flux)
}

type Catalog struct {
	Sources []Source
}

type Filter struct {
	RaMin   float64
	RaMax   float64
	DecMin  float64
	DecMax  float64
	FluxMin float64
	FluxMax float64
}

func NewFilter() *Filter {
	return &Filter{
		RaMin:   math.Inf(-1),
		RaMax:   math.Inf(1),
		DecMin:  math.Inf(-1),
		DecMax:  math.Inf(1),
		FluxMin: math.Inf(-1),
		FluxMax: math.Inf(1),
	}
}

func NewCatalogFromText(filename string) (*Catalog, error) {
	fp, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	sources := make([]Source, 0, 0)
	scanner := bufio.NewScanner(fp)

	var ra_h, ra_m, ra_s, dec_d, dec_m, dec_s, flux float64

	fmt.Printf(fmt.Sprintf("Scanning %s\n", filename))
	for scanner.Scan() {
		var line string = scanner.Text()
		if len(line) == 0 {
			continue
		}
		switch line[0] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			data := strings.Fields(line)
			if ra_h, err = strconv.ParseFloat(data[0], 64); err != nil {
				return nil, err
			}
			if ra_m, err = strconv.ParseFloat(data[1], 64); err != nil {
				return nil, err
			}
			if ra_s, err = strconv.ParseFloat(data[2], 64); err != nil {
				return nil, err
			}
			if dec_d, err = strconv.ParseFloat(data[3], 64); err != nil {
				return nil, err
			}
			if dec_m, err = strconv.ParseFloat(data[4], 64); err != nil {
				return nil, err
			}
			if dec_s, err = strconv.ParseFloat(data[5], 64); err != nil {
				return nil, err
			}
			if flux, err = strconv.ParseFloat(data[7], 64); err != nil {
				return nil, err
			}
			ra := &coordinate.HMSToDeg(ra_h, ra_m, ra_s)
			dec := &coordinate.DMSToDeg(dec_d, dec_m, dec_s)
			source := Source{
				Coord: &coordinate.J2000{Ra: ra, Dec: dec},
				Flux:  flux,
			}
			sources = append(sources, source)
		}
	}
	fmt.Printf("Done. %d sources found.\n", len(sources))
	return &Catalog{Sources: sources}, nil
}

func (cat *Catalog) Filter(filter *Filter) *Catalog {
	sources := make([]Source, 0, 0)
	for _, source := range cat.Sources {
		ra := source.Coord.Ra
		dec := source.Coord.Dec
		flux := source.Flux
		if filter.RaMin <= ra && ra <= filter.RaMax {
			if filter.DecMin <= dec && dec <= filter.DecMax {
				if filter.FluxMin <= flux && flux <= filter.FluxMax {
					sources = append(sources, source)
				}
			}
		}
	}
	return &Catalog{Sources: sources}
}

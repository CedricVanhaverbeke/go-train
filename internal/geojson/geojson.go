package geojson

import "math"

type GeoJson struct {
	Type     string    `json:"type"`
	Features []feature `json:"features"`
}

func newGeometry() geometry {
	return geometry{
		Type:        "LineString",
		Coordinates: [][]float64{},
	}
}

type geometry struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

func newFeature() feature {
	return feature{
		Type:     "Feature",
		Geometry: newGeometry(),
	}
}

type feature struct {
	Type     string `json:"type"`
	Geometry geometry
}

func New() GeoJson {
	return GeoJson{
		Type:     "FeatureCollection",
		Features: []feature{newFeature()},
	}
}

func (g *GeoJson) Add(lon, lat, elevation float64) {
	g.Features[0].Geometry.Coordinates = append(
		g.Features[0].Geometry.Coordinates,
		[]float64{lon, lat, elevation})
}

// distance returns the distance of a segment
// with a given start and end index
//
// it uses the Haversine formula
func (g *GeoJson) distance(x int, y int) float64 {
	if g.Features == nil {
		return 0.0
	}

	if g.Features[0].Geometry.Coordinates == nil {
		return 0.0
	}

	if len(g.Features[0].Geometry.Coordinates) < x ||
		len(g.Features[0].Geometry.Coordinates) < y {
		return 0.0
	}

	var d float64
	for i := x; i < y-1; i++ {
		c1 := g.Features[0].Geometry.Coordinates[i]
		c2 := g.Features[0].Geometry.Coordinates[i+1]
		d += haversine(c1[0], c1[1], c2[0], c2[1])
	}
	return d
}

var EARTH_RADIUS = 6371e3

// stolen frim here https://www.movable-type.co.uk/scripts/latlong.html
func haversine(lng1 float64, lat1 float64, lng2 float64, lat2 float64) float64 {
	r1 := lat1 * (math.Pi / 180)
	r2 := lat2 * (math.Pi / 180)

	d1 := (lat2 - lat1) * (math.Pi / 180)
	d2 := (lng2 - lng1) * (math.Pi / 180)

	a := math.Pow(math.Sin(d1/2), 2) + (math.Cos(r1)*math.Cos(r2))*math.Pow(math.Sin(d2/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EARTH_RADIUS * c
}

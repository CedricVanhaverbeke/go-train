package geojson

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

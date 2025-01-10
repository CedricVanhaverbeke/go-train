package gpx

import (
	"encoding/xml"
	"math"
	"time"
)

var XSI = `@xmlns:xsi:http://www.w3.org/2001/XMLSchema-instance
@xsi:schemaLocation:http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd
@creator:StravaGPX
@version:1.1
@xmlns:http://www.topografix.com/GPX/1/1
@xmlns:gpxtpx:http://www.garmin.com/xmlschemas/TrackPointExtension/v1
@xmlns:gpxx:http://www.garmin.com/xmlschemas/GpxExtensions/v3
`

var VIRTUAL_RIDE = "VirtualRide"

type metadata struct {
	Text string `xml:",chardata"`
	Time string `xml:"time"`
}

type trk struct {
	Text   string `xml:",chardata"`
	Name   string `xml:"name"`
	Type   string `xml:"type"`
	Trkseg struct {
		Text  string  `xml:",chardata"`
		Trkpt []trkpt `xml:"trkpt"`
	} `xml:"trkseg"`
}

type Gpx struct {
	XMLName        xml.Name `xml:"gpx"`
	Text           string   `xml:",chardata"`
	Xsi            string   `xml:"xsi,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	Creator        string   `xml:"creator,attr"`
	Version        string   `xml:"version,attr"`
	Xmlns          string   `xml:"xmlns,attr"`
	Gpxtpx         string   `xml:"gpxtpx,attr"`
	Gpxx           string   `xml:"gpxx,attr"`
	Metadata       metadata `xml:"metadata"`
	Trk            trk      `xml:"trk"`
}

func New(name string) Gpx {
	return Gpx{
		Xsi: XSI,

		Metadata: metadata{
			Time: time.Now().Format(time.RFC3339),
		},
		Trk: trk{
			Name: name,
			Type: VIRTUAL_RIDE,
		},
	}
}

func (gpx *Gpx) AddTrackpoint(trackPoint trkpt) {
	gpx.Trk.Trkseg.Trkpt = append(gpx.Trk.Trkseg.Trkpt, trackPoint)
}

// Distance returns the distance of
// a geojson file in meters
func (g *Gpx) Distance() float64 {
	return g.distance(0, len(g.Trk.Trkseg.Trkpt))
}

// distance returns the distance of a segment
// with a given start and end index
//
// it uses the Haversine formula
func (g *Gpx) distance(x int, y int) float64 {
	if g.Trk.Trkseg.Trkpt == nil {
		return 0.0
	}

	if len(g.Trk.Trkseg.Trkpt) < x ||
		len(g.Trk.Trkseg.Trkpt) < y {
		return 0.0
	}

	var d float64
	for i := x; i < y-1; i++ {
		c1 := g.Trk.Trkseg.Trkpt[i]
		c2 := g.Trk.Trkseg.Trkpt[i+1]
		d += haversine(c1.Lon, c1.Lat, c2.Lon, c2.Lat)
	}
	return d
}

var EARTH_RADIUS = 6371e3

// stolen from here https://www.movable-type.co.uk/scripts/latlong.html
func haversine(lng1 float64, lat1 float64, lng2 float64, lat2 float64) float64 {
	r1 := lat1 * (math.Pi / 180)
	r2 := lat2 * (math.Pi / 180)

	d1 := (lat2 - lat1) * (math.Pi / 180)
	d2 := (lng2 - lng1) * (math.Pi / 180)

	a := math.Pow(math.Sin(d1/2), 2) + (math.Cos(r1)*math.Cos(r2))*math.Pow(math.Sin(d2/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EARTH_RADIUS * c
}

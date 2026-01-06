package gpx

import (
	"encoding/xml"
	"fmt"
	"math"
	"time"

	"overlay/internal/angle"
	"overlay/internal/physics"
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
	return g.distance(0, len(g.Trk.Trkseg.Trkpt)-1)
}

// distance returns the distance of a segment
// with a given start and end index, end index is included
//
// it uses the Haversine formula
func (g *Gpx) distance(i int, j int) float64 {
	if g.Trk.Trkseg.Trkpt == nil {
		return 0.0
	}

	if len(g.Trk.Trkseg.Trkpt) < i ||
		len(g.Trk.Trkseg.Trkpt) < j {
		return 0.0
	}

	var d float64
	for z := i; z <= j-1; z++ {
		c1 := g.Trk.Trkseg.Trkpt[z]
		c2 := g.Trk.Trkseg.Trkpt[z+1]
		d += haversine(c1.Lon, c1.Lat, c2.Lon, c2.Lat)
	}
	return d
}

// slope returns the slope between two points of
// the gpx. it returns it in degrees
func (g *Gpx) Slope(i int, j int) float64 {
	if g.Trk.Trkseg.Trkpt == nil {
		return 0.0
	}

	if len(g.Trk.Trkseg.Trkpt) < i ||
		len(g.Trk.Trkseg.Trkpt) < j {
		return 0.0
	}

	var s float64
	for z := i; z <= j-1; z++ {
		c1 := g.Trk.Trkseg.Trkpt[z]
		c2 := g.Trk.Trkseg.Trkpt[z+1]

		el := c2.Ele - c1.Ele
		distance := g.distance(i, j)

		// make a right triangle, the slope in degrees is
		// tan(alpha) = el / distance
		tans := el / distance
		curs := math.Atan(tans)

		s += curs
	}

	return s / float64(j-i)
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

// Speed determines the initial speed of a route
// this is important to determine the progression
// the biker has made in a route. or something like that
// I'll see
func (g *Gpx) Speed(distance float64, power int) float64 {
	// Get slope for the current distance
	_, _, _, i, j := g.CoordInfo(distance)
	slope := g.Slope(i, j)
	fmt.Println("slope: ", angle.ToDegrees(slope))

	return physics.CalculateSpeed(float64(power), slope)
}

// CoordInfo returns lat/lng coordinates based on the driven distance
func (g *Gpx) CoordInfo(distance float64) (lat float64, lng float64, ele float64, i int, j int) {
	i = 1

	if distance == 0.0 {
		return g.Trk.Trkseg.Trkpt[0].Lat, g.Trk.Trkseg.Trkpt[0].Lon, g.Trk.Trkseg.Trkpt[0].Ele, 0, 1
	}

	distance = math.Mod(distance, g.Distance()) // Ensure distance wraps correctly

	// Special case: If distance matches total track distance, return last point
	if distance == g.Distance() {
		lastIndex := len(g.Trk.Trkseg.Trkpt) - 1
		return g.Trk.Trkseg.Trkpt[lastIndex].Lat,
			g.Trk.Trkseg.Trkpt[lastIndex].Lon,
			g.Trk.Trkseg.Trkpt[lastIndex].Ele,
			lastIndex - 1, lastIndex
	}

	// Find segment where distance fits
	for g.distance(0, i) < distance {
		i++
	}

	segmentD := g.distance(i-1, i)

	// Ensure segment distance is valid
	if segmentD == 0 {
		return g.Trk.Trkseg.Trkpt[i-1].Lat, g.Trk.Trkseg.Trkpt[i-1].Lon, g.Trk.Trkseg.Trkpt[i-1].Ele, i - 1, i
	}

	d := distance - g.distance(0, i-1)
	percentage := d / segmentD

	// Interpolate lat/lon/ele
	pt1 := g.Trk.Trkseg.Trkpt[i-1]
	pt2 := g.Trk.Trkseg.Trkpt[i]

	latD := (pt2.Lat - pt1.Lat) * percentage
	lngD := (pt2.Lon - pt1.Lon) * percentage
	eleD := (pt2.Ele - pt1.Ele) * percentage

	return pt1.Lat + latD, pt1.Lon + lngD, pt1.Ele + eleD, i - 1, i
}

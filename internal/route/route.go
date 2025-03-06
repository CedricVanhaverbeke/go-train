package route

import (
	"encoding/xml"
	"io"
	"overlay/pkg/gpx"
	"strings"
)

// createRoundTrip will make sure the created
// route is a roundTrip
func createRoundTrip(in gpx.Gpx) gpx.Gpx {
	i := len(in.Trk.Trkseg.Trkpt) - 1
	start := in.Trk.Trkseg.Trkpt[0]
	end := in.Trk.Trkseg.Trkpt[i]

	if start.Lat == end.Lat && start.Lon == end.Lon {
		return in
	}

	in.Trk.Trkseg.Trkpt = append(in.Trk.Trkseg.Trkpt, gpx.NewTrackpoint(start.Lat, start.Lon))
	return in
}

func New() gpx.Gpx {
	return gpx.Gpx{}
}

func NewFromFile(reader io.Reader) (gpx.Gpx, error) {
	var g gpx.Gpx
	err := xml.NewDecoder(reader).Decode(&g)
	if err != nil {
		return gpx.Gpx{}, err
	}

	g = createRoundTrip(g)

	return g, nil
}

func NewExample() gpx.Gpx {
	// example should never panic
	g, _ := NewFromFile(strings.NewReader(example))
	return g
}

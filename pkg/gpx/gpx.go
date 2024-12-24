package gpx

import (
	"encoding/xml"
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
		Text  string `xml:",chardata"`
		Trkpt []struct {
			Text       string `xml:",chardata"`
			Lat        string `xml:"lat,attr"`
			Lon        string `xml:"lon,attr"`
			Ele        string `xml:"ele"`
			Time       string `xml:"time"`
			Extensions struct {
				Text                string `xml:",chardata"`
				Power               string `xml:"power"`
				TrackPointExtension struct {
					Text string `xml:",chardata"`
					Cad  string `xml:"cad"`
				} `xml:"TrackPointExtension"`
			} `xml:"extensions"`
		} `xml:"trkpt"`
	} `xml:"trkseg"`
}

// Gpx was generated 2024-12-24 16:45:28 by https://xml-to-go.github.io/ in Ukraine.
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

func NewGpx(name string) Gpx {
	return Gpx{
		Xsi: XSI,

		Metadata: metadata{
			Time: time.Now().Format(time.RFC3339),
		},
		Trk: trk{
			Name: name,
			Type: VIRTUAL_RIDE,

			// todo: add trk seg
		},
	}
}

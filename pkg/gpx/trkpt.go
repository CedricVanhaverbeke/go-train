package gpx

import "time"

type trkpt struct {
	Text       string  `xml:",chardata"`
	Lat        float64 `xml:"lat,attr"`
	Lon        float64 `xml:"lon,attr"`
	Ele        float64 `xml:"ele"`
	Time       string  `xml:"time"`
	Extensions struct {
		Text                string `xml:",chardata"`
		Power               int    `xml:"power"`
		TrackPointExtension struct {
			Text string `xml:",chardata"`
			Hr   int    `xml:"hr,omitempty"`
			Cad  int    `xml:"cad,omitempty"`
		} `xml:"TrackPointExtension"`
	} `xml:"extensions"`
}

type trkOpt = func(trkpt *trkpt)

func NewTrackpoint(lat float64, lng float64, opts ...trkOpt) trkpt {
	pt := trkpt{}

	time := time.Now().Format(time.RFC3339)
	pt.Time = time
	pt.Lat = lat
	pt.Lon = lng

	for _, opt := range opts {
		opt(&pt)
	}

	return pt
}

func WithPower(power int) trkOpt {
	return func(tp *trkpt) {
		tp.Extensions.Power = power
	}
}

func WithCadence(cad int) trkOpt {
	return func(tp *trkpt) {
		tp.Extensions.TrackPointExtension.Cad = cad
	}
}

func WithHr(hr int) trkOpt {
	return func(tp *trkpt) {
		tp.Extensions.TrackPointExtension.Hr = hr
	}
}

func WithElevation(el float64) trkOpt {
	return func(tp *trkpt) {
		tp.Ele = el
	}
}

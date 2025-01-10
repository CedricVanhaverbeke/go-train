package route

import (
	"encoding/xml"
	"overlay/pkg/gpx"
	"strings"
)

func New() gpx.Gpx {
	var g gpx.Gpx
	err := xml.NewDecoder(strings.NewReader(example)).Decode(&g)
	if err != nil {
		panic(err) // panic for now, should never happen in this case
	}
	return g
}

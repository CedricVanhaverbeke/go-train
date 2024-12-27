package route

import (
	"encoding/json"
	"overlay/internal/geojson"
	"strings"
)

// TODO add elevation to this
var helloWorldRoute = `{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            3.739833247792518,
            51.020659989232115
          ],
          [
            3.752806393718231,
            51.033190690479756
          ],
          [
            3.742844870954258,
            51.037852412282376
          ],
          [
            3.746898979055686,
            51.019129845986896
          ],
          [
            3.7462039890951075,
            51.02729002642997
          ],
          [
            3.754775531939117,
            51.01956703492431
          ],
          [
            3.7606829466016904,
            51.02903844965229
          ],
          [
            3.7507214238377458,
            51.02743573088412
          ],
          [
            3.7623045898422447,
            51.0197127636539
          ],
          [
            3.7702969743866106,
            51.035303091278394
          ],
          [
            3.760451283282265,
            51.04032876124137
          ],
          [
            3.7702969743866106,
            51.01985849192542
          ],
          [
            3.778868517230592,
            51.04171254565961
          ],
          [
            3.7708761326862827,
            51.043169116185425
          ],
          [
            3.7785210222503167,
            51.020951439362676
          ],
          [
            3.7787526855709075,
            51.030131180674516
          ],
          [
            3.7829226253320485,
            51.03115103971808
          ],
          [
            3.786976733434642,
            51.030859653710024
          ],
          [
            3.7925366531169686,
            51.028164246285314
          ],
          [
            3.7925366531169686,
            51.02423012708846
          ],
          [
            3.7894091982954876,
            51.02182577876198
          ],
          [
            3.7859342484937883,
            51.02153433412781
          ],
          [
            3.7786368539100295,
            51.021680056673944
          ],
          [
            3.7940424646977817,
            51.0325350982387
          ],
          [
            3.798328236119829,
            51.02364726621752
          ],
          [
            3.8034248291621395,
            51.030641113001536
          ],
          [
            3.808868917184725,
            51.024448698025765
          ],
          [
            3.8107222237446763,
            51.03624999830993
          ],
          [
            3.82218955809347,
            51.036031483020025
          ],
          [
            3.825780339554882,
            51.035303091278394
          ],
          [
            3.828212804416978,
            51.03355490438386
          ],
          [
            3.8307611009374796,
            51.03042257126259
          ],
          [
            3.8307611009374796,
            51.02729002642997
          ],
          [
            3.8279811410963305,
            51.02561439228592
          ],
          [
            3.8199887565531867,
            51.024448698025765
          ],
          [
            3.8150079951706175,
            51.024448698025765
          ],
          [
            3.81373384691031,
            51.026780057233225
          ],
          [
            3.812691361969428,
            51.03005833274105
          ],
          [
            3.8139655102297354,
            51.03377343135247
          ],
          [
            3.816166311771184,
            51.03493889111394
          ],
          [
            3.818251281651726,
            51.035303091278394
          ],
          [
            3.8253170129148657,
            51.035812966699694
          ],
          [
            3.8308769325972207,
            51.043533251661245
          ],
          [
            3.8356260306604213,
            51.026561497288384
          ],
          [
            3.8307611009374796,
            51.043533251661245
          ],
          [
            3.8455875534243944,
            51.04156688608833
          ],
          [
            3.834351882398863,
            51.03523025147433
          ],
          [
            3.850220819826717,
            51.028674200253135
          ],
          [
            3.8633097974121426,
            51.04367890505006
          ],
          [
            3.851958294726927,
            51.046664698600665
          ],
          [
            3.860993164212175,
            51.02911129918877
          ],
          [
            3.877557091599357,
            51.031369578022066
          ],
          [
            3.87431380511822,
            51.04105707398213
          ],
          [
            3.8643522823530247,
            51.03974610280525
          ],
          [
            3.867827232154724,
            51.02998548469347
          ],
          [
            3.8776729232590412,
            51.031296732035116
          ],
          [
            3.868753885434728,
            51.058897159851796
          ],
          [
            3.879642061479842,
            51.02532297146908
          ],
          [
            3.747246474041674,
            51.01643375644721
          ],
          [
            3.7399490794578867,
            51.02058712641298
          ],
          [
            3.739833247792518,
            51.020659989232115
          ]
        ],
        "type": "LineString"
      }
    }
  ]
}
`

type Route struct {
	Geojson geojson.GeoJson `json:"geojson"`
}

// New loads the helloworld route for now
func New() Route {
	var g geojson.GeoJson
	err := json.NewDecoder(strings.NewReader(helloWorldRoute)).Decode(&g)
	if err != nil {
		panic(err) // panic for now, should never happen in this case
	}
	return Route{
		Geojson: g,
	}
}

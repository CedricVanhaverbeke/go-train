package state

type Metrics struct {
	Ftp   int
	Power int
	Speed int // in m/h -> so 30 000m/h = 30km/u
	Hr    int
}

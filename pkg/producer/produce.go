package producer

type Metrics struct {
	Power chan int
	Hr    chan int
}

func NewMetrics(power chan int, hr chan int) Metrics {
	return Metrics{
		Power: power,
		Hr:    hr,
	}
}

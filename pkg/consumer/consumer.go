package consumer

type Metrics struct {
	Power chan int
}

func NewMetrics(powerChan chan int) Metrics {
	return Metrics{
		Power: powerChan,
	}
}

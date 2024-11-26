package trainer

import (
	"fmt"
	"math/rand"
	"overlay/pkg/consumer"
	"overlay/pkg/producer"
	"time"
)

type RandomTrainer struct {
	cM producer.Metrics
}

func NewRandomTrainer(metr producer.Metrics) RandomTrainer {
	return RandomTrainer{
		cM: metr,
	}

}

func (rt RandomTrainer) Consume(metr consumer.Metrics) {
	for p := range metr.Power {
		fmt.Printf("Trainer should adjust to %dW\n", p)
	}
}

func (rt RandomTrainer) Produce() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		go func() {
			rt.cM.Power <- getPower()
			rt.cM.Hr <- getHeartRate()
		}()
	}

}

func getPower() int {
	return rand.Intn(200) + 100
}

func getHeartRate() int {
	return rand.Intn(100) + 80
}

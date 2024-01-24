package networksimulator

import (
	"math/rand"
	"time"
)

//type networksimulator struct{}

func NetworkDelay() {

	rand.Seed(time.Now().UnixNano())
	//1: randomNumber := rand.Intn(100)
	//2: randomNumber := rand.Intn(150)
	randomNumber := rand.Intn(200)

	// 1:
	//if rand.Float64() < 0.7 {
	//	randomNumber += 50
	//}
	//2
	//if rand.Float64() < 0.6 {
	//	randomNumber += 80
	//}

	if rand.Float64() < 0.5 {
		randomNumber += 80
	}

	time.Sleep(time.Duration(randomNumber) * time.Millisecond)
}

package application

import (
	"fmt"
	"math/rand"
	"time"
)

type ResiliencyService struct {
}

func (r *ResiliencyService) GenerateResiliency(minDelaySecond, maxDelaySecond int32, statusCodes []uint32) (string, uint32) {
	delay := rand.Intn(int(maxDelaySecond-minDelaySecond+1) + int(minDelaySecond))
	delaySecond := time.Duration(delay) * time.Second
	time.Sleep(delaySecond)
	idx := rand.Intn(len(statusCodes))
	str := fmt.Sprintf("The time now is %s, execution delayed for %d seconds", time.Now().Format("15:04:05.000"), delay)
	return str, statusCodes[idx]
}

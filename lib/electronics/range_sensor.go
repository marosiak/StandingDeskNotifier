package electronics

import (
	"errors"
	"github.com/shanghuiyang/rpi-devices/dev"
)

type RangeSensor struct {
	hcsr04 *dev.HCSR04
}

func NewRangeSensor(triggerPin, echoPin int8) *RangeSensor {
	return &RangeSensor{hcsr04: dev.NewHCSR04(triggerPin, echoPin)}
}

func (r *RangeSensor) DistInCM() (float32, error) {
	dist, err := r.hcsr04.Dist()
	if err != nil {
		return 0, err
	}
	if dist > 350 {
		return 0, errors.New("invalid distance")
	}
	return float32(dist), nil
}

func (r *RangeSensor) Dist64InCM() (float64, error) {
	return r.hcsr04.Dist()
}

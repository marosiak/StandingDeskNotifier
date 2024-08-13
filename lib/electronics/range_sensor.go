package electronics

import "github.com/shanghuiyang/rpi-devices/dev"

type RangeSensor struct {
	hcsr04 *dev.HCSR04
}

func NewRangeSensor(triggerPin, echoPin int8) *RangeSensor {
	return &RangeSensor{hcsr04: dev.NewHCSR04(17, 27)}
}

func (r *RangeSensor) DistInCM() (float32, error) {
	dist, err := r.hcsr04.Dist()
	if err != nil {
		return 0, err
	}
	return float32(dist), nil
}

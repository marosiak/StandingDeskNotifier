package electronics

import "github.com/shanghuiyang/rpi-devices/dev"

type Buzzer struct {
	buzzer dev.Buzzer
}

func NewBuzzer(pin uint8) *Buzzer {
	return &Buzzer{
		buzzer: dev.NewBuzzerImp(pin, dev.High),
	}
}

func (s *Buzzer) Beep(n int, intervalMs int) {
	s.buzzer.Beep(n, intervalMs)
}

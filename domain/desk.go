package domain

import (
	clock_lib "DeskNotifier/lib/clock"
	"time"
)

type Desk struct {
	TimeSpentUp   *time.Duration
	TimeSpentDown *time.Duration

	updatedAt       time.Time
	distanceToFloor float32

	topPosition    float32
	midPoint       float32
	bottomPosition float32
}

func NewDesk(bottomPosition float32, topPosition float32) *Desk {
	return &Desk{
		TimeSpentUp:    new(time.Duration),
		TimeSpentDown:  new(time.Duration),
		bottomPosition: bottomPosition,
		topPosition:    topPosition,
		midPoint:       (topPosition + bottomPosition) / 2,
		updatedAt:      time.Now(),
	}

}
func (d *Desk) IsLoaded() bool {
	return d.distanceToFloor != 0
}
func (d *Desk) IsHigh() bool {
	return d.distanceToFloor >= d.midPoint
}

func (d *Desk) IsLow() bool {
	return d.distanceToFloor < d.midPoint
}
func (d *Desk) UpdateCurrentPosition(currentPosition float32) {
	clock := clock_lib.Get()

	d.distanceToFloor = currentPosition
	since := clock.Since(d.updatedAt)
	if d.IsHigh() {
		*d.TimeSpentUp += since
	} else if d.IsLow() {
		*d.TimeSpentDown += since
	}
	d.updatedAt = clock.Now()
}

func (d *Desk) GetTimeSpentUp() time.Duration {
	return *d.TimeSpentUp
}

func (d *Desk) GetTimeSpentDown() time.Duration {
	return *d.TimeSpentDown
}

func (d *Desk) GetTimeUntilStand(targetDuration time.Duration) time.Duration {
	return d.GetTimeSpentDown() - targetDuration
}

func (d *Desk) GetTimeUntilSit(targetDuration time.Duration) time.Duration {
	return d.GetTimeSpentUp() - targetDuration
}

func (d *Desk) ResetTimeRecords() {
	d.TimeSpentUp = new(time.Duration)
	d.TimeSpentDown = new(time.Duration)
	d.updatedAt = time.Now()
}

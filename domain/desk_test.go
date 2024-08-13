package domain

import (
	"DeskNotifier/lib/clock"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDeskTimeSpentDown(t *testing.T) {
	bottomPosition := float32(95)
	topPosition := float32(120)

	desk := NewDesk(bottomPosition, topPosition)

	type positionDuration struct {
		position float32
		duration time.Duration
	}

	type test struct {
		name                 string
		positions            []positionDuration
		expectedDurationUp   time.Duration
		expectedDurationDown time.Duration
	}

	tests := []test{
		{
			name: "up and down",
			positions: []positionDuration{
				{position: bottomPosition, duration: time.Minute * 1},
				{position: topPosition, duration: time.Minute * 5},
				{position: bottomPosition, duration: time.Minute * 2},
			},
			expectedDurationUp:   time.Minute * 5,
			expectedDurationDown: time.Minute * 3,
		},
		{
			name: "not at exact up position",
			positions: []positionDuration{
				{position: topPosition, duration: time.Minute * 1},
				{position: topPosition + 5, duration: time.Minute * 1},
				{position: topPosition - 5, duration: time.Minute * 1},
			},
			expectedDurationUp: time.Minute * 3,
		},
		{
			name: "top position extremely high",
			positions: []positionDuration{
				{position: topPosition, duration: time.Minute * 1},
				{position: topPosition + 10000, duration: time.Minute * 1},
			},
			expectedDurationUp: time.Minute * 2,
		},
		{
			name: "bottom extremely low",
			positions: []positionDuration{
				{position: bottomPosition, duration: time.Minute * 1},
				{position: bottomPosition - 100, duration: time.Minute * 1},
			},
			expectedDurationDown: time.Minute * 2,
		},
	}

	fakeClock := clockwork.NewFakeClock()
	clock.InitMock(fakeClock)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, pd := range tt.positions {
				fakeClock.Advance(pd.duration)
				desk.UpdateCurrentPosition(pd.position)
			}

			assert.Equal(t, int(tt.expectedDurationUp.Minutes()), int(desk.TimeSpentUp.Minutes()))
			assert.Equal(t, int(tt.expectedDurationDown.Minutes()), int(desk.TimeSpentDown.Minutes()))
		})
	}
}

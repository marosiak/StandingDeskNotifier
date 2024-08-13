package main

import (
	"DeskNotifier/config"
	"DeskNotifier/domain"
	"DeskNotifier/lib/electronics"
	"log/slog"
	"time"
)

type App struct {
	cfg         *config.Config
	desk        *domain.Desk
	rangeSensor *electronics.RangeSensor
	speaker     *electronics.Buzzer
}

func NewApp() *App {
	cfg := config.NewConfig()
	return &App{
		cfg:         cfg,
		desk:        domain.NewDesk(cfg.DeskBottomPosition, cfg.DeskTopPosition),
		rangeSensor: electronics.NewRangeSensor(cfg.RangeSensorTriggerPin, cfg.RangeSensorEchoPin),
		speaker:     electronics.NewBuzzer(cfg.BuzzerPin),
	}
}

func (app *App) run() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			app.checkDeskStatus()
		}
	}
}

func (app *App) checkDeskStatus() {
	distance, err := app.rangeSensor.DistInCM()
	if err != nil {
		slog.Error("reading distance", slog.Any("error", err))
		return
	}

	app.desk.UpdateCurrentPosition(distance)

	if app.desk.IsLow() {
		app.handleSittingTooLong()
	} else {
		app.handleStandingTooLong()
	}
}

func (app *App) handleSittingTooLong() {
	secondsSittingTooLong := (app.desk.GetTimeSpentDown() - app.cfg.DurationToSit).Seconds()

	switch {
	case secondsSittingTooLong > 25:
		app.speaker.Beep(4, 3000)
	case secondsSittingTooLong > 0:
		app.speaker.Beep(1, 100)
		time.Sleep(time.Second * 25)
	}
}

func (app *App) handleStandingTooLong() {
	if app.desk.GetTimeSpentDown().Minutes() > app.cfg.DurationToStand.Minutes() {
		if app.cfg.NotifyToSit {
			app.speaker.Beep(1, 100)
		}
		app.desk.ResetTimeRecords()
	}
}

func main() {
	app := NewApp()
	app.run()
}

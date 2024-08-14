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
	slog.Info("App is running")

	for {
		select {
		case <-ticker.C:
			distance, err := app.rangeSensor.DistInCM()
			if err != nil {
				slog.Error("reading distance", slog.Any("error", err))
				continue
			}
			slog.Info("distance read", slog.Any("distance", distance))
			app.checkDeskStatus(distance)
		}
	}
}

func (app *App) checkDeskStatus(distance float32) {
	app.desk.UpdateCurrentPosition(distance)

	if app.desk.IsLow() {
		app.handleSittingTooLong()
	} else {
		app.handleStandingTooLong()
	}
}

func (app *App) handleSittingTooLong() {
	secondsSittingTooLong := (app.desk.GetTimeSpentDown() - app.cfg.DurationToSit.Duration()).Seconds()

	if secondsSittingTooLong <= 0 {
		return
	}
	app.speaker.Beep(1, 300)
	if secondsSittingTooLong < 45 {
		time.Sleep(time.Second * 50)
	}
}

func (app *App) handleStandingTooLong() {
	if app.desk.GetTimeSpentDown().Minutes() > app.cfg.DurationToStand.Duration().Minutes() {
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

package main

import (
	"DeskNotifier/config"
	"DeskNotifier/domain"
	"DeskNotifier/lib/electronics"
	"DeskNotifier/server"
	"log/slog"
	"time"
)

type App struct {
	cfg         *config.Config
	desk        *domain.Desk
	rangeSensor *electronics.RangeSensor
	speaker     *electronics.Buzzer

	firstChanceToStandUpCommunicatedAt *time.Time
}

func NewApp() *App {
	cfg := config.Get()
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
			slog.Info("rangeSensor", slog.Any("distance", distance))
			app.checkDeskStatus(distance)
		}
	}
}

func (app *App) checkDeskStatus(distance float32) {
	app.desk.UpdateCurrentPosition(distance)

	if app.desk.IsLow() {
		app.handleSittingTooLong()
	} else {
		app.handleStandingEnough()
	}
}

func (app *App) RegisterFirstSignalToStandUp() {
	tmp := time.Now()
	app.firstChanceToStandUpCommunicatedAt = &tmp
	app.speaker.Beep(1, 350)

}

func (app *App) handleSittingTooLong() {
	secondsSittingTooLong := (app.desk.GetTimeSpentDown() - app.cfg.DurationToSit.Duration()).Seconds()

	if secondsSittingTooLong <= 0 {
		return
	}

	if app.firstChanceToStandUpCommunicatedAt == nil {
		app.RegisterFirstSignalToStandUp()
		return
	}

	if time.Since(*app.firstChanceToStandUpCommunicatedAt) > time.Second*45 {
		app.speaker.Beep(1, 450)
	}
}

func (app *App) handleStandingEnough() {
	if app.desk.GetTimeSpentUp().Minutes() > app.cfg.DurationToStand.Duration().Minutes() {
		if app.cfg.NotifyToSit {
			app.speaker.Beep(1, 100)
		}
		slog.Info("Restarting timers because stood enugh time")
		app.desk.ResetTimeRecords()
		app.firstChanceToStandUpCommunicatedAt = nil
	}
}

func main() {
	app := NewApp()
	if app.cfg.HttpServerEnabled {
		go app.run()
		server.New(app.desk).Start()
	} else {
		app.run()
	}
}

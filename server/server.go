package server

import (
	"DeskNotifier/config"
	"DeskNotifier/domain"
	"DeskNotifier/templates"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"net/http"
	"time"

	"log/slog"
)

type Server struct {
	fiberApp *fiber.App
	desk     *domain.Desk
}

func New(desk *domain.Desk) *Server {
	engine := html.NewFileSystem(http.FS(templates.FS), ".gohtml")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	api := &Server{
		fiberApp: app,
		desk:     desk,
	}
	app.Get("/", api.Home)
	return api
}

func (a *Server) Start() {
	err := a.fiberApp.Listen(":8080")
	if err != nil {
		slog.Error("starting server", slog.Any("error", err))
		return
	}
}

func (a *Server) Home(c *fiber.Ctx) error {
	cfg := config.Get()
	return c.Render("home", fiber.Map{
		"Title":                 "DeskNotifier",
		"IsStanding":            a.desk.IsHigh(),
		"IsLoaded":              a.desk.IsLoaded(),
		"StandingDuration":      formatDuration(a.desk.GetTimeSpentUp()),
		"SittingDuration":       formatDuration(a.desk.GetTimeSpentDown()),
		"RemainingSittingTime":  formatDuration(a.desk.GetTimeUntilStand(cfg.DurationToSit.Duration()).Abs()),
		"RemainingStandingTime": formatDuration(a.desk.GetTimeUntilSit(cfg.DurationToStand.Duration()).Abs()),
		"RefreshInterval":       cfg.AutoRefreshPageDelayMs,
	})
}

func formatDuration(d time.Duration) string {
	rounded := d.Round(time.Second)
	minutes := int(rounded / time.Minute)
	seconds := int(rounded % time.Minute / time.Second)

	if seconds > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%dm", minutes)
}

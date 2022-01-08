package main

import (
	"task-scheduler/internal/api"
	"task-scheduler/internal/configs"
	"task-scheduler/internal/emailService"
	"task-scheduler/internal/platform/datastore"
	"task-scheduler/internal/platform/logger"
	"task-scheduler/internal/server/http"
	"task-scheduler/internal/tasks"
	"task-scheduler/internal/users"

	"github.com/joho/godotenv"
)

func main() {
	l := logger.New("task-scheduler", "v1.0.0", 1)

	err := godotenv.Load()
	if err != nil {
		l.Fatal("Error loading .env file")
	}
	cfg, err := configs.NewService()
	if err != nil {
		l.Fatal(err.Error())
		return
	}

	dscfg, err := cfg.Datastore()
	if err != nil {
		l.Fatal(err.Error())
		return
	}

	pqdriver, err := datastore.NewService(dscfg)
	if err != nil {
		l.Fatal(err.Error())
		return
	}

	us, err := users.NewService(l, pqdriver)
	if err != nil {
		l.Fatal(err.Error())
		return
	}

	ts, err := tasks.NewService(l, pqdriver)
	if err != nil {
		l.Fatal(err.Error())
		return
	}

	es, err := emailService.NewService(l, pqdriver)
	if err != nil {
		l.Fatal(err.Error())
		return
	}

	a, err := api.NewService(l, us, ts, es)
	if err != nil {
		l.Fatal(err.Error())
		return
	}

	httpCfg, err := cfg.HTTP()
	if err != nil {
		l.Fatal(err.Error())
		return
	}

	h, err := http.NewService(
		httpCfg,
		a,
	)
	if err != nil {
		l.Fatal(err.Error())
		return
	}

	h.Start()
}

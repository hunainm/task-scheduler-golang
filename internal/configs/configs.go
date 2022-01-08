package configs

import (
	"os"
	"strings"
	"time"

	"task-scheduler/internal/platform/datastore"
	"task-scheduler/internal/server/http"
)

type Configs struct {
}

func (cfg *Configs) HTTP() (*http.Config, error) {
	return &http.Config{
		TemplatesBasePath: strings.TrimSpace(os.Getenv("TEMPLATES_BASEPATH")),
		Port:              "8080",
		ReadTimeout:       time.Second * 5,
		WriteTimeout:      time.Second * 5,
		DialTimeout:       time.Second * 3,
	}, nil
}

func (cfg *Configs) Datastore() (*datastore.Config, error) {
	return &datastore.Config{
		Host:   os.Getenv("DBHOST"),
		Port:   "5432",
		Driver: "postgres",

		StoreName: "mjonszhl",
		Username:  os.Getenv("DBUSER"),
		Password:  os.Getenv("DBPASS"),

		SSLMode: "",

		ConnPoolSize: 10,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		IdleTimeout:  time.Second * 60,
		DialTimeout:  time.Second * 10,
	}, nil
}

func NewService() (*Configs, error) {
	return &Configs{}, nil
}

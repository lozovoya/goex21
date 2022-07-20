package main

import (
	"GoEx21/internal/api/httpserver"
	v1 "GoEx21/internal/api/httpserver/v1"
	"GoEx21/internal/domain/usecase/company"
	"GoEx21/internal/repository/postgres"
	"context"
	"github.com/caarlos0/env"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
)

type Params = struct {
	Port    string `env:"PORT" envDefault:"8888"`
	Host    string `env:"HOST" envDefault:"0.0.0.0"`
	DSN     string `env:"DSN" envDefault:"postgres://app:pass@localhost:5433/goex21"`
	User    string `env:"User" envDefault:"user"`
	Pass    string `env:"HOST" envDefault:"pass"`
	Country string `env:"COUNTRY" envDefault:"Cyprus"`
}

func main() {
	var config Params
	err := env.Parse(&config)
	if err != nil {
		log.Printf("Config load utils: %v", err)
		os.Exit(1)
	}
	if err = execute(config); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(config Params) (err error) {
	logrus := log.New()
	logrus.SetFormatter(&log.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.ReportCaller = true
	lg := logrus.WithFields(log.Fields{
		"app": "GoEx21",
	})

	companyPool, err := pgxpool.Connect(context.Background(), config.DSN)
	if err != nil {
		return err
	}
	companyRepo := postgres.NewCompanyRepo(companyPool)
	companyUsecase := company.NewCompanyUsecase(companyRepo)
	companyController := v1.NewCompanyController(companyUsecase, lg)

	var creds = map[string]string{config.User: config.Pass}
	router := httpserver.NewRouter(chi.NewRouter(), companyController, creds, config.Country, lg)
	server := http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: &router,
	}
	return server.ListenAndServe()
}

package main

import (
	"GoEx21/app/api/httpserver"
	v1 "GoEx21/app/api/httpserver/v1"
	"GoEx21/app/domain/usecase/company"
	"GoEx21/app/repository/postgres"
	"context"
	"github.com/getsentry/sentry-go"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type Params = struct {
	Port    string `env:"PORT" envDefault:"8888"`
	Host    string `env:"HOST" envDefault:"0.0.0.0"`
	DSN     string `env:"DSN" envDefault:"postgres://app:pass@localhost:5433/goex21"`
	User    string `env:"USER" envDefault:"user"`
	Pass    string `env:"PASS" envDefault:"pass"`
	Country string `env:"COUNTRY" envDefault:"Cyprus"`
	AmqpUrl string `env:"AMQP_URL" envDefault:"amqp://quest:quest@localhost:5672/"`
}

func main() {
	var config Params
	err := env.Parse(&config)
	if err != nil {
		log.Printf("Config load utils: %v", err)
		os.Exit(1)
	}
	if err = execute(&config); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(config *Params) (err error) {
	logrus := log.New()
	logrus.SetFormatter(&log.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.ReportCaller = true
	lg := logrus.WithFields(log.Fields{
		"app": "GoEx21",
	})

	err = sentry.Init(sentry.ClientOptions{
		Dsn: "https://f1bf877bcf344d61b6aa4a8ba88b660b@o4504002857664512.ingest.sentry.io/4504002971303936",
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	sentry.CaptureMessage("It started!")

	// commented for demo
	//amqpConn, err := amqp.Dial(config.AmqpUrl)
	//if err != nil {
	//	lg.Infof("Tried to connect to Rabbit")
	//}
	//amqpCh, err := amqpConn.Channel()
	//if err != nil {
	//	lg.Infof("Tried to connect to Rabbit")
	//}
	//defer amqpCh.Close()
	//remind erService := rabbitmq.NewCompanyEvent(amqpCh)

	companyPool, err := pgxpool.Connect(context.Background(), config.DSN)
	if err != nil {
		return err
	}
	companyRepo := postgres.NewCompanyRepo(companyPool)
	companyUsecase := company.NewCompanyUsecase(companyRepo, nil)
	CompanyHandlers := v1.NewCompanyController(companyUsecase, lg)

	var creds = map[string]string{config.User: config.Pass}
	router := httpserver.NewRouter(chi.NewRouter(), CompanyHandlers, creds, config.Country, lg)
	server := http.Server{
		Addr:              net.JoinHostPort(config.Host, config.Port),
		Handler:           &router,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}
	return server.ListenAndServe()
}

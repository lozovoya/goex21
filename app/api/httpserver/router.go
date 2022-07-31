package httpserver

import (
	v1 "GoEx21/app/api/httpserver/v1"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

func NewRouter(
	mux *chi.Mux,
	companyController *v1.CompanyController,
	creds map[string]string,
	contry string,
	lg *logrus.Entry) chi.Mux {
	mux.Use(middleware.Logger)
	mux.Route("/api/v1", func(router chi.Router) {
		companyRouter(router, companyController, creds, contry, lg)
	})
	lg.Infof("Router is started")
	return *mux
}

func companyRouter(
	router chi.Router,
	controller *v1.CompanyController,
	creds map[string]string,
	country string,
	lg *logrus.Entry) {
	router.
		//With(middleware.BasicAuth("goex21", creds), mw.CountryChecker(country, lg)).
		Post("/companies", controller.AddCompany)
	router.
		Get("/companies", controller.SearchCompany)
	router.
		Put("/companies", controller.EditCompany)
	router.
		//With(middleware.BasicAuth("goex21", creds), mw.CountryChecker(country, lg)).
		Delete("/companies", controller.DeleteCompany)
}

package httpserver

import (
	"GoEx21/internal/api/httpserver/mw"
	v1 "GoEx21/internal/api/httpserver/v1"

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
	router.With(
		middleware.BasicAuth("goex21", creds),
		mw.CountryChecker(country, lg)).Post(
		"/company/add", controller.AddCompany)
	router.Get("/company/search", controller.SearchCompany)
	router.Post("/company/edit", controller.EditCompany)
	router.With(
		middleware.BasicAuth("goex21", creds),
		mw.CountryChecker(country, lg)).Get(
		"/company/delete", controller.DeleteCompany)
}

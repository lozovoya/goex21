package v1

import (
	"GoEx21/internal/domain/model"
	"GoEx21/internal/domain/usecase"
	"GoEx21/internal/utils"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

type CompanyController struct {
	company usecase.CompanyUsecaseInterface
	lg      *logrus.Entry
}

func NewCompanyController(company usecase.CompanyUsecaseInterface, lg *logrus.Entry) *CompanyController {
	return &CompanyController{company: company, lg: lg}
}

func (c *CompanyController) AddCompany(writer http.ResponseWriter, request *http.Request) {
	var data *model.Company
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		c.lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	result, err := c.company.AddCompany(request.Context(), data)
	if (err != nil) || (len(result) == 0) {
		c.lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var reply = struct {
		Length    int `json:"length"`
		Companies []model.Company
	}{
		Length:    len(result),
		Companies: result,
	}
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(reply)
	if err != nil {
		c.lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (c *CompanyController) SearchCompany(writer http.ResponseWriter, request *http.Request) {
	values := request.URL.Query()
	conditions, err := c.getParameters(values)
	if (err != nil) && !errors.Is(err, utils.ErrNoParameters) {
		c.lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	reply, err := c.company.SearchCompany(request.Context(), conditions)
	if err != nil {
		if errors.Is(err, utils.ErrNoRecords) {
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		c.lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(reply)
	if err != nil {
		c.lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}

func (c *CompanyController) EditCompany(writer http.ResponseWriter, request *http.Request) {
	values := request.URL.Query()
	conditions, err := c.getParameters(values)
	if err != nil {
		c.lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	var data *model.Company
	err = json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		c.lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	result, err := c.company.EditCompany(request.Context(), conditions, data)
	if (err != nil) || (len(result) == 0) {
		if errors.Is(err, utils.ErrNoRecords) {
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		c.lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var reply = struct {
		Length    int `json:"length"`
		Companies []model.Company
	}{
		Length:    len(result),
		Companies: result,
	}
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(reply)
	if err != nil {
		c.lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (c *CompanyController) DeleteCompany(writer http.ResponseWriter, request *http.Request) {
	values := request.URL.Query()
	conditions, err := c.getParameters(values)
	if err != nil {
		c.lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	result, err := c.company.DeleteCompany(request.Context(), conditions)
	if (err != nil) || (len(result) == 0) {
		if errors.Is(err, utils.ErrNoRecords) {
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		c.lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var reply = struct {
		Length    int `json:"length"`
		Companies []model.Company
	}{
		Length:    len(result),
		Companies: result,
	}
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(reply)
	if err != nil {
		c.lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (c *CompanyController) getParameters(values url.Values) (map[string]string, error) {
	var conditions = make(map[string]string)
	if len(values) == 0 {
		return conditions, utils.ErrNoParameters
	}
	for k, v := range values {
		conditions[k] = v[0]
	}
	return conditions, nil
}

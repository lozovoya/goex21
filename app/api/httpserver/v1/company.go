package v1

import (
	"GoEx21/app/domain/model"
	"GoEx21/app/domain/usecase"
	"GoEx21/app/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/url"
	"strconv"

	"github.com/sirupsen/logrus"
)

type CompanyHandlers struct {
	company usecase.CompanyUsecaseInterface
	lg      *logrus.Entry
}

func NewCompanyController(company usecase.CompanyUsecaseInterface, lg *logrus.Entry) *CompanyHandlers {
	return &CompanyHandlers{company: company, lg: lg}
}

func (c *CompanyHandlers) AddCompany(writer http.ResponseWriter, request *http.Request) {
	data, err := c.unmarshalCompany(request)
	if err != nil {
		c.lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	result, err := c.company.AddCompany(request.Context(), data)
	if err != nil {
		c.lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var reply = struct {
		ID int64 `json:"id"`
	}{
		ID: result.ID,
	}
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(reply)
	if err != nil {
		c.lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (c *CompanyHandlers) SearchCompany(writer http.ResponseWriter, request *http.Request) {
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

func (c *CompanyHandlers) GetCompanyByID(writer http.ResponseWriter, request *http.Request) {
	companyID, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil {
		c.lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	reply, err := c.company.GetCompanyByID(request.Context(), companyID)
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

func (c *CompanyHandlers) EditCompany(writer http.ResponseWriter, request *http.Request) {
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

func (c *CompanyHandlers) EditCompanyByID(writer http.ResponseWriter, request *http.Request) {
	companyID, err := strconv.Atoi(chi.URLParam(request, "id"))
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
	result, err := c.company.EditCompanyByID(request.Context(), companyID, data)
	if err != nil {
		if errors.Is(err, utils.ErrNoRecords) {
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		c.lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var reply = struct {
		Company *model.Company
	}{
		Company: result,
	}
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(reply)
	if err != nil {
		c.lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (c *CompanyHandlers) DeleteCompany(writer http.ResponseWriter, request *http.Request) {
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

func (c *CompanyHandlers) DeleteCompanyByID(writer http.ResponseWriter, request *http.Request) {
	companyID, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil {
		c.lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	result, err := c.company.DeleteCompanyByID(request.Context(), companyID)
	if err != nil {
		if errors.Is(err, utils.ErrNoRecords) {
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		c.lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var reply = struct {
		Company *model.Company
	}{
		Company: result,
	}
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(reply)
	if err != nil {
		c.lg.Error(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

//func (c *CompanyHandlers) getParameters(values url.Values) (map[string]string, error) {
//	var conditions = make(map[string]string)
//	if len(values) == 0 {
//		return conditions, utils.ErrNoParameters
//	}
//	for k, v := range values {
//		conditions[k] = v[0]
//	}
//	return conditions, nil
//}

func (c *CompanyHandlers) getParameters(values url.Values) (*model.Conditions, error) {
	var err error
	var conditions model.Conditions
	if len(values) == 0 {
		return &conditions, utils.ErrNoParameters
	}
	if v, ok := values["name"]; ok {
		conditions.Name.Value = v[0]
		conditions.Name.IsExist = true
	}
	if v, ok := values["code"]; ok {
		conditions.Code.Value = v[0]
		conditions.Code.IsExist = true
	}
	if v, ok := values["country"]; ok {
		conditions.Country.Value = v[0]
		conditions.Country.IsExist = true
	}
	if v, ok := values["website"]; ok {
		conditions.Website.Value = v[0]
		conditions.Website.IsExist = true
	}
	if v, ok := values["phone"]; ok {
		conditions.Phone.Value = v[0]
		conditions.Phone.IsExist = true
	}
	if v, ok := values["isactive"]; ok {
		conditions.IsActive.Value, err = strconv.ParseBool(v[0])
		if err != nil {
			return &conditions, fmt.Errorf("v1.getParameters: %w", err)
		}
		conditions.Code.IsExist = true
	}
	return &conditions, nil
}

func (c *CompanyHandlers) unmarshalCompany(request *http.Request) (*model.Company, error) {
	var data *model.Company
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		c.lg.Error(err)
		return data, fmt.Errorf("v1.unmarshalCompany: %w", err)
	}
	return data, err
}

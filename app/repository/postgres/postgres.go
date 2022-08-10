package postgres

import (
	"GoEx21/app/domain/model"
	"GoEx21/app/repository"
	"GoEx21/app/utils"
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type companyRepo struct {
	pool *pgxpool.Pool
}

const ALL_COMPANY_FIELDS = "id, name, code, country, website, phone, isactive"

func NewCompanyRepo(pool *pgxpool.Pool) repository.CompanyRepositoryInterface {
	return &companyRepo{pool: pool}
}

func (c *companyRepo) AddCompany(ctx context.Context, company *model.Company) (*model.Company, error) {
	dbReq := `INSERT 
			  INTO companies (name, code, country, website, phone) 
			  VALUES ($1, $2, $3, $4, $5) 
			  RETURNING id`
	err := c.pool.QueryRow(ctx,
		dbReq,
		company.Name,
		company.Code,
		company.Country,
		company.Website,
		company.Phone).Scan(&company.ID)
	if err != nil {
		return company, fmt.Errorf("postgres.AddCompany: %w", err)
	}
	return company, nil
}

//func (c *companyRepo) SearchCompany(ctx context.Context, conditions map[string]string) ([]model.Company, error) {
//	conditionsReq := c.conditionsBuilder(conditions)
//	dbReq := "SELECT " + ALL_COMPANY_FIELDS + " FROM companies WHERE " + conditionsReq + " isactive=true"
//	var companies []model.Company
//	companies, err := c.getCompany(ctx, dbReq)
//	if err != nil {
//		return companies, fmt.Errorf("postgres.SearchCompany: %w", err)
//	}
//	return companies, nil
//}

func (c *companyRepo) SearchCompany(ctx context.Context, conditions *model.Conditions) ([]model.Company, error) {
	var companies []model.Company
	conditionsReq, err := c.conditionsBuilder(conditions)
	if err != nil {
		return companies, fmt.Errorf("postgres.SearchCompany: %w", err)
	}
	dbReq := "SELECT " + ALL_COMPANY_FIELDS + " FROM companies WHERE " + conditionsReq + " isactive=true"
	companies, err = c.getCompany(ctx, dbReq)
	if err != nil {
		return companies, fmt.Errorf("postgres.SearchCompany: %w", err)
	}
	return companies, nil
}

func (c *companyRepo) EditCompany(
	ctx context.Context,
	conditions *model.Conditions,
	company *model.Company) ([]model.Company, error) {

	var result []model.Company
	conditionsReq, err := c.conditionsBuilder(conditions)
	if err != nil {
		return result, fmt.Errorf("postgres.EditCompany: %w", err)
	}
	dbReq := "UPDATE companies SET "
	setters, err := c.setBuilder(company)
	if err != nil {
		return result, fmt.Errorf("postgres.EditCompany: %w", err)
	}
	dbReq = fmt.Sprintf(" %s %s WHERE %s isactive=true RETURNING %s",
		dbReq, setters, conditionsReq, ALL_COMPANY_FIELDS)
	result, err = c.getCompany(ctx, dbReq)
	if err != nil {
		return result, fmt.Errorf("postgres.EditCompany: %w", err)
	}
	return result, nil
}

func (c *companyRepo) SetCompanyInactive(ctx context.Context, conditions *model.Conditions) ([]model.Company, error) {
	var result []model.Company
	conditionsReq, err := c.conditionsBuilder(conditions)
	if err != nil {
		return result, fmt.Errorf("postgres.SetCompanyInactive: %w", err)
	}
	dbReq := fmt.Sprintf(`UPDATE companies 
								 SET isactive=false 
								 WHERE  %s isactive=true 
								 RETURNING id`, conditionsReq)
	result, err = c.getCompany(ctx, dbReq)
	if err != nil {
		return result, fmt.Errorf("postgres.SetCompanyInActive: %w", err)
	}
	return result, nil
}

func (c *companyRepo) GetCompanyByID(ctx context.Context, id int) (*model.Company, error) {
	var company model.Company
	dbReq := "SELECT " + ALL_COMPANY_FIELDS + " FROM companies WHERE id=$1 AND isactive=true"
	err := c.pool.QueryRow(ctx, dbReq, id).Scan(&company.ID,
		&company.Name,
		&company.Code,
		&company.Country,
		&company.Website,
		&company.Phone,
		&company.IsActive)
	if err != nil {
		return &company, fmt.Errorf("postgres.GetCompanyByID: %w", err)
	}
	return &company, nil
}

func (c *companyRepo) EditCompanyByID(ctx context.Context, id int, company *model.Company) (*model.Company, error) {
	//TODO implement me
	panic("implement me")
}

func (c *companyRepo) SetCompanyInactiveCompanyByID(ctx context.Context, id int) (*model.Company, error) {
	//TODO implement me
	panic("implement me")
}

func (c *companyRepo) getCompany(ctx context.Context, dbReq string, args ...interface{}) ([]model.Company, error) {
	var companies = make([]model.Company, 0)
	rows, err := c.pool.Query(ctx, dbReq, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return companies, utils.ErrNoRecords
		}
		return companies, fmt.Errorf("postgres.getCompany: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var company model.Company
		err = rows.Scan(&company.ID,
			&company.Name,
			&company.Code,
			&company.Country,
			&company.Website,
			&company.Phone,
			&company.IsActive)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return companies, utils.ErrNoRecords
			}
			return companies, fmt.Errorf("postgres.getCompany: %w", err)
		}
		companies = append(companies, company)
	}
	return companies, nil
}

//func (c *companyRepo) conditionsBuilder(conditions map[string]string) string {
//	result := " "
//	for k, v := range conditions {
//		result = fmt.Sprintf(" %s %s='%s' AND ", result, k, v) //!!!!! STOP очень плохо, sql injection
//	}
//	return result
//}

func (c *companyRepo) conditionsBuilder(conditions *model.Conditions) (string, error) {
	result := " "
	if conditions.Name.IsExist {
		validParameter := regexp.MustCompile(`^[\w- ]+$`)
		if !validParameter.MatchString(conditions.Name.Value) {
			return result, fmt.Errorf("postgres.conditionsBuilder: wrong query parameter")
		}
		result = fmt.Sprintf(" %s names='%s' AND ", result, conditions.Name.Value)
	}
	if conditions.Code.IsExist {
		validParameter := regexp.MustCompile(`^[\w- ]+$`)
		if !validParameter.MatchString(conditions.Code.Value) {
			return result, fmt.Errorf("postgres.conditionsBuilder: wrong query parameter")
		}
		result = fmt.Sprintf(" %s code='%s' AND ", result, conditions.Code.Value)
	}
	if conditions.Country.IsExist {
		validParameter := regexp.MustCompile(`^[\w- ]+$`)
		if !validParameter.MatchString(conditions.Country.Value) {
			return result, fmt.Errorf("postgres.conditionsBuilder: wrong query parameter")
		}
		result = fmt.Sprintf(" %s country='%s' AND ", result, conditions.Country.Value)
	}
	if conditions.Website.IsExist {
		validParameter := regexp.MustCompile(`^[\w-.]+$`)
		if !validParameter.MatchString(conditions.Website.Value) {
			return result, fmt.Errorf("postgres.conditionsBuilder: wrong query parameter")
		}
		result = fmt.Sprintf(" %s website='%s' AND ", result, conditions.Website.Value)
	}
	if conditions.Phone.IsExist {
		validParameter := regexp.MustCompile(`^[\d-+]+$`)
		if !validParameter.MatchString(conditions.Phone.Value) {
			return result, fmt.Errorf("postgres.conditionsBuilder: wrong query parameter")
		}
		result = fmt.Sprintf(" %s phone='%s' AND ", result, conditions.Phone.Value)
	}
	if conditions.IsActive.IsExist {
		result = fmt.Sprintf(" %s isactive='%s' AND ", result, conditions.IsActive.Value)
	}
	return result, nil
}

func (c *companyRepo) setBuilder(company *model.Company) (string, error) {
	setters := " "
	if company.Name != "" {
		validParameter := regexp.MustCompile(`^[\w- ]+$`)
		if !validParameter.MatchString(company.Name) {
			return setters, fmt.Errorf("postgres.setBuilder: wrong query parameter")
		}
		setters = fmt.Sprintf(" %s name='%s', ", setters, company.Name)
	}
	if company.Code != "" {
		validParameter := regexp.MustCompile(`^[\w- ]+$`)
		if !validParameter.MatchString(company.Code) {
			return setters, fmt.Errorf("postgres.setBuilder: wrong query parameter")
		}
		setters = fmt.Sprintf(" %s code='%s', ", setters, company.Code)
	}
	if company.Country != "" {
		validParameter := regexp.MustCompile(`^[\w- ]+$`)
		if !validParameter.MatchString(company.Country) {
			return setters, fmt.Errorf("postgres.setBuilder: wrong query parameter")
		}
		setters = fmt.Sprintf(" %s country='%s', ", setters, company.Country)
	}
	if company.Website != "" {
		validParameter := regexp.MustCompile(`^[\w-.]+$`)
		if !validParameter.MatchString(company.Website) {
			return setters, fmt.Errorf("postgres.setBuilder: wrong query parameter")
		}
		setters = fmt.Sprintf(" %s website='%s', ", setters, company.Website)
	}
	if company.Phone != "" {
		validParameter := regexp.MustCompile(`^[\d-+]+$`)
		if !validParameter.MatchString(company.Phone) {
			return setters, fmt.Errorf("postgres.setBuilder: wrong query parameter")
		}
		setters = fmt.Sprintf(" %s phone='%s', ", setters, company.Phone)
	}
	result := setters + " modified=CURRENT_TIMESTAMP "
	return result, nil
}

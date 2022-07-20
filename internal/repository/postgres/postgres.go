package postgres

import (
	"GoEx21/internal/domain/model"
	"GoEx21/internal/repository"
	"GoEx21/internal/utils"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type companyRepo struct {
	pool *pgxpool.Pool
}

const ALL_COMPANY_FIELDS = "id, name, code, country, website, phone"

func NewCompanyRepo(pool *pgxpool.Pool) repository.CompanyRepositoryInterface {
	return &companyRepo{pool: pool}
}

func (c *companyRepo) AddCompany(ctx context.Context, company *model.Company) (int64, error) {
	dbReq := `INSERT INTO companies (name, code, country, website, phone, isactive) 
			VALUES ($1, $2, $3, $4, $5, $6) 
			RETURNING id`
	var id int64
	err := c.pool.QueryRow(ctx, dbReq,
		company.Name, company.Code, company.Country, company.Website, company.Phone, true).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("postgres.AddCompany: %w", err)
	}
	return id, nil
}

func (c *companyRepo) SearchCompany(ctx context.Context, conditions map[string]string) ([]model.Company, error) {
	conditionsReq := c.conditionsBuilder(conditions)
	dbReq := "SELECT " + ALL_COMPANY_FIELDS + " FROM companies WHERE " + conditionsReq + " isactive=true"
	var companies []model.Company
	companies, err := c.getCompany(ctx, dbReq)
	if err != nil {
		return companies, fmt.Errorf("postgres.SearchCompany: %w", err)
	}
	return companies, nil
}

func (c *companyRepo) EditCompany(
	ctx context.Context,
	conditions map[string]string,
	company *model.Company) error {
	conditionsReq := c.conditionsBuilder(conditions)
	dbReq := "UPDATE companies SET "
	setters := c.setBuilder(company)
	dbReq = fmt.Sprintf(" %s %s WHERE %s isactive=true", dbReq, setters, conditionsReq)
	result, err := c.pool.Exec(ctx, dbReq)
	if err != nil {
		return fmt.Errorf("postgres.EditCompany: %w", err)
	}
	if result.RowsAffected() == 0 {
		return utils.ErrNoRecords
	}
	return nil
}

func (c *companyRepo) SetCompanyInActive(ctx context.Context, conditions map[string]string) error {
	conditionsReq := c.conditionsBuilder(conditions)
	dbReq := "UPDATE companies SET isactive=false WHERE " + conditionsReq + " isactive=true"
	result, err := c.pool.Exec(ctx, dbReq)
	if err != nil {
		return fmt.Errorf("postgres.SetCompanyInActive: %w", err)
	}
	if result.RowsAffected() == 0 {
		return utils.ErrNoRecords
	}
	return nil
}

func (c *companyRepo) getCompany(ctx context.Context, dbReq string) ([]model.Company, error) {
	var companies = make([]model.Company, 0)
	rows, err := c.pool.Query(ctx, dbReq)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return companies, utils.ErrNoRecords
		}
		return companies, fmt.Errorf("postgres.getCompany: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var company model.Company
		err = rows.Scan(&company.ID, &company.Name, &company.Code, &company.Country, &company.Website, &company.Phone)
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

func (c *companyRepo) conditionsBuilder(conditions map[string]string) string {
	result := " "
	for k, v := range conditions {
		result = fmt.Sprintf(" %s %s='%s' AND ", result, k, v)
	}
	return result
}
func (c *companyRepo) setBuilder(company *model.Company) string {
	setters := " "
	if company.Name != "" {
		setters = fmt.Sprintf(" %s name='%s', ", setters, company.Name)
	}
	if company.Code != "" {
		setters = fmt.Sprintf(" %s code='%s', ", setters, company.Code)
	}
	if company.Country != "" {
		setters = fmt.Sprintf(" %s country='%s', ", setters, company.Country)
	}
	if company.Website != "" {
		setters = fmt.Sprintf(" %s website='%s', ", setters, company.Website)
	}
	if company.Phone != "" {
		setters = fmt.Sprintf(" %s phone='%s', ", setters, company.Phone)
	}
	result := setters + " modified=CURRENT_TIMESTAMP "
	return result
}

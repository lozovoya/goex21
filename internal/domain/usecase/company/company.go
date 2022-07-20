package company

import (
	"GoEx21/internal/domain/model"
	"GoEx21/internal/domain/usecase"
	"GoEx21/internal/reminder"
	"GoEx21/internal/repository"
	"context"
	"fmt"
)

type CompanyUsecase struct {
	repoCompany repository.CompanyRepositoryInterface
	reminder    reminder.ReminderInterface
}

func NewCompanyUsecase(
	repoCompany repository.CompanyRepositoryInterface,
	reminder reminder.ReminderInterface) usecase.CompanyUsecaseInterface {
	return &CompanyUsecase{repoCompany: repoCompany, reminder: reminder}
}

func (c *CompanyUsecase) AddCompany(ctx context.Context, company *model.Company) (int64, error) {
	id, err := c.repoCompany.AddCompany(ctx, company)
	if err != nil {
		return id, fmt.Errorf("company.AddCompany: %w", err)
	}
	return id, nil
}

func (c *CompanyUsecase) SearchCompany(ctx context.Context, conditions map[string]string) ([]model.Company, error) {
	companies, err := c.repoCompany.SearchCompany(ctx, conditions)
	if err != nil {
		return companies, fmt.Errorf("company.SearchCompany: %w", err)
	}
	return companies, nil
}

func (c *CompanyUsecase) EditCompany(ctx context.Context, conditions map[string]string, company *model.Company) error {
	err := c.repoCompany.EditCompany(ctx, conditions, company)
	if err != nil {
		return fmt.Errorf("company.EditCompany: %w", err)
	}
	return nil
}

func (c *CompanyUsecase) DeleteCompany(ctx context.Context, conditions map[string]string) error {
	err := c.repoCompany.SetCompanyInActive(ctx, conditions)
	if err != nil {
		return fmt.Errorf("company.DeleteCompany: %w", err)
	}
	return nil
}

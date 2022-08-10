package company

import (
	"GoEx21/app/domain/model"
	"GoEx21/app/domain/usecase"
	"GoEx21/app/reminder"
	"GoEx21/app/repository"
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

func (c *CompanyUsecase) AddCompany(ctx context.Context, company *model.Company) (*model.Company, error) {
	result, err := c.repoCompany.AddCompany(ctx, company)
	if (err != nil) || (result.ID == 0) {
		return result, fmt.Errorf("company.AddCompany: %w", err)
	}
	//commented for demo
	//err = c.reminder.SendEvent(result)
	//if err != nil {
	//	return result, fmt.Errorf("company.AddCompany: %w", err)
	//}
	fmt.Printf("Send message to rabbit: %+v \n", result)
	return result, nil
}

func (c *CompanyUsecase) SearchCompany(ctx context.Context, conditions *model.Conditions) ([]model.Company, error) {
	companies, err := c.repoCompany.SearchCompany(ctx, conditions)
	if err != nil {
		return companies, fmt.Errorf("company.SearchCompany: %w", err)
	}
	return companies, nil
}

func (c *CompanyUsecase) EditCompany(
	ctx context.Context,
	conditions *model.Conditions,
	company *model.Company) ([]model.Company, error) {
	result, err := c.repoCompany.EditCompany(ctx, conditions, company)
	if err != nil {
		return result, fmt.Errorf("company.EditCompany: %w", err)
	}
	//commented for demo
	//err = c.reminder.SendEvent(result)
	//if err != nil {
	//	return result, fmt.Errorf("company.EditCompany: %w", err)
	//}
	fmt.Printf("Send message to rabbit: %+v \n", result)
	return result, nil
}

func (c *CompanyUsecase) DeleteCompany(ctx context.Context, conditions *model.Conditions) ([]model.Company, error) {
	result, err := c.repoCompany.SetCompanyInactive(ctx, conditions)
	if err != nil {
		return result, fmt.Errorf("company.DeleteCompany: %w", err)
	}
	//commented for demo
	//err = c.reminder.SendEvent(result)
	//if err != nil {
	//	return result, fmt.Errorf("company.DeleteCompany: %w", err)
	//}
	fmt.Printf("Send message to rabbit: %+v \n", result)
	return result, nil
}

func (c *CompanyUsecase) GetCompanyByID(ctx context.Context, id int) (*model.Company, error) {
	companies, err := c.repoCompany.GetCompanyByID(ctx, id)
	if err != nil {
		return companies, fmt.Errorf("company.GetCompanyByID: %w", err)
	}
	return companies, nil
}

func (c *CompanyUsecase) EditCompanyByID(ctx context.Context, id int, company *model.Company) (*model.Company, error) {
	result, err := c.repoCompany.EditCompanyByID(ctx, id, company)
	if err != nil {
		return result, fmt.Errorf("company.EditCompanyByID: %w", err)
	}
	//commented for demo
	//err = c.reminder.SendEvent(result)
	//if err != nil {
	//	return result, fmt.Errorf("company.EditCompanyByID: %w", err)
	//}
	fmt.Printf("Send message to rabbit: %+v \n", result)
	return result, nil
}

func (c *CompanyUsecase) DeleteCompanyByID(ctx context.Context, id int) (*model.Company, error) {
	result, err := c.repoCompany.SetCompanyInactiveCompanyByID(ctx, id)
	if err != nil {
		return result, fmt.Errorf("company.DeleteCompanyByID: %w", err)
	}
	//commented for demo
	//err = c.reminder.SendEvent(result)
	//if err != nil {
	//	return result, fmt.Errorf("company.DeleteCompanyByID: %w", err)
	//}
	fmt.Printf("Send message to rabbit: %+v \n", result)
	return result, nil
}

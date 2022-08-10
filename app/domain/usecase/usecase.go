package usecase

import (
	"GoEx21/app/domain/model"
	"context"
)

type CompanyUsecaseInterface interface {
	AddCompany(ctx context.Context, company *model.Company) (*model.Company, error)
	SearchCompany(ctx context.Context, conditions *model.Conditions) ([]model.Company, error)
	GetCompanyByID(ctx context.Context, id int) (*model.Company, error)
	EditCompany(ctx context.Context, conditions *model.Conditions, company *model.Company) ([]model.Company, error)
	EditCompanyByID(ctx context.Context, id int, company *model.Company) (*model.Company, error)
	DeleteCompany(ctx context.Context, conditions *model.Conditions) ([]model.Company, error)
	DeleteCompanyByID(ctx context.Context, id int) (*model.Company, error)
}

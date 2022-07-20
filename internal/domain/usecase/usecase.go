package usecase

import (
	"GoEx21/internal/domain/model"
	"context"
)

type CompanyUsecaseInterface interface {
	AddCompany(ctx context.Context, company *model.Company) (int64, error)
	SearchCompany(ctx context.Context, conditions map[string]string) ([]model.Company, error)
	EditCompany(ctx context.Context, conditions map[string]string, company *model.Company) error
	DeleteCompany(ctx context.Context, conditions map[string]string) error
}

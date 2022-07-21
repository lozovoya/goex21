package reminder

import (
	"GoEx21/internal/domain/model"
)

type ReminderInterface interface {
	SendEvent(company []model.Company) error
}

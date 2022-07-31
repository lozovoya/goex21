package reminder

import (
	"GoEx21/app/domain/model"
)

type ReminderInterface interface {
	SendEvent(company []model.Company) error
}

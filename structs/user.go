package structs

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var validate *validator.Validate

type User struct {
	gorm.Model
	Name      string `json:"name" gorm:"not null" validate:"required"`
	BirthDate string `json:"birthDate" gorm:"not null" validate:"required,birthdate"`
	Gender    string `json:"gender" gorm:"not null" validate:"required"`
	Email     string `json:"email" gorm:"unique;not null" validate:"required,email"`
	Password  string `json:"password" gorm:"not null" validate:"required"`
	Role      string `json:"role" gorm:"default:member"`
}

func init() {
	validate = validator.New()
	validate.RegisterValidation("birthdate", validateBirthDate)
}

func validateBirthDate(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	birthDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}
	return birthDate.Before(time.Now())
}

func (u *User) Validate() error {
	return validate.Struct(u)
}

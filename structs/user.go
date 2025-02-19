package structs

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"gorm.io/gorm"
)

var (
	validate   *validator.Validate
	translator ut.Translator
)

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
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	translator, _ = uni.GetTranslator("en")

	enTranslations.RegisterDefaultTranslations(validate, translator)

	validate.RegisterValidation("birthdate", validateBirthDate)
	validate.RegisterTranslation("birthdate", translator, func(ut ut.Translator) error {
		return ut.Add("birthdate", "{0} must be before today", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("birthdate", fe.Field())
		return t
	})

	validate.RegisterTranslation("required", translator, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})
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
	err := validate.Struct(u)
	if err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, err.Translate(translator))
		}
		return fmt.Errorf(strings.Join(errorMessages, ", "))
	}
	return nil
}

package user

import (
	"net/mail"

	"github.com/go-playground/validator"
)

type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
	DeletedAt string `json:"-"`
}

func New() *User {
	return &User{}
}

func (u *User) Validate() error {
	v := validator.New()
	v.RegisterValidation("email", validateEmail)

	return v.Struct(u)
}

func validateEmail(fl validator.FieldLevel) bool {
	_, err := mail.ParseAddress(fl.Field().String())
	return err == nil
}

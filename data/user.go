package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
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

func (u *User) Validate() error{
	v := validator.New()
	v.RegisterValidation("email", validateEmail)

	return v.Struct(u)
}

func validateEmail(fl validator.FieldLevel) bool {
	_, err := mail.ParseAddress(fl.Field().String())
	return err == nil
}

func (u *User) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(u)
}

type Users []*User

func (u *Users) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

func GetUsers() Users {
	return userList
}

func AddUser(u *User) Users {
	u.Id = getNextId()
	userList = append(userList, u)
	return userList
}

func GetUser(id int) *User {
	return userList[id]
}

func UpdateUser(id int, u *User) error {
	pos, err := findUser(id)
	u.Id = id
	if err != nil {
		return err
	}
	userList[pos] = u
	return nil
}

func DeleteUser(id int) error {
	pos, err := findUser(id)
	if err != nil {
		return err
	}
	userList[pos] = userList[len(userList) - 1]
	userList = userList[:len(userList) - 1]
	return nil
}

var ErrUserNotFound = fmt.Errorf("User not found")

func findUser(id int) (int, error){
	for i , u := range userList {
		if u.Id == id {
			return i, nil
		}
	}
	return -1, ErrUserNotFound
}

func getNextId() int {
	last := userList[len(userList)-1].Id
	return last + 1
}

var userList = Users{
	{
		Id:        1,
		FirstName: "Dush",
		LastName:  "Desh",
		Email:     "dushdesh@email.com",
		CreatedAt: time.Now().UTC().String(),
		UpdatedAt: time.Now().UTC().String(),
		DeletedAt: "",
	},
	{
		Id:        2,
		FirstName: "Neha",
		LastName:  "Kumar",
		Email:     "neha@email.com",
		CreatedAt: time.Now().UTC().String(),
		UpdatedAt: time.Now().UTC().String(),
		DeletedAt: "",
	},
}

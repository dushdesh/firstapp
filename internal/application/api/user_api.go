package api

import (
	"github.com/dushdesh/firstapp/internal/ports"
	"github.com/dushdesh/firstapp/internal/application/core/user"
)

type UserApp struct {
	userRepo ports.UserRepository
}

func NewUserApp(repo ports.UserRepository) *UserApp {
	return &UserApp{userRepo: repo}
}

func (uapi *UserApp) CreateUser(u *user.User) error {
	err := uapi.userRepo.Create(u)
	return err
}

func (uapi *UserApp) UpdateUser(id int, u *user.User) error {
	err := uapi.userRepo.Update(id, u)
	return err
}

func (uapi *UserApp) DeleteUser(id int) error {
	err := uapi.userRepo.Delete(id)
	return err
}

func (uapi *UserApp) FindAllUsers() ([]*user.User, error) {
	users, err := uapi.userRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (uapi *UserApp) FindUser(id int) (*user.User, error) {
	user, err := uapi.userRepo.Find(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

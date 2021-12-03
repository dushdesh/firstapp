package ports

import (
	u "github.com/dushdesh/firstapp/internal/application/core/user"
)

type DBPort interface {
	CloseDBConnection()
}

type UserRepository interface {
	Create(*u.User) error
	Update(int, *u.User) error
	Delete(int) error
	FindAll() ([]*u.User, error)
	Find(int) (*u.User, error)
}

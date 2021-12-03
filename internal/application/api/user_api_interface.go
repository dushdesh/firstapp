package api

import "github.com/dushdesh/firstapp/internal/application/core/user"

type UserApi interface {
	Validate(*user.User) error
}



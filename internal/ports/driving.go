package ports

type UserAPIPort interface {
	Create(*user.User) error
	Update(*user.User) error
	Delete(int) error
	FindAll() ([]*user.User, error)
	Find(int) (*user.User, error)
}

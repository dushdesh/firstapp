package db

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	u "github.com/dushdesh/firstapp/internal/application/core/user"
)


func (da *Adapter) Create(user *u.User) error {
	query, args, err := sq.Insert("users").
		Columns("first_name", "last_name", "email", "created_at", "updated_at").
		Values(user.FirstName, user.LastName, user.Email, time.Now().UTC().String(), time.Now().UTC().String()).
		ToSql()

	if err != nil {
		return err
	}

	_, err = da.db.Exec(query, args...)

	if err != nil {
		return err
	}

	return nil
}

func (da *Adapter) Update(id int, user *u.User) error {
	query, args, err := sq.Update("users").
		Set("first_name", user.FirstName).
		Set("last_name", user.LastName).
		Set("email", user.Email).
		Set("updated_at", time.Now().UTC().String()).
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = da.db.Exec(query, args...)

	if err != nil {
		return err
	}

	return nil
}

func (da *Adapter) Delete(id int) error {
	query, args, err := sq.Delete("users").
		Where(sq.Eq{"id": id}).ToSql()

	if err != nil {
		return err
	}

	_, err = da.db.Exec(query, args...)

	if err != nil {
		return err
	}

	return nil
}

func (da *Adapter) FindAll() ([]*u.User, error) {
	query, _ , err := sq.Select("*").From("users").ToSql()

	var users []*u.User

	if err != nil {
		return users, err
	}

	rows, err := da.db.Query(query)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user *u.User

		err := rows.Scan(
			user.Id,
			user.FirstName,
			user.LastName,
			user.CreatedAt,
			user.UpdatedAt,
			user.DeletedAt,
		)
		if err != nil {
			return users, err
		}
		users = append(users, user)
		if err = rows.Err(); err != nil {
			return users, err
		}
	}

	return users, nil

}

func (da *Adapter) Find(id int) (*u.User, error) {

	var user *u.User

	sql, args, err := sq.Select("*").From("users").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}

    if err = da.db.QueryRow(sql, args...).Scan(&user); err != nil {
        return nil, err
    }

    return user, nil
}

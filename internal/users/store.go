package users

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/bnkamalesh/errors"
	"github.com/jackc/pgx/v4/pgxpool"
)

type store interface {
	Create(ctx context.Context, u *User) error
	GetUser(ctx context.Context, email string) (*User, error)
}

type userStore struct {
	qbuilder  squirrel.StatementBuilderType
	pqdriver  *pgxpool.Pool
	tableName string
}

func (us *userStore) Create(ctx context.Context, u *User) error {
	query, args, err := us.qbuilder.Insert(us.tableName).SetMap(map[string]interface{}{
		"fullName":  u.Name,
		"email":     u.Email,
		"pwd":       u.Password,
		"createdAt": u.CreatedAt,
		"updatedAt": u.UpdatedAt,
	}).ToSql()
	if err != nil {
		return errors.InternalErr(err, errors.DefaultMessage)
	}

	_, err = us.pqdriver.Exec(ctx, query, args...)
	if err != nil {
		return errors.InternalErr(err, "Email is already in use")
	}

	return nil
}

func (us *userStore) GetUser(ctx context.Context, email string) (*User, error) {
	query, args, err := us.qbuilder.Select(
		"id",
		"fullName",
		"pwd",
		"createdAt",
		"updatedAt",
	).From(
		us.tableName,
	).Where(
		squirrel.Eq{
			"email": email,
		},
	).ToSql()

	if err != nil {
		return nil, errors.InternalErr(err, errors.DefaultMessage)
	}

	user := new(User)
	id := new(sql.NullInt64)
	fullname := new(sql.NullString)
	pwd := new(sql.NullString)

	row := us.pqdriver.QueryRow(ctx, query, args...)
	err = row.Scan(
		id,
		fullname,
		pwd,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, errors.InternalErr(err, err.Error())
	}

	user.UID = id.Int64
	user.Name = fullname.String
	user.Email = email
	user.Password = pwd.String

	return user, nil
}

func newStore(pqdriver *pgxpool.Pool) (*userStore, error) {
	return &userStore{
		pqdriver:  pqdriver,
		qbuilder:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		tableName: "Users",
	}, nil
}

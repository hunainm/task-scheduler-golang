package api

import (
	"context"

	"task-scheduler/internal/users"
)

func (a *API) Register(ctx context.Context, u *users.User) (*users.User, error) {
	u, err := a.users.Register(ctx, u)
	if err != nil {
		a.logger.Error(err)
		return nil, err
	}

	return u, nil
}

func (a *API) Login(ctx context.Context, email string, password string) (users.JWT, error) {
	token, err := a.users.Login(ctx, email, password)
	if err != nil {
		a.logger.Error(err)
		return token, err
	}
	return token, nil
}

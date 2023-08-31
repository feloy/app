package mock

import (
	"context"

	"github.com/0xdod/go-realworld/conduit"
)

type UserService struct {
	CreateUserFn   func(*conduit.User) error
	UserByEmailFn  func(string) *conduit.User
	AuthenticateFn func() *conduit.User
}

func (m *UserService) CreateUser(_ context.Context, user *conduit.User) error {
	return m.CreateUserFn(user)
}

func (m *UserService) UserByID(_ context.Context, id uint) (*conduit.User, error) {
	return nil, nil
}

func (m *UserService) UserByEmail(_ context.Context, email string) (*conduit.User, error) {
	return m.UserByEmailFn(email), nil
}

func (m *UserService) Authenticate(_ context.Context, email, password string) (*conduit.User, error) {
	return m.AuthenticateFn(), nil
}

func (us *UserService) UpdateUser(ctx context.Context, user *conduit.User, patch conduit.UserPatch) error {
	return nil
}

func (us *UserService) DeleteUser(ctx context.Context, id uint) error {
	return nil
}

func (us *UserService) Users(ctx context.Context, uf conduit.UserFilter) ([]*conduit.User, error) {
	return nil, nil
}

func (m *UserService) FollowUser(_ context.Context, user, follower *conduit.User) error {
	panic("not implemeted")
}

func (m *UserService) UnFollowUser(_ context.Context, user, follower *conduit.User) error {
	panic("not implmented")
}

func (us *UserService) UserByUsername(ctx context.Context, uname string) (*conduit.User, error) {
	panic("not implemented")
}

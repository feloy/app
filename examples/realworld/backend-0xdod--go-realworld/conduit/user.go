package conduit

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint      `json:"-"`
	Email        string    `json:"email,omitempty"`
	Username     string    `json:"username,omitempty"`
	Bio          string    `json:"bio,omitempty"`
	Image        string    `json:"image,omitempty"`
	Token        string    `json:"token,omitempty"`
	Following    []*User   `json:"-"`
	Followers    []*User   `json:"-"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"-" db:"created_at"`
	UpdatedAt    time.Time `json:"-" db:"updated_at"`
}

type Following struct {
	ID          uint  `json:"id"`
	FollowingID uint  `json:"following_id" db:"following_id"`
	Following   *User `json:"following"`
	FollowerID  uint  `json:"follower_id" db:"follower_id"`
	Follower    *User `json:"follower"`
}

type Profile struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

func (u *User) Profile() *Profile {
	return &Profile{
		Username: u.Username,
		Bio:      u.Bio,
		Image:    u.Image,
	}
}

func (u *User) ProfileWithFollow(_u *User) *Profile {
	return &Profile{
		Username:  u.Username,
		Bio:       u.Bio,
		Image:     u.Image,
		Following: _u.IsFollowing(u),
	}
}

func (me *User) IsFollowing(user *User) bool {
	for _, u := range user.Followers {
		if me.Username == u.Username {
			return true
		}
	}
	return false
}

var AnonymousUser User

type UserFilter struct {
	ID       *uint
	Email    *string
	Username *string

	Limit  int
	Offset int
}

type UserPatch struct {
	Email        *string `json:"email"`
	Username     *string `json:"username"`
	Image        *string `json:"image"`
	Bio          *string `json:"bio"`
	PasswordHash *string `json:"-" db:"password_hash"`
}

func (u *User) SetPassword(password string) error {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		// return better error message
		return err
	}

	u.PasswordHash = string(hashBytes)

	return nil
}

func (u User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))

	return err == nil
}

func (u *User) IsAnonymous() bool {
	return u == &AnonymousUser
}

type UserService interface {
	Authenticate(ctx context.Context, email, password string) (*User, error)

	CreateUser(context.Context, *User) error

	UserByID(context.Context, uint) (*User, error)

	UserByEmail(context.Context, string) (*User, error)

	UserByUsername(context.Context, string) (*User, error)

	Users(context.Context, UserFilter) ([]*User, error)

	UpdateUser(context.Context, *User, UserPatch) error

	// FollowUser creates a following relationship as follower follows user
	FollowUser(ctx context.Context, user, follower *User) error

	// UnFollowUser removes the following relationship as follower unfollows user
	UnFollowUser(ctx context.Context, user, follower *User) error

	DeleteUser(context.Context, uint) error
}

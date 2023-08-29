package conduit

import (
	"context"
	"time"
)

type Comment struct {
	ID            uint      `json:"id"`
	ArticleID     uint      `json:"-" db:"article_id"`
	Article       *Article  `json:"-"`
	AuthorID      uint      `json:"-" db:"author_id"`
	Author        *User     `json:"-"`
	AuthorProfile *Profile  `json:"author"`
	Body          string    `json:"body"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
}

func (c *Comment) SetAuthorProfile(currentUser *User) {
	c.AuthorProfile = &Profile{
		Username: c.Author.Username,
		Bio:      c.Author.Bio,
		Image:    c.Author.Image,
	}

	c.AuthorProfile.Following = currentUser.IsFollowing(c.Author)
}

type CommentFilter struct {
	ID        *uint
	ArticleID *uint
	AuthorID  *uint

	Limit  int
	Offset int
}

type CommentService interface {
	CreateComment(context.Context, *Comment) error
	CommentByID(context.Context, uint) (*Comment, error)
	Comments(context.Context, CommentFilter) ([]*Comment, error)
	DeleteComment(context.Context, uint) error
}

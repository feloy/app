package conduit

import (
	"context"
	"time"
)

type Article struct {
	ID             uint      `json:"-"`
	Title          string    `json:"title"`
	Body           string    `json:"body"`
	Description    string    `json:"description"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount int64     `json:"favoritesCount" db:"favorites_count"`
	FavoritedBy    []*User   `json:"-"`
	Slug           string    `json:"slug"`
	AuthorID       uint      `json:"-" db:"author_id"`
	Author         *User     `json:"-"`
	AuthorProfile  *Profile  `json:"author"`
	Tags           []*Tag    `json:"tagList"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

func (a *Article) SetAuthorProfile(currentUser *User) {
	a.AuthorProfile = &Profile{
		Username: a.Author.Username,
		Bio:      a.Author.Bio,
		Image:    a.Author.Image,
	}

	a.AuthorProfile.Following = currentUser.IsFollowing(a.Author)
}

func (a *Article) UserHasFavorite(currentUser *User) bool {
	for _, fav := range a.FavoritedBy {
		if fav.ID == currentUser.ID {
			return true
		}
	}

	return false
}

type ArticleFilter struct {
	ID             *uint
	Title          *string
	Description    *string
	AuthorID       *uint
	AuthorUsername *string
	Tag            *string
	Slug           *string
	FavoritedBy    *string

	Limit  int
	Offset int
}

type ArticlePatch struct {
	Title       *string
	Body        *string
	Description *string
	Slug        *string
	Tags        []Tag
}

type ArticleService interface {
	CreateArticle(context.Context, *Article) error
	ArticleBySlug(context.Context, string) (*Article, error)
	Articles(context.Context, ArticleFilter) ([]*Article, error)
	ArticleFeed(context.Context, *User, ArticleFilter) ([]*Article, error)
	FavoriteArticle(ctx context.Context, userID uint, article *Article) error
	UnfavoriteArticle(ctx context.Context, userID uint, article *Article) error
	UpdateArticle(context.Context, *Article, ArticlePatch) error
	DeleteArticle(context.Context, uint) error
}

func (a *Article) AddTags(_tags ...string) {
	for _, t := range _tags {
		a.Tags = append(a.Tags, &Tag{Name: t})
	}
}

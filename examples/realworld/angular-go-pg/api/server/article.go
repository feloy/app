package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/0xdod/go-realworld/conduit"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	goslug "github.com/gosimple/slug"
)

func (s *Server) createArticle() http.HandlerFunc {
	type Input struct {
		Article struct {
			Title       string   `json:"title" validate:"required"`
			Description string   `json:"description"`
			Body        string   `json:"body" validate:"required"`
			Tags        []string `json:"tagList"`
		} `json:"article"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		input := Input{}

		if err := readJSON(r.Body, &input); err != nil {
			badRequestError(w)
			return
		}

		if err := validate.Struct(input.Article); err != nil {
			validationError(w, err)
			return
		}

		article := conduit.Article{
			Title:       input.Article.Title,
			Body:        input.Article.Body,
			Slug:        slug.Make(input.Article.Title),
			Description: input.Article.Description,
		}

		article.AddTags(input.Article.Tags...)
		user := userFromContext(r.Context())
		article.Author = user

		if user.IsAnonymous() {
			invalidAuthTokenError(w)
			return
		}

		if err := s.articleService.CreateArticle(r.Context(), &article); err != nil {
			serverError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, M{"article": article})
	}
}

func (s *Server) listArticles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		filter := conduit.ArticleFilter{}

		if v := query.Get("author"); v != "" {
			filter.AuthorUsername = &v
		}

		if v := query.Get("tag"); v != "" {
			filter.Tag = &v
		}

		if v := query.Get("favorited"); v != "" {
			filter.FavoritedBy = &v
		}

		articles, err := s.articleService.Articles(r.Context(), filter)

		if err != nil {
			serverError(w, err)
			return
		}
		user := userFromContext(r.Context())
		for _, a := range articles {
			a.SetAuthorProfile(user)
			a.Favorited = a.UserHasFavorite(user)
		}

		writeJSON(w, http.StatusOK, M{"articles": articles})
	}
}

func (s *Server) articleFeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		filter := conduit.ArticleFilter{}
		limit, _ := strconv.Atoi(query.Get("limit"))
		filter.Limit = limit
		ctx := r.Context()
		articles, err := s.articleService.ArticleFeed(ctx, userFromContext(ctx), filter)

		if err != nil {
			serverError(w, err)
			return
		}

		for _, a := range articles {
			a.SetAuthorProfile(userFromContext(ctx))
		}

		writeJSON(w, http.StatusOK, M{"articles": articles})
	}
}

func (s *Server) getArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		filter := conduit.ArticleFilter{}

		if slug, exists := vars["slug"]; exists {
			filter.Slug = &slug
		}

		articles, err := s.articleService.Articles(r.Context(), filter)

		if err != nil {
			serverError(w, err)
			return
		}

		var article *conduit.Article

		if len(articles) > 0 {
			article = articles[0]
			article.SetAuthorProfile(userFromContext(r.Context()))
		}

		writeJSON(w, http.StatusOK, M{"article": article})
	}
}

func (s *Server) updateArticle() http.HandlerFunc {
	type Input struct {
		Article struct {
			Title       *string `json:"title,omitempty"`
			Description *string `json:"description,omitempty"`
			Body        *string `json:"body,omitempty"`
		} `json:"article,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		input := Input{}

		if err := readJSON(r.Body, &input); err != nil {
			badRequestError(w)
			return
		}

		slug := mux.Vars(r)["slug"]

		article, err := s.articleService.ArticleBySlug(r.Context(), slug)

		if err != nil {
			switch {
			case errors.Is(err, conduit.ErrNotFound):
				err := ErrorM{"article": []string{"requested article not found"}}
				notFoundError(w, err)
			default:
				serverError(w, err)
			}
			return
		}

		user := userFromContext(r.Context())

		if user.ID != article.AuthorID {
			err := ErrorM{"article": []string{"forbidden request"}}
			errorResponse(w, http.StatusForbidden, err)
			return
		}

		patch := conduit.ArticlePatch{
			Title:       input.Article.Title,
			Body:        input.Article.Body,
			Description: input.Article.Description,
		}

		if patch.Title != nil {
			*patch.Slug = goslug.Make(*patch.Title)
		}

		if err := s.articleService.UpdateArticle(r.Context(), article, patch); err != nil {
			serverError(w, err)
			return
		}

		article.SetAuthorProfile(user)

		writeJSON(w, http.StatusOK, M{"article": article})
	}
}

func (s *Server) listTags() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tags, err := s.tagService.Tags(r.Context(), conduit.TagFilter{})

		if err != nil {
			serverError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, M{"tags": tags})
	}
}

func (s *Server) deleteArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := mux.Vars(r)["slug"]

		article, err := s.articleService.ArticleBySlug(r.Context(), slug)

		if err != nil {
			switch {
			case errors.Is(err, conduit.ErrNotFound):
				err := ErrorM{"article": []string{"requested article not found"}}
				notFoundError(w, err)
			default:
				serverError(w, err)
			}
			return
		}

		user := userFromContext(r.Context())

		if user.ID != article.AuthorID {
			err := ErrorM{"article": []string{"forbidden request"}}
			errorResponse(w, http.StatusForbidden, err)
			return
		}

		if err := s.articleService.DeleteArticle(r.Context(), article.ID); err != nil {
			serverError(w, err)
			return
		}

		writeJSON(w, http.StatusNoContent, nil)
	}
}

func (s *Server) addComment() http.HandlerFunc {
	type Input struct {
		Comment struct {
			Body string `json:"body" validate:"required"`
		} `json:"comment" validate:"required"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		input := Input{}

		if err := readJSON(r.Body, &input); err != nil {
			badRequestError(w)
			return
		}

		if err := validate.Struct(input.Comment); err != nil {
			validationError(w, err)
			return
		}

		ctx := r.Context()
		user := userFromContext(ctx)
		slug := mux.Vars(r)["slug"]
		article, err := s.articleService.ArticleBySlug(ctx, slug)

		if err != nil {
			switch {
			case errors.Is(err, conduit.ErrNotFound):
				err := ErrorM{"article": []string{"requested article not found"}}
				notFoundError(w, err)
			default:
				serverError(w, err)
			}
			return
		}

		comment := conduit.Comment{
			Body:          input.Comment.Body,
			AuthorID:      user.ID,
			ArticleID:     article.ID,
			AuthorProfile: user.Profile(),
		}

		if err := s.commentService.CreateComment(ctx, &comment); err != nil {
			serverError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, M{"comment": comment})
	}
}

func (s *Server) favoriteAction(action string) http.HandlerFunc {
	const (
		Favorite   = "favorite"
		UnFavorite = "unfavorite"
	)
	return func(w http.ResponseWriter, r *http.Request) {
		slug := mux.Vars(r)["slug"]
		article, err := s.articleService.ArticleBySlug(r.Context(), slug)

		if err != nil {
			if errors.Is(err, conduit.ErrNotFound) {
				notFoundError(w, ErrorM{"article": []string{"article not found"}})
			} else {
				serverError(w, err)
			}
			return
		}

		user := userFromContext(r.Context())

		switch isFavorited := article.UserHasFavorite(user); action {
		case Favorite:
			if !isFavorited {
				err = s.articleService.FavoriteArticle(r.Context(), user.ID, article)
			}
		case UnFavorite:
			if isFavorited {
				err = s.articleService.UnfavoriteArticle(r.Context(), user.ID, article)
			}
		}

		if err != nil {
			serverError(w, err)
			return
		}

		article.SetAuthorProfile(user)
		writeJSON(w, http.StatusOK, M{"article": article})
	}
}

func (s *Server) getArticleComments() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		slug := mux.Vars(r)["slug"]
		article, err := s.articleService.ArticleBySlug(r.Context(), slug)

		if err != nil {
			if errors.Is(err, conduit.ErrNotFound) {
				notFoundError(w, ErrorM{"article": []string{"article not found"}})
			} else {
				serverError(w, err)
			}
			return
		}

		cf := conduit.CommentFilter{ArticleID: &article.ID}
		comments, err := s.commentService.Comments(r.Context(), cf)

		if err != nil {
			serverError(w, err)
			return
		}

		user := userFromContext(r.Context())

		for _, c := range comments {
			c.SetAuthorProfile(user)
		}

		writeJSON(w, http.StatusOK, M{"comments": comments})
	}
}

func (s *Server) deleteComment() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idInt, _ := strconv.Atoi(vars["comment_id"])

		comment, err := s.commentService.CommentByID(r.Context(), uint(idInt))

		if err != nil {
			if errors.Is(err, conduit.ErrNotFound) {
				notFoundError(w, ErrorM{"comment": []string{"comment not found"}})
			} else {
				serverError(w, err)
			}
			return
		}

		if comment.AuthorID != userFromContext(r.Context()).ID {
			errorResponse(w, http.StatusForbidden, "not permitted to delete comment")
			return
		}

		if err := s.commentService.DeleteComment(r.Context(), comment.ID); err != nil {
			serverError(w, err)
			return
		}

		writeJSON(w, http.StatusNoContent, nil)
	}
}

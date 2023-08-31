package server

import (
	"os"

	"github.com/rs/cors"
)

const (
	MustAuth     = true
	OptionalAuth = false
)

func (s *Server) routes() {
	s.router.Use(cors.AllowAll().Handler)
	s.router.Use(Logger(os.Stdout))
	apiRouter := s.router.PathPrefix("/api/v1").Subrouter()

	optionalAuth := apiRouter.PathPrefix("").Subrouter()
	optionalAuth.Use(s.authenticate(OptionalAuth))
	{
		optionalAuth.Handle("/profiles/{username}", s.getProfile()).Methods("GET")
		optionalAuth.Handle("/tags", s.listTags()).Methods("GET")
	}

	noAuth := apiRouter.PathPrefix("").Subrouter()
	{
		noAuth.Handle("/health", healthCheck())
		noAuth.Handle("/users", s.createUser()).Methods("POST")
		noAuth.Handle("/users/login", s.loginUser()).Methods("POST")
	}

	authApiRoutes := apiRouter.PathPrefix("").Subrouter()
	authApiRoutes.Use(s.authenticate(MustAuth))
	{
		authApiRoutes.Handle("/user", s.getCurrentUser()).Methods("GET")
		authApiRoutes.Handle("/user", s.updateUser()).Methods("PUT", "PATCH")
		authApiRoutes.Handle("/articles", s.createArticle()).Methods("POST")
		authApiRoutes.Handle("/articles", s.listArticles()).Methods("GET")
		authApiRoutes.Handle("/articles/feed", s.articleFeed()).Methods("GET")
		authApiRoutes.Handle("/articles/{slug}", s.getArticle()).Methods("GET")
		authApiRoutes.Handle("/articles/{slug}", s.updateArticle()).Methods("PUT", "PATCH")
		authApiRoutes.Handle("/articles/{slug}", s.deleteArticle()).Methods("DELETE")
		authApiRoutes.Handle("/articles/{slug}/comment", s.addComment()).Methods("POST")
		authApiRoutes.Handle("/articles/{slug}/comments", s.getArticleComments()).Methods("GET")
		authApiRoutes.Handle("/articles/{slug}/comments/{comment_id}", s.deleteComment()).Methods("DELETE")
		authApiRoutes.Handle("/articles/{slug}/favorite", s.favoriteAction("favorite")).Methods("POST")
		authApiRoutes.Handle("/articles/{slug}/favorite", s.favoriteAction("unfavorite")).Methods("DELETE")
		authApiRoutes.Handle("/profiles/{username}/follow", s.followAction("follow")).Methods("POST")
		authApiRoutes.Handle("/profiles/{username}/follow", s.followAction("unfollow")).Methods("DELETE")
	}
}

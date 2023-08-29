package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/0xdod/go-realworld/conduit"
	"github.com/0xdod/go-realworld/mock"
	"github.com/gorilla/mux"
)

func Test_createUser(t *testing.T) {
	userStore := &mock.UserService{}
	srv := testServer()
	srv.userService = userStore

	input := `{
		"user": {
			"email":    "e@mail.com",
			"username": "username",
			"password": "passwerd"
		}
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(input))
	w := httptest.NewRecorder()

	var user conduit.User
	userStore.CreateUserFn = func(u *conduit.User) error {
		user = *u
		return nil
	}

	srv.router.ServeHTTP(w, req)
	expectedResp := userResponse(&user)
	gotResp := M{}
	extractResponseBody(w.Body, &gotResp)

	if code := w.Code; code != http.StatusCreated {
		t.Errorf("expected status code of 201, but got %d", code)
	}

	if !reflect.DeepEqual(expectedResp, gotResp) {
		t.Errorf("expected response %v, but got %v", expectedResp, gotResp)
	}
}

func Test_loginUser(t *testing.T) {
	userStore := &mock.UserService{}
	srv := testServer()
	srv.userService = userStore

	userStore.AuthenticateFn = func() *conduit.User {
		user := &conduit.User{
			Email: "e@mail.com",
		}
		return user
	}

	input := `{
		"user": {
			"email":    "e@mail.com",
			"password": "passwerd"
		}
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", strings.NewReader(input))
	w := httptest.NewRecorder()

	srv.router.ServeHTTP(w, req)

	if code := w.Code; code != http.StatusOK {
		t.Errorf("expected status code of 200, but got %d", code)
	}
}

func testServer() *Server {
	srv := &Server{
		router: mux.NewRouter(),
	}
	srv.routes()
	return srv
}

func extractResponseBody(body io.Reader, v interface{}) {
	mm := M{}
	_ = readJSON(body, &mm)
	byt, err := json.Marshal(mm["user"])
	if err != nil {
		panic(err)
	}
	json.Unmarshal(byt, v)
}

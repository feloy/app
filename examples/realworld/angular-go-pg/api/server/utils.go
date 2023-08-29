package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/0xdod/go-realworld/conduit"
	"github.com/golang-jwt/jwt"
)

// M is a generic map
type M map[string]interface{}

func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	jsonBytes, err := json.Marshal(data)

	if err != nil {
		serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(jsonBytes)

	if err != nil {
		log.Println(err)
	}
}

func readJSON(body io.Reader, input interface{}) error {
	return json.NewDecoder(body).Decode(input)
}

var hmacSampleSecret = []byte("sample-secret")

func generateUserToken(user *conduit.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
	})

	tokenString, err := token.SignedString(hmacSampleSecret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func parseUserToken(tokenStr string) (userClaims M, err error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, conduit.ErrUnAuthorized
		}

		return hmacSampleSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, nil
	}

	return M(claims), nil
}

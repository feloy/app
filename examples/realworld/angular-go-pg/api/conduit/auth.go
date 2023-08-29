package conduit

// TODO: roll-out stateful auth tokens, as with the JWTs it's quite tricky to revoke tokens

// import (
// 	"context"
// 	"crypto/rand"
// 	"crypto/sha256"
// 	"encoding/base32"
// 	"time"
// )

// const (
// 	ScopeAuthorization = "authorization"
// )

// type Token struct {
// 	ID        uint
// 	PlainText string
// 	UserID    uint
// 	Hash      string
// 	Expiry    time.Time
// 	Scope     string
// }

// func generateToken(userID uint, ttl time.Duration, scope string) (*Token, error) {
// 	token := &Token{
// 		UserID: userID,
// 		Expiry: time.Now().Add(ttl),
// 		Scope:  scope,
// 	}

// 	randomBytes := make([]byte, 16)

// 	_, err := rand.Read(randomBytes)
// 	if err != nil {
// 		return nil, err
// 	}

// 	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
// 	hash := sha256.Sum256([]byte(token.PlainText))
// 	token.Hash = string(hash[:])
// 	return token, nil
// }

// type TokenService interface {
// 	CreateToken(context.Context, *Token) error
// 	DeleteToken(ctx context.Context, userID uint, scope string) error
// }

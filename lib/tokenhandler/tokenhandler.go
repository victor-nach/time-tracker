package tokenhandler

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	AuthTokenDuration    = 15 * time.Hour
	RefreshTokenDuration = 48 * time.Hour
)

var (
	ErrInvalidSigningMethod = errors.New("invalid token signing method")
	ErrInvalidToken         = errors.New("invalid token")
)

type Claims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

type TokenHandler interface {
	ValidateToken(token string) (*Claims, error)
	NewToken(userId string, expirationTime time.Time) (string, error)
}

type tokenHandler struct {
	jwtSecret string
}

// validate interface implementation
var _ TokenHandler = &tokenHandler{}

func New(secret string) TokenHandler {
	return &tokenHandler{
		jwtSecret: secret,
	}
}

//NewToken ...
func (t *tokenHandler) NewToken(userId string, expirationTime time.Time) (string, error) {
	claims := &Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(t.jwtSecret))
	return tokenString, err
}

//ValidateToken ...
func (t *tokenHandler) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	keyFunc := func(token *jwt.Token) (i interface{}, e error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(t.jwtSecret), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}
	return &Claims{
		UserId: claims.UserId,
	}, nil
}

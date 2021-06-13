package middlewares

import (
	"context"
	"fmt"
	"github.com/victor-nach/time-tracker/lib/tokenhandler"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type ctxKey struct {
	Name string
}

var AuthContextKey = &ctxKey{Name: "AuthKey "}

type AuthMiddleware struct {
	tokenHandler tokenhandler.TokenHandler
	logger       *zap.Logger
}

func NewAuthMiddleware(tokenHandler tokenhandler.TokenHandler, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		tokenHandler: tokenHandler,
		logger:       logger,
	}
}

func (A AuthMiddleware) HandleAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			next.ServeHTTP(w, r)
			return
		}

		fmt.Println("authorization", authorization)

		jwtToken := ""
		sp := strings.Split(authorization, " ")
		if len(sp) > 1 {
			jwtToken = sp[1]
		}

		if jwtToken == "" {
			next.ServeHTTP(w, r)
			return
		}

		fmt.Println("token", jwtToken)

		claims, err := A.tokenHandler.ValidateToken(jwtToken)
		if err != nil {
			A.logger.Error("failed to validate token", zap.Error(err))
			next.ServeHTTP(w, r)
			return
		}

		fmt.Println("claims", claims)
		ctx := context.WithValue(r.Context(), AuthContextKey, tokenhandler.Claims{
			UserId: claims.UserId,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

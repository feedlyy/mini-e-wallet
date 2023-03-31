package middleware

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"mini-e-wallet/domain"
	"net/http"
	"strings"
)

type Middleware struct {
	tokenRepo domain.TokenRepository
}

func NewMiddleware(t domain.TokenRepository) Middleware {
	handler := &Middleware{
		tokenRepo: t,
	}

	return *handler
}

// AuthMiddleware function to authenticate requests using the provided token
func (m *Middleware) AuthMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.Background()
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logrus.Error("Middleware | Empty auth header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract the token from the header
		splitHeader := strings.Split(authHeader, " ")
		if len(splitHeader) != 2 || strings.ToLower(splitHeader[0]) != "token" {
			logrus.Error("Middleware | Empty token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		token := splitHeader[1]

		checkToken, err := m.tokenRepo.GetByToken(ctx, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if checkToken == (domain.Tokens{}) {
			logrus.Error("Middleware | Empty data")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add the token to the request context
		ctx = context.WithValue(ctx, "token", token)

		// If authentication succeeded, call the next handler with the modified context
		next(w, r.WithContext(ctx), ps)
	}
}
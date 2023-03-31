package middleware

import (
	"context"
	"github.com/julienschmidt/httprouter"
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
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract the token from the header
		splitHeader := strings.Split(authHeader, " ")
		if len(splitHeader) != 2 || strings.ToLower(splitHeader[0]) != "token" {
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
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// If authentication succeeded, call the next handler
		next(w, r, ps)
	}
}

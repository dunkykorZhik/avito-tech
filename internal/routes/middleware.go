package routes

import (
	"context"
	"errors"
	"net/http"
	"strings"

	errs "github.com/dunkykorZhik/avito-tech/internal/errors"
	"github.com/dunkykorZhik/avito-tech/internal/service"
	"go.uber.org/zap"
)

func LoggingMiddleware(logger *zap.SugaredLogger) func(handlefuncWithError) http.HandlerFunc {
	return func(h handlefuncWithError) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			err := h(w, r)
			if err != nil {
				var statusErr *errs.StatusError
				if errors.As(err, &statusErr) {
					errorJSON(w, err.Error(), statusErr.Status)

				} else {
					logger.Errorf("error : %v, the path: %v%v", err.Error(), r.Method, r.URL.RequestURI())
					errorJSON(w, "internal server error", http.StatusInternalServerError)
				}
				return
			}
		}

	}

}

func AuthMiddleware(service service.User) func(handlefuncWithError) handlefuncWithError {
	return func(hwe handlefuncWithError) handlefuncWithError {
		return func(w http.ResponseWriter, r *http.Request) error {
			token, ok := getBearerToken(r)
			if !ok {
				return errs.ErrUnAuth
			}
			userID, username, err := service.ParseToken(token)
			if err != nil {
				return err
			}
			ctx := context.WithValue(r.Context(), userCtx, userInfo{
				userID,
				username,
			})
			return hwe(w, r.WithContext(ctx))

		}
	}
}

func getBearerToken(r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", false
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {

		return "", false
	}

	token := parts[1]
	return token, true
}

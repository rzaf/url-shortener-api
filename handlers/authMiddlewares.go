package handlers

import (
	"context"
	"errors"
	"github.com/rzaf/url-shortener-api/helpers"
	"github.com/rzaf/url-shortener-api/models"
	"net/http"
)

func AuthApiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")
		if apiKey == "" {
			helpers.WriteJsonError(w, errors.New("header X-API-KEY required"), 401)
			return
		}
		user, _ := models.GetUserByApiKey(apiKey)
		if user == nil {
			helpers.WriteJsonError(w, errors.New("invalid api key"), 401)
			return
		}
		ctx := context.WithValue(r.Context(), models.User{}, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

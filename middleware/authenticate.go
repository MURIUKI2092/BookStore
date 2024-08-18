package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"context"
	"net/http"
	"strings"
	"BookStore/models"
	"BookStore/helpers"
)
var jwtKey = []byte("your_secret_key")
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")[1]

		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
				return
			}
			helpers.RespondWithError(w, http.StatusBadRequest, "Bad request")
			return
		}

		if !token.Valid {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Attach the UUID to the context
		ctx := context.WithValue(r.Context(), "uuid", claims.UUID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})


}


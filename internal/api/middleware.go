package api

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"strings"
)

// JWTMiddleware is a middleware that checks for a valid JWT token in the request.
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		// Extract the token from the Authorization header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Read the JWT secret key from the mounted secret file
		jwtSecretKey, err := os.ReadFile("/run/secrets/jwt_secret")
		if err != nil {
			http.Error(w, "Failed to read JWT secret key", http.StatusInternalServerError)
			return
		}

		// Parse and verify the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecretKey, nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Token is valid, you can now use the claims
			r.Header.Set("EthereumAddress", claims["ethereumAddress"].(string))
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
		}
	})
}

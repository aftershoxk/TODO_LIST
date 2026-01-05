package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"crypto/sha256"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type signinRequest struct {
	Password string `json:"password"`
}

type signinResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

type Claims struct {
	Hash string `json:"hash"`
	jwt.RegisteredClaims
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "method not allowed")
		return
	}

	passEnv := os.Getenv("TODO_PASSWORD")
	if passEnv == "" {
		writeError(w, 400, "password auth disabled")
		return
	}

	var req signinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid json")
		return
	}

	if req.Password != passEnv {
		writeJSON(w, 401, signinResponse{
			Error: "incorrect password",
		})
		return
	}

	hash := hashPassword(passEnv)
	claims := Claims{
		Hash: hash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(passEnv))
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}

	writeJSON(w, 200, signinResponse{
		Token: signed,
	})
}

func hashPassword(p string) string {
	sum := sha256.Sum256([]byte(p))
	return fmt.Sprintf("%x", sum[:])
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		passEnv := os.Getenv("TODO_PASSWORD")
		if passEnv == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(
			tokenStr,
			claims,
			func(t *jwt.Token) (interface{}, error) {
				return []byte(passEnv), nil
			},
		)

		if err != nil || !token.Valid {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		if claims.Hash != hashPassword(passEnv) {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

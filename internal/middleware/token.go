package middleware

import (
	"errors"
	"fmt"
	"gophermart/internal/logger"
	"gophermart/internal/model"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

const cookieName = "authToken"

type makeAuthResponseWriter struct {
	http.ResponseWriter
	signatureKey string
	authToken    string
}

func (r *makeAuthResponseWriter) Write(b []byte) (int, error) {
	return r.ResponseWriter.Write(b)
}

func (r *makeAuthResponseWriter) WriteHeader(statusCode int) {
	if r.authToken == "" {
		var err error
		userID := r.Header().Get(string(model.UserIDKey))
		r.authToken, err = MakeAuthToken(r.signatureKey, userID)
		if err != nil {
			logger.Log.Error("WithAuth middleware. WriteHeader MakeAuthToken error", zap.String("error: ", err.Error()), zap.String("UserID: ", userID))
		}
	}

	cookie := http.Cookie{Name: cookieName, Value: r.authToken}
	http.SetCookie(r.ResponseWriter, &cookie)

	logger.Log.Info("MakeAuth middleware. WriteHeader with cookie", zap.String("cookie name: ", cookie.Name), zap.String("cookie value: ", cookie.Value))

	r.ResponseWriter.WriteHeader(statusCode)
}

type checkAuthResponseWriter struct {
	http.ResponseWriter
	authToken string
}

func (r *checkAuthResponseWriter) Write(b []byte) (int, error) {
	return r.ResponseWriter.Write(b)
}

func (r *checkAuthResponseWriter) WriteHeader(statusCode int) {
	cookie := http.Cookie{Name: cookieName, Value: r.authToken}
	http.SetCookie(r.ResponseWriter, &cookie)

	logger.Log.Info("CheckAuth middleware. WriteHeader with cookie", zap.String("cookie name: ", cookie.Name), zap.String("cookie value: ", cookie.Value))

	http.SetCookie(r.ResponseWriter, &cookie)

	r.ResponseWriter.WriteHeader(statusCode)
}

// Claims — структура утверждений, которая включает стандартные утверждения
// и одно пользовательское — UserID
type claims struct {
	jwt.RegisteredClaims
	UserID string
}

const tokenExp = time.Hour * 3

func getUserID(key, tokenString string) (string, error) {
	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})
	if err != nil {
		return "", fmt.Errorf("token-GetUserId-ParseWithClaims-err: %w", err)
	}

	if !token.Valid {
		return "", errors.New("token-GetUserId-TokenIsNotValid")
	}

	return claims.UserID, nil
}

func MakeAuthToken(key, userID string) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		// собственное утверждение
		UserID: userID,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("token-MakeAuthToken-signedToken-err: %w", err)
	}

	// возвращаем строку токена
	return tokenString, nil
}

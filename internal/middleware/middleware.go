package middleware

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gophermart/internal/logger"
	"gophermart/internal/model"
	"net/http"
	"strings"
	"time"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// WithLogging — middleware-логер для входящих HTTP-запросов.
func WithLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.Log.Info("got incoming HTTP request",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
			zap.Duration("duration", duration),
			zap.Int("status", responseData.status),
			zap.Int("size", responseData.size),
		)
	})
}

// WithGzip - middleware поддерживающий gzip компрессию и декомпрессию
func WithGzip(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w
		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		clientSupportsGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		logger.Log.Info("withGzip middleware", zap.String("Accept-Encoding", r.Header.Get("Accept-Encoding")))
		if clientSupportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		clientSentGzip := strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
		if clientSentGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				http.Error(w, "Add gzip compress error", http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, r)

	})
}

// WithCheckAuth - middleware который чекает авторизацию
func WithCheckAuth(key string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Log.Info("WithCheckAuth middleware")

			tokenWithUser, err := r.Cookie(cookieName)
			if tokenWithUser != nil {
				logger.Log.Info("WithAuth middleware. tokenWithUser != nil", zap.String(" tokenWithUser.Value", tokenWithUser.Value))
			} else {
				logger.Log.Info("WithAuth middleware. tokenWithUser == nil")
				http.Error(w, "tokenWithUser == nil", http.StatusUnauthorized)
				return
			}

			if err != nil && !errors.Is(err, http.ErrNoCookie) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if err != nil {
				logger.Log.Info("WithAuth middleware. err from r.Cookie() != nil", zap.String("error: ", err.Error()))
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			var userID string
			if tokenWithUser.Value != "" {
				logger.Log.Info("WithAuth middleware. tokenWithUser.Value != ''", zap.String(" tokenWithUser.Value: ", tokenWithUser.Value))
				userID, err = getUserID(key, tokenWithUser.Value)
				logger.Log.Info("WithAuth middleware. token.GetUserID", zap.String("userID: ", userID))
				if err != nil {
					http.Error(w, "Get userID from token error", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "tokenWithUser.Value == ''", http.StatusUnauthorized)
				return
			}

			if userID == "" {
				http.Error(w, "Absent user_id at authToken", http.StatusUnauthorized)
				return
			}

			logger.Log.Info("Известный юзер", zap.String("user_id", userID))

			aw := checkAuthResponseWriter{
				ResponseWriter: w,
				authToken:      tokenWithUser.Value,
			}

			userForContext := model.ContextKey(userID)
			ctx := context.WithValue(r.Context(), model.UserIDKey, userForContext)
			r = r.WithContext(ctx)

			h.ServeHTTP(&aw, r)
		})
	}
}

// WithMakeAuth - middleware который навешивает куку для авторизации
func WithMakeAuth(key string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Log.Info("WithMakeAuth middleware")

			aw := makeAuthResponseWriter{
				ResponseWriter: w,
				signatureKey:   key,
			}

			h.ServeHTTP(&aw, r)
		})
	}
}

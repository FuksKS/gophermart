package handlers

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"gophermart/internal/logger"
	"gophermart/internal/model"
	"io"
	"net/http"
	"strconv"
)

func (h *GmHandler) register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Log.Error("register reading request body error", zap.String("error", err.Error()))
			http.Error(w, "register reading request body error", http.StatusInternalServerError)
			return
		}

		var req model.LogoPass
		err = json.Unmarshal(body, &req)
		if err != nil {
			logger.Log.Error("register unmarshal body error", zap.String("error", err.Error()))
			http.Error(w, "register unmarshal body error", http.StatusBadRequest)
			return
		}

		userID, err := h.gmService.AddAuthInfo(ctx, req.Login, req.Password, h.passKey)
		if err != nil {
			if errors.Is(err, model.ErrLoginAlreadyExist) {
				logger.Log.Error("register AddAuthInfo error", zap.String("login", req.Login), zap.String("error", "login already exist"))
				http.Error(w, "login already exist", http.StatusConflict)
				return
			}
			logger.Log.Error("register AddAuthInfo error", zap.String("login", req.Login), zap.String("error", err.Error()))
			http.Error(w, "register AddAuthInfo error", http.StatusInternalServerError)
			return
		}

		if userID == 0 {
			logger.Log.Error("register AddAuthInfo error", zap.String("login", req.Login), zap.String("error", "login is already in use"))
			http.Error(w, "user_id = 0", http.StatusConflict)
			return
		}

		// закидываем юзера в хедеры чтоб потом навесить куку
		w.Header().Set(string(model.UserIDKey), strconv.FormatInt(userID, 10))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func (h *GmHandler) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Log.Error("login reading request body error", zap.String("error", err.Error()))
			http.Error(w, "login reading request body error", http.StatusInternalServerError)
			return
		}

		var req model.LogoPass
		err = json.Unmarshal(body, &req)
		if err != nil {
			logger.Log.Error("login unmarshal body error", zap.String("error", err.Error()))
			http.Error(w, "login unmarshal body error", http.StatusBadRequest)
			return
		}

		userID, err := h.gmService.GetAuthInfo(ctx, req.Login, req.Password, h.passKey)
		if err != nil {
			if errors.Is(err, model.ErrWrongLogin) {
				logger.Log.Error("login GetAuthInfo error", zap.String("login", req.Login), zap.String("error", err.Error()))
				http.Error(w, "login does not exist", http.StatusUnauthorized)
				return
			} else if errors.Is(err, model.ErrWrongPas) {
				logger.Log.Error("login GetAuthInfo error", zap.String("login", req.Login), zap.String("error", err.Error()))
				http.Error(w, "wrong password", http.StatusUnauthorized)
				return
			} else {
				logger.Log.Error("login GetAuthInfo error", zap.String("login", req.Login), zap.String("error", err.Error()))
				http.Error(w, "login GetAuthInfo error", http.StatusInternalServerError)
				return
			}
		}

		// закидываем юзера в хедеры чтоб потом навесить куку
		w.Header().Set(string(model.UserIDKey), strconv.FormatInt(userID, 10))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

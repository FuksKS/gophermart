package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gophermart/internal/logger"
	"gophermart/internal/luhnalgorithm"
	"gophermart/internal/model"
	"io"
	"net/http"
	"strconv"
)

func (h *GmHandler) addOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Log.Error("addOrder reading request body error", zap.String("error", err.Error()))
			http.Error(w, "addOrder reading request body error", http.StatusInternalServerError)
			return
		}

		orderID := string(body)

		isCorrect, err := luhnalgorithm.LuhnCheck(orderID)
		if err != nil && errors.Is(err, model.NotANumberError) {
			logger.Log.Error("addOrder LuhnCheck error", zap.String("error", err.Error()))
			fmt.Println("addOrder LuhnCheck error. order id is not a number")
			http.Error(w, "order id is not a number", http.StatusUnprocessableEntity)
			return
		}

		if !isCorrect {
			logger.Log.Error("addOrder LuhnCheck error", zap.String("error", "incorrect order id"))
			fmt.Println("addOrder LuhnCheck error. incorrect order id")
			http.Error(w, "incorrect order id", http.StatusUnprocessableEntity)
			return
		}

		userID, ok := ctx.Value(model.UserIDKey).(model.ContextKey)
		if !ok {
			logger.Log.Error("addOrder get user_id from context error")
			http.Error(w, "addOrder get user_id from context error", http.StatusInternalServerError)
			return
		}

		userInt64, err := strconv.ParseInt(string(userID), 10, 64)
		if err != nil {
			logger.Log.Error("addOrder parse user_id to int64", zap.String("error", err.Error()))
			http.Error(w, "addOrder parse user_id to int64", http.StatusInternalServerError)
			return
		}

		err = h.gmService.AddOrder(ctx, orderID, userInt64)
		if err != nil {
			if errors.Is(err, model.AlreadyUploadedByThisUser) {
				logger.Log.Warn("AddOrder error", zap.String("error", err.Error()))
				w.WriteHeader(http.StatusOK)
				return
			} else if errors.Is(err, model.AlreadyUploadedByAnotherUser) {
				logger.Log.Error("AddOrder error", zap.String("error", err.Error()))
				http.Error(w, "order id has already been uploaded by another user", http.StatusConflict)
				return
			} else {
				logger.Log.Error("AddOrder error", zap.String("error", err.Error()))
				http.Error(w, "AddOrder error", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusAccepted)
	}
}

func (h *GmHandler) getOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := ctx.Value(model.UserIDKey).(model.ContextKey)
		if !ok {
			logger.Log.Error("getOrders get user_id from context error")
			http.Error(w, "getOrders get user_id from context error", http.StatusInternalServerError)
			return
		}

		userInt64, err := strconv.ParseInt(string(userID), 10, 64)
		if err != nil {
			logger.Log.Error("getOrders parse user_id to int64", zap.String("error", err.Error()))
			http.Error(w, "getOrders parse user_id to int64", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		orders, err := h.gmService.GetOrders(ctx, userInt64)
		if err != nil {
			logger.Log.Error("getOrders error", zap.String("error", err.Error()))
			http.Error(w, "getOrders error", http.StatusInternalServerError)
			return
		}

		if len(orders) == 0 {
			logger.Log.Info("getOrders user has no orders", zap.String("user_id", string(userID)))
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.WriteHeader(http.StatusOK)
		resp, err := json.Marshal(orders)
		if err != nil {
			http.Error(w, "getOrders marshal response error", http.StatusInternalServerError)
			return
		}

		w.Write(resp)
	}
}

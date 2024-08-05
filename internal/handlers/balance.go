package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"gophermart/internal/logger"
	"gophermart/internal/luhnalgorithm"
	"gophermart/internal/model"
	"io"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

func (h *GmHandler) getBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := ctx.Value(model.UserIDKey).(model.ContextKey)
		if !ok {
			logger.Log.Error("getBalance get user_id from context error")
			http.Error(w, "getBalance get user_id from context error", http.StatusInternalServerError)
			return
		}

		userInt64, err := strconv.ParseInt(string(userID), 10, 64)
		if err != nil {
			logger.Log.Error("getBalance parse user_id to int64", zap.String("error", err.Error()))
			http.Error(w, "getBalance parse user_id to int64", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		balance, err := h.gmService.GetBalance(ctx, userInt64)
		if err != nil {
			logger.Log.Error("getBalance error", zap.String("error", err.Error()))
			http.Error(w, "getBalance error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		resp, err := json.Marshal(balance)
		if err != nil {
			http.Error(w, "getBalance marshal response error", http.StatusInternalServerError)
			return
		}

		w.Write(resp)
	}
}

func (h *GmHandler) withdraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Log.Error("withdraw reading request body error", zap.String("error", err.Error()))
			http.Error(w, "withdraw reading request body error", http.StatusInternalServerError)
			return
		}

		var req model.Withdraw
		err = json.Unmarshal(body, &req)
		if err != nil {
			logger.Log.Error("withdraw unmarshal body error", zap.String("error", err.Error()))
			http.Error(w, "withdraw unmarshal body error", http.StatusBadRequest)
			return
		}

		isCorrect, err := luhnalgorithm.LuhnCheck(req.OrderID)
		if err != nil && errors.Is(err, model.ErrNotANumber) {
			logger.Log.Error("withdraw LuhnCheck error", zap.String("error", err.Error()))
			fmt.Println("withdraw LuhnCheck error. order id is not a number")
			http.Error(w, "order id is not a number", http.StatusUnprocessableEntity)
			return
		}

		if !isCorrect {
			logger.Log.Error("withdraw LuhnCheck error", zap.String("error", "incorrect order id"))
			fmt.Println("withdraw LuhnCheck error. incorrect order id")
			http.Error(w, "incorrect order id", http.StatusUnprocessableEntity)
			return
		}

		userID, ok := ctx.Value(model.UserIDKey).(model.ContextKey)
		if !ok {
			logger.Log.Error("withdraw get user_id from context error")
			http.Error(w, "withdraw get user_id from context error", http.StatusInternalServerError)
			return
		}

		userInt64, err := strconv.ParseInt(string(userID), 10, 64)
		if err != nil {
			logger.Log.Error("withdraw parse user_id to int64", zap.String("error", err.Error()))
			http.Error(w, "withdraw parse user_id to int64", http.StatusInternalServerError)
			return
		}

		req.UserID = userInt64

		err = h.gmService.Withdraw(ctx, req)
		if err != nil {
			if errors.Is(err, model.ErrOrderAlreadyUploaded) {
				logger.Log.Error("Withdraw error", zap.String("error", err.Error()))
				http.Error(w, "order id has already been uploaded", http.StatusConflict)
				return
			} else if errors.Is(err, model.ErrNotEnoughMoney) {
				logger.Log.Error("Withdraw error", zap.String("error", err.Error()))
				http.Error(w, "not enough money", http.StatusPaymentRequired)
				return
			} else {
				logger.Log.Error("Withdraw error", zap.String("error", err.Error()))
				http.Error(w, "Withdraw error", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	}
}

func (h *GmHandler) getWithdrawals() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := ctx.Value(model.UserIDKey).(model.ContextKey)
		if !ok {
			logger.Log.Error("getWithdrawals get user_id from context error")
			http.Error(w, "getWithdrawals get user_id from context error", http.StatusInternalServerError)
			return
		}

		userInt64, err := strconv.ParseInt(string(userID), 10, 64)
		if err != nil {
			logger.Log.Error("getWithdrawals parse user_id to int64", zap.String("error", err.Error()))
			http.Error(w, "getWithdrawals parse user_id to int64", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		withdrawals, err := h.gmService.GetWithdrawals(ctx, userInt64)
		if err != nil {
			logger.Log.Error("getWithdrawals error", zap.String("error", err.Error()))
			http.Error(w, "getWithdrawals error", http.StatusInternalServerError)
			return
		}

		if len(withdrawals) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.WriteHeader(http.StatusOK)
		resp, err := json.Marshal(withdrawals)
		if err != nil {
			http.Error(w, "getWithdrawals marshal response error", http.StatusInternalServerError)
			return
		}

		w.Write(resp)
	}
}

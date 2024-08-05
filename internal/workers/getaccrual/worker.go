package getaccrual

import (
	"context"
	"errors"
	"fmt"
	"gophermart/internal/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type accrualWorker struct {
	storager       storager
	accrualService accrualService
}

func New(storager storager, accrualService accrualService) *accrualWorker {
	accrualWorker := accrualWorker{
		storager:       storager,
		accrualService: accrualService,
	}
	return &accrualWorker
}

func (w *accrualWorker) Process(ctx context.Context) error {
	orderID, err := w.storager.GetOrderForAccrual(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Log.Info("getAccrualWorker-noRowsForProcessing")
			time.Sleep(200 * time.Millisecond)
			return nil
		}
		logger.Log.Error("getAccrualWorker-storager-GetOrderForAccrual-err", zap.Error(err))
		return err
	}

	code, accrual, err := w.accrualService.GetAccrual(orderID)
	if err != nil {
		logger.Log.Error("getAccrualWorker-accrualService-GetAccrual-err", zap.Error(err), zap.String("order_id", orderID))
		return err
	}

	logger.Log.Info("getAccrualWorker-accrualService-GetAccrual", zap.String("order_id", orderID), zap.String("status code", strconv.Itoa(code)), zap.String("status", accrual.Status))

	switch code {
	case http.StatusNoContent:
		logger.Log.Warn("getAccrualWorker-accrualService-GetAccrual-StatusNoContent", zap.String("order_id", orderID))
		return nil
	case http.StatusTooManyRequests:
		logger.Log.Warn("getAccrualWorker-accrualService-GetAccrual-StatusTooManyRequests", zap.String("order_id", orderID))
		time.Sleep(2 * time.Second)
		return nil
	case http.StatusOK:
	default:
		err = fmt.Errorf("getAccrualWorker accrualService incorrect responce code: %d", code)
		logger.Log.Error(err.Error(), zap.String("order_id", orderID))
		return err
	}

	err = w.storager.SetAccrual(ctx, accrual)
	if err != nil {
		logger.Log.Error("getAccrualWorker-storager-SetAccrual-err", zap.Error(err))
		return err
	}

	return nil
}

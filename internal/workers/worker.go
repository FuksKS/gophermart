package workers

import (
	"context"
	"go.uber.org/zap"
	"gophermart/internal/logger"
	"strconv"
	"time"
)

type Worker interface {
	Process(ctx context.Context) error
}

func Start(ctx context.Context, w Worker, schedule time.Duration, workerNumber int) {
	go run(ctx, w, schedule, workerNumber)
}

func run(ctx context.Context, w Worker, period time.Duration, workerNumber int) {
	for {
		time.Sleep(period)
		select {
		case <-ctx.Done():
			return
		default:
			if err := w.Process(ctx); err != nil {
				logger.Log.Error("Worker-error", zap.Error(err), zap.String("worker number", strconv.Itoa(workerNumber)))
			}
		}
	}
}

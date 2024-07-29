package main

import (
	"context"
	"go.uber.org/zap"
	"gophermart/internal/accrual"
	"gophermart/internal/config"
	"gophermart/internal/logger"
	"gophermart/internal/pg"
	"gophermart/internal/service"
	"gophermart/internal/workers"
	"gophermart/internal/workers/getaccrual"
	"os"
	"os/signal"
	"syscall"
)

const (
	workerCount    = 2
	workerSchedule = 0
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cfg := config.Init()

	if err := logger.Init(logger.LoggerLevelINFO); err != nil {
		logger.Log.Fatal(err.Error(), zap.String("init client", "logger Initialize"))
	}
	logger.Log.Info("Step 1", zap.String("init client", "config Initialized"))

	db, err := pg.NewConnect(ctx, cfg.DBURI)
	if err != nil {
		logger.Log.Fatal(err.Error(), zap.String("init client", "db Initialize"))
	}
	logger.Log.Info("Step 2", zap.String("init client", "db Initialized"), zap.String("cfg.DBURI", cfg.DBURI))

	serv := service.New(db)
	logger.Log.Info("Step 3", zap.String("init client", "service Initialized"))

	accrualClient := accrual.NewClient(cfg.AccrualAddr)
	accrualWorker := getaccrual.New(serv, accrualClient)

	for i := 0; i < workerCount; i++ {
		workers.Start(ctx, accrualWorker, workerSchedule, i)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

	cancel()

	logger.Log.Info("Terminated client. Goodbye")
}

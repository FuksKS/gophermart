package main

import (
	"context"
	"go.uber.org/zap"
	"gophermart/internal/accrual"
	"gophermart/internal/config"
	"gophermart/internal/handlers"
	"gophermart/internal/logger"
	"gophermart/internal/pg"
	"gophermart/internal/service"
	"gophermart/internal/workers"
	"gophermart/internal/workers/getaccrual"
	"net/http"
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
		logger.Log.Fatal(err.Error(), zap.String("init", "logger Initialize"))
	}
	logger.Log.Info("Step 1", zap.String("init", "config Initialized"))

	db, err := pg.NewConnect(ctx, cfg.DBURI)
	if err != nil {
		logger.Log.Fatal(err.Error(), zap.String("init", "db Initialize"))
	}
	logger.Log.Info("Step 2", zap.String("init", "db Initialized"), zap.String("cfg.DBURI", cfg.DBURI))

	serv := service.New(db)
	logger.Log.Info("Step 3", zap.String("init", "service Initialized"))

	handler, err := handlers.New(serv, cfg.SignatureKey, cfg.PassKey)
	if err != nil {
		logger.Log.Fatal(err.Error(), zap.String("init", "set handler"))
	}
	logger.Log.Info("Step 4", zap.String("init", "handler Initialized"))

	accrualClient := accrual.NewClient(cfg.AccrualAddr)
	accrualWorker := getaccrual.New(serv, accrualClient)

	for i := 0; i < workerCount; i++ {
		workers.Start(ctx, accrualWorker, workerSchedule, i)
	}

	logger.Log.Info("Running server", zap.String("address", cfg.HTTPAddr))
	go func() {
		if err := http.ListenAndServe(cfg.HTTPAddr, handler.InitRouter()); err != nil {
			logger.Log.Fatal(err.Error(), zap.String("event", "start server"))
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

	logger.Log.Info("Stop server", zap.String("address", cfg.HTTPAddr))

	cancel()

	logger.Log.Info("Terminated. Goodbye")
}
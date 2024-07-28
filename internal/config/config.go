package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"log"
)

const (
	defaultAddr         = "localhost:8080"
	defaultDBURI        = "postgres://fuks:pass@localhost:5432/accrual"
	defaultAccrualAddr  = "http://localhost:8000"
	defaultSignatureKey = "super_secret"
	lenPassKey          = 32
	defaultPassKey      = "myverystrongpasswordo32bitlength"
)

type Config struct {
	HTTPAddr     string `env:"RUN_ADDRESS"`
	DBURI        string `env:"DATABASE_URI"`
	AccrualAddr  string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	SignatureKey string `env:"SIGNATURE_KEY"` // ключ для подписи кук при авторизации
	PassKey      []byte `env:"PASS_KEY"`      // симметричный ключ для паролей юзеров
}

func Init() *Config {
	var cfg Config
	if err := envConfig(&cfg); err != nil {
		log.Fatal(err)
	}

	flagAddr, flagDBURI, flagAccrualAddr, flagSignKey, flagPassKey := flagConfig()
	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = flagAddr
	}
	if cfg.DBURI == "" {
		cfg.DBURI = flagDBURI
	}
	if cfg.AccrualAddr == "" {
		cfg.AccrualAddr = flagAccrualAddr
	}
	if cfg.SignatureKey == "" {
		cfg.SignatureKey = flagSignKey
	}

	if cfg.PassKey == nil {
		cfg.PassKey = []byte(flagPassKey)
	} else {
		if len(cfg.PassKey) != lenPassKey {
			cfg.PassKey = []byte(flagPassKey)
		}
	}

	return &cfg
}

func flagConfig() (flagAddr, flagDBDSN, flagAccrualAddr, flagSignKey, flagPassKey string) {
	flag.StringVar(&flagAddr, "a", defaultAddr, "адрес запуска HTTP-сервера")
	flag.StringVar(&flagDBDSN, "d", defaultDBURI, "строка с адресом подключения к БД")
	flag.StringVar(&flagAccrualAddr, "r", defaultAccrualAddr, "адрес системы расчёта начислений")
	flag.StringVar(&flagSignKey, "sk", defaultSignatureKey, "ключ для подписи кук при авторизации")
	flag.StringVar(&flagPassKey, "pk", defaultPassKey, "симметричный ключ для паролей юзеров длинной 32 байта")

	flag.Parse()
	return
}

func envConfig(cfg *Config) error {
	if err := env.Parse(cfg); err != nil {
		return fmt.Errorf("InitConfig-envConfig-err: %w", err)
	}
	return nil
}

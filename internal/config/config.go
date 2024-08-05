package config

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	defaultAddr          = "localhost:8080"
	defaultDBURI         = "postgres://fuks:pass@localhost:5432/gophermart"
	defaultAccrualAddr   = "http://localhost:8000"
	defaultSignatureKey  = "super_secret"
	lenPassKey           = 32
	defaultPassKey       = "myverystrongpasswordo32bitlength"
	defaultDBTimeout     = 3 * time.Second
	defaultClientTimeout = 3 * time.Second
)

type Config2 struct {
	HTTPAddr      string        `env:"RUN_ADDRESS"`
	DBURI         string        `env:"DATABASE_URI"`
	AccrualAddr   string        `env:"ACCRUAL_SYSTEM_ADDRESS"`
	SignatureKey  string        `env:"SIGNATURE_KEY"` // ключ для подписи кук при авторизации
	PassKey       []byte        `env:"PASS_KEY"`      // симметричный ключ для паролей юзеров
	DBTimeout     time.Duration `env:"DB_TIMEOUT"`    // таймаут на запросы в базу
	ClientTimeout time.Duration `env:"CLIEN_TIMEOUT"` // клиентский таймаут
}

type Config struct {
	ServerConfig ServerConfig
	DBConfig     DBConfig
	ClientConfig ClientConfig
}

type ServerConfig struct {
	HTTPAddr     string `env:"RUN_ADDRESS"`
	SignatureKey string `env:"SIGNATURE_KEY"` // ключ для подписи кук при авторизации
	PassKey      []byte `env:"PASS_KEY"`      // симметричный ключ для паролей юзеров
}

type DBConfig struct {
	DBURI     string        `env:"DATABASE_URI"`
	DBTimeout time.Duration `env:"DB_TIMEOUT"` // таймаут на запросы в базу
}

type ClientConfig struct {
	AccrualAddr   string        `env:"ACCRUAL_SYSTEM_ADDRESS"`
	ClientTimeout time.Duration `env:"CLIEN_TIMEOUT"` // клиентский таймаут
}

func Init() *Config {
	var cfg Config
	if err := envConfig(&cfg); err != nil {
		log.Fatal(err)
	}

	flagAddr, flagDBURI, flagAccrualAddr, flagSignKey, flagPassKey := flagConfig()
	if cfg.ServerConfig.HTTPAddr == "" {
		cfg.ServerConfig.HTTPAddr = flagAddr
	}
	if cfg.DBConfig.DBURI == "" {
		cfg.DBConfig.DBURI = flagDBURI
	}
	if cfg.ClientConfig.AccrualAddr == "" {
		cfg.ClientConfig.AccrualAddr = flagAccrualAddr
	}
	if cfg.ServerConfig.SignatureKey == "" {
		cfg.ServerConfig.SignatureKey = flagSignKey
	}

	if cfg.ServerConfig.PassKey == nil {
		cfg.ServerConfig.PassKey = []byte(flagPassKey)
	} else {
		if len(cfg.ServerConfig.PassKey) != lenPassKey {
			cfg.ServerConfig.PassKey = []byte(flagPassKey)
		}
	}

	if cfg.DBConfig.DBTimeout == time.Duration(0) {
		cfg.DBConfig.DBTimeout = defaultDBTimeout
	}

	if cfg.ClientConfig.ClientTimeout == time.Duration(0) {
		cfg.ClientConfig.ClientTimeout = defaultClientTimeout
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

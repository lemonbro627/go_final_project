package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultPassword = "12345678"
	defaultPort     = "7540"
)

type Config struct {
	AppPassword         string
	EncryptionSecretKey string // секретный ключ шифрования
	ApiPort             string
	DbPath              string
}

// NewConfig конструктор конфигурации приложения
func NewConfig() (*Config, error) {
	appPass := os.Getenv("TODO_PASSWORD") //пароль для входа в панель приложения
	encKey := "superpupersecret"          //пароль для шифрования JWT токена
	apiPort := os.Getenv("TODO_PORT")     //порт для прослушивания веб-сервером
	dbPath := os.Getenv("TODO_DBFILE")    //путь у файлу базы данных

	if appPass == "" {
		appPass = defaultPassword
	}

	if appPass == "" || encKey == "" {
		return nil, fmt.Errorf("invalid config")
	}

	if apiPort == "" {
		apiPort = defaultPort
	}

	if dbPath == "" {
		appPath, err := os.Executable()
		if err != nil {
			return nil, err
		}
		dbPath = filepath.Join(filepath.Dir(appPath), "scheduler.db")
	}

	return &Config{AppPassword: appPass, EncryptionSecretKey: encKey, ApiPort: apiPort, DbPath: dbPath}, nil
}

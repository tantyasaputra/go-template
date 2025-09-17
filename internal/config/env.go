package config

import (
	"fmt"
	"go-template/internal/log"
	"os"
	"path"
	"strings"

	"github.com/joho/godotenv"
)

// Config is struct for all env
type Config struct {
	DEVMODE string
	PORT    string
	DB      string
	DBT     string // Test DB for unit testing
	// Removed Keycloak-related configuration
	DatalakeProjectId              string
	AllowedServiceAccountClientIds []string
}

func SetDevelopmentEnv() error {
	log.Init()

	LoadEnv()

	if strings.TrimSpace(os.Getenv("PORT")) != "" {
		return nil
	}
	var err error
	dbname := "template" // replace with BU name
	err = os.Setenv("DEVMODE", "true")
	if err != nil {
		return err
	}
	err = os.Setenv("PORT", "5000")
	if err != nil {
		return err
	}
	err = os.Setenv("CONNECTION_STRING", "postgres://postgres:admin@localhost/"+dbname+"?sslmode=disable")
	if err != nil {
		return err
	}
	// only for testing
	err = os.Setenv("CONNECTION_STRING_TEST", "postgres://postgres:admin@localhost/"+dbname+"_test?sslmode=disable")
	if err != nil {
		return err
	}
	return nil
}

// GetEnv return all env on config struct
func GetEnv() *Config {
	return &Config{
		DEVMODE:                        os.Getenv("DEVMODE"),
		PORT:                           os.Getenv("PORT"),
		DB:                             os.Getenv("CONNECTION_STRING"),
		DBT:                            os.Getenv("CONNECTION_STRING_TEST"),
		DatalakeProjectId:              os.Getenv("DATALAKE_PROJECT_ID"),
		AllowedServiceAccountClientIds: strings.Split(os.Getenv("ALLOWED_SERVICE_ACCOUNT_CLIENT_IDS"), " "),
	}
}

func getEnvPath() (string, error) {
	directory, err := os.Getwd()
	if err != nil {
		return "", err
	}

	filepath := searchup(directory, ".env")
	return filepath, nil
}

func searchup(dir string, filename string) string {
	if dir == "/" || dir == "" || dir == "." {
		return ""
	}

	if _, err := os.Stat(path.Join(dir, filename)); err == nil {
		return path.Join(dir, filename)
	}

	return searchup(path.Dir(dir), filename)
}

func LoadEnv() {
	path, err := getEnvPath()

	if err != nil || strings.TrimSpace(path) == "" {
		log.Warnw("Environment Loading", "event", "config", "service", "IKN_B2B", "message", "No .env file found, env should be injected by other methods")
		return
	}

	log.Infow("Environment Loading", "event", "config", "message", fmt.Sprint("loading env from ", path))

	err = godotenv.Load(path)

	if err != nil {
		log.Warnw("Environment Loading", "event", "config", "message", "failed to load env")
	}
}

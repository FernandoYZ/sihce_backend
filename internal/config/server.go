package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	JWT      JWTConfig      `yaml:"jwt"`
	Security SecurityConfig `yaml:"security"`
	App      AppConfig      `yaml:"app"`
}

type JWTConfig struct {
	AccessSecret            string `yaml:"access_secret"`
	RefreshSecret           string `yaml:"refresh_secret"`
	AccessTokenExpiration   int    `yaml:"access_token_expiration_seconds"`
	RefreshTokenExpiration  int    `yaml:"refresh_token_expiration_seconds"`
}

type SecurityConfig struct {
	SessionSecret  string `yaml:"session_secret"`
	HashSaltRounds int    `yaml:"hash_salt_rounds"`
}

type AppConfig struct {
	Port        int      `yaml:"port"`
	CorsOrigins []string `yaml:"cors_origins"`
	LogLevel    string   `yaml:"log_level"`
	AppEnv      string   `yaml:"app_env"`
}

var (
	cfg     *Config
	cfgOnce sync.Once
)

func Cargar(ruta string) (*Config, error) {
	var err error

	cfgOnce.Do(func() {
		archivo, readErr := os.ReadFile(ruta)
		if readErr != nil {
			err = fmt.Errorf("error al leer config: %w", readErr)
			return
		}

		contenido := os.ExpandEnv(string(archivo))

		var c Config
		if parseErr := yaml.Unmarshal([]byte(contenido), &c); parseErr != nil {
			err = fmt.Errorf("error al parsear config: %w", parseErr)
			return
		}

		cfg = &c
	})

	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func Obtener() *Config {
	return cfg
}

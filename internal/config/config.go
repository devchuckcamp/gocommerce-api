package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/devchuckcamp/goauthx"
	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Driver          string // postgres, mysql, sqlserver
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret          string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	JWTIssuer          string
	JWTAudience        string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	GoogleOAuthEnabled bool
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file (optional)
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 120*time.Second),
		},
		Database: DatabaseConfig{
			Driver:          getEnv("DB_DRIVER", "postgres"),
			DSN:             getEnv("DB_DSN", ""),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Auth: AuthConfig{
			JWTSecret:          getEnv("JWT_SECRET", ""),
			AccessTokenExpiry:  getDurationEnv("JWT_ACCESS_TOKEN_EXPIRY", 15*time.Minute),
			RefreshTokenExpiry: getDurationEnv("JWT_REFRESH_TOKEN_EXPIRY", 7*24*time.Hour),
			JWTIssuer:          getEnv("JWT_ISSUER", "gocommerce-api"),
			JWTAudience:        getEnv("JWT_AUDIENCE", "gocommerce-api-users"),
			GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/google/callback"),
			GoogleOAuthEnabled: getEnv("GOOGLE_CLIENT_ID", "") != "" && getEnv("GOOGLE_CLIENT_SECRET", "") != "",
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Database.DSN == "" {
		return fmt.Errorf("DB_DSN is required")
	}

	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	if len(c.Auth.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}

	validDrivers := map[string]bool{
		"postgres":  true,
		"mysql":     true,
		"sqlserver": true,
	}
	if !validDrivers[c.Database.Driver] {
		return fmt.Errorf("invalid DB_DRIVER: %s (must be postgres, mysql, or sqlserver)", c.Database.Driver)
	}

	return nil
}

// ToGoAuthXConfig converts our config to goauthx.Config
func (c *Config) ToGoAuthXConfig() *goauthx.Config {
	var driver goauthx.DatabaseDriver
	switch c.Database.Driver {
	case "postgres":
		driver = goauthx.Postgres
	case "mysql":
		driver = goauthx.MySQL
	case "sqlserver":
		driver = goauthx.SQLServer
	}

	return &goauthx.Config{
		Database: goauthx.DatabaseConfig{
			Driver:          driver,
			DSN:             c.Database.DSN,
			MaxOpenConns:    c.Database.MaxOpenConns,
			MaxIdleConns:    c.Database.MaxIdleConns,
			ConnMaxLifetime: c.Database.ConnMaxLifetime,
		},
		JWT: goauthx.JWTConfig{
			Secret:            c.Auth.JWTSecret,
			AccessTokenExpiry: c.Auth.AccessTokenExpiry,
			Issuer:            c.Auth.JWTIssuer,
			Audience:          c.Auth.JWTAudience,
		},
		Password: goauthx.PasswordConfig{
			MinLength:  8,
			BcryptCost: 12,
		},
		Token: goauthx.TokenConfig{
			RefreshTokenExpiry: c.Auth.RefreshTokenExpiry,
			RefreshTokenLength: 64,
		},
		OAuth: goauthx.OAuthConfig{
			Google: goauthx.GoogleOAuthConfig{
				ClientID:     c.Auth.GoogleClientID,
				ClientSecret: c.Auth.GoogleClientSecret,
				RedirectURL:  c.Auth.GoogleRedirectURL,
				Enabled:      c.Auth.GoogleOAuthEnabled,
			},
		},
	}
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

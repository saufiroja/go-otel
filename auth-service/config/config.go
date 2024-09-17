package config

import (
	"github.com/joho/godotenv"
	"github.com/saufiroja/go-otel/auth-service/pkg/logging"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"sync"
)

type AppConfig struct {
	App struct {
		Env string
	}
	Http struct {
		Port string
	}
	Postgres struct {
		Name string
		User string
		Pass string
		Host string
		Port string
		SSL  string
	}
	Jwt struct {
		Secret string
	}
	Otel struct {
		OTLPEndpoint string
	}
}

var appConfig *AppConfig
var lock = &sync.Mutex{}

func NewAppConfig(logging logging.Logger) *AppConfig {
	// add config file path in .env
	_ = godotenv.Load("../.env")

	if appConfig == nil {
		lock.Lock()
		defer lock.Unlock()
		if appConfig == nil {
			logging.LogInfo("Creating AppConfig first time")
			appConfig = &AppConfig{}

			appConfig.initApp()
			appConfig.initHttp()
			appConfig.initPostgres()
			appConfig.initJwt()
			appConfig.initOtel()
		} else {
			logging.LogInfo("AppConfig already created")
		}
	} else {
		logging.LogInfo("AppConfig already created")
	}

	return appConfig
}

func (c *AppConfig) initApp() {
	c.App.Env = os.Getenv("GO_ENV")
	switch cases.Lower(language.English).String(c.App.Env) {
	case "development":
		c.App.Env = "development"
	case "staging":
		c.App.Env = "staging"
	case "testing":
		c.App.Env = "testing"
	case "production":
		c.App.Env = "production"
	default:
		c.App.Env = "development"
	}
}

func (c *AppConfig) initHttp() {
	c.Http.Port = os.Getenv("HTTP_PORT")
	if c.Http.Port == "" {
		c.Http.Port = "8080"
	}
}

func (c *AppConfig) initPostgres() {
	c.Postgres.Host = os.Getenv("DB_HOST")
	c.Postgres.Port = os.Getenv("DB_PORT")
	c.Postgres.User = os.Getenv("DB_USER")
	c.Postgres.Pass = os.Getenv("DB_PASS")
	c.Postgres.Name = os.Getenv("DB_NAME")
	c.Postgres.SSL = os.Getenv("DB_SSL_MODE")
}

func (c *AppConfig) initJwt() {
	c.Jwt.Secret = os.Getenv("JWT_SECRET")
	if c.Jwt.Secret == "" {
		c.Jwt.Secret = "secret"
	}
}

func (c *AppConfig) initOtel() {
	c.Otel.OTLPEndpoint = os.Getenv("OTEL_ENDPOINT")
	if c.Otel.OTLPEndpoint == "" {
		c.Otel.OTLPEndpoint = "localhost:4317"
	}
}

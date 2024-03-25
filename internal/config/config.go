package config

import (
	"embed"
	"fmt"
	"log"
	"os"
	"time"

	logging "gitlab.com/tour/internal/pkg/logger"
	"gitlab.com/tour/pkg/logger"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type AppMode string

const (
	DEVELOPMENT AppMode = "DEVELOPMENT"
	PRODUCTION  AppMode = "PRODUCTION"
)

var (
	TimeoutDuration time.Duration
	CacheTimeout    time.Duration
)

//go:embed configs
var configs embed.FS

type EmailNotificationConfig struct {
	RandomCodeLiveTime        string `yaml:"random_code_live_time"`
	RandomCodeErrorBlockTime  string `yaml:"random_code_error_block_time"`
	RandomCodeErrorTriesCount int    `yaml:"random_code_error_tries_count"`

	SenderEmail    string `env:"SENDER_EMAIL"`
	SenderPassword string `env:"SENDER_PASSWORD"`

	ProductName string `env:"PRODUCT_NAME"`

	SmtpHost string `yaml:"smtp_host"`
	SmtpPort string `yaml:"smtp_port"`
}

type Config struct {
	Logging logger.LoggingConfig `yaml:"logging"`
	Mode    string               `env:"APPLICATION_MODE" envDefault:"development"`

	Project struct {
		Name           string        `env:"PROJECT_NAME" yaml:"name"`
		Version        string        `env:"PROJECT_VERSION" yaml:"version"`
		Timeout        time.Duration `env:"PROJECT_TIMEOUT" yaml:"timeout"`
		SwaggerEnabled bool          `env:"PROJECT_SWAGGER_ENABLED" yaml:"swagger_enabled"`
		CacheTimeout   time.Duration `env:"PROJECT_CACHE_TIMEOUT" yaml:"cache_timeout"`
	} `yaml:"project"`

	Http struct {
		Host string `env:"HTTP_HOST" yaml:"host"`
		Port int    `env:"HTTP_PORT" yaml:"port"`

		URL string `env:"HTTP_URL" yaml:"url"`
	} `yaml:"http"`

	Grpc struct {
		Host string `env:"GRPC_HOST" yaml:"host"`
		Port int    `env:"GRPC_PORT" yaml:"port"`

		ChatGrpcHost string `env:"CHAT_SERVICE_HOST"`
		ChatGrpcPort int    `env:"CHAT_SERVICE_PORT"`

		URL string `env:"GRPC_URL" yaml:"url"`
	} `yaml:"grpc"`

	JWT struct {
		AccessLifeTime            int    `yaml:"access_life_time"`
		RefreshLifeTime           int    `yaml:"refresh_life_time"`
		OrganizationTokenLifeTime string `yaml:"organization_token_life_time"` // todo: implement later
		SecretToken               string `env:"TOKEN_SECRET"`
	} `yaml:"jwt"`

	Organization struct {
		OrganizationKeyLifeTime string `yaml:"organization_token_life_time"`
		OrganizationSecretKey   string `env:"ORGANIZATION_KEY_GENERATOR_SECRET"`
	} `yaml:"organization"`

	PSQL struct {
		URI string `env:"PSQL_URI"`
	}
	RedisAddress string `env:"REDIS_ADDRESS"`

	EmailNotificationConfig EmailNotificationConfig `yaml:"notification_service"`

	Casbin struct {
		ConfigPath string `env:" ../transport/grpc/middleware/rbac.conf"`
		PolicyPath string `env:" ../transport/grpc/middleware/policy_effect.csv"`
	}
}

func Load() *Config {
	var cfg Config
	err := godotenv.Load(".env")
	if err != nil && !os.IsNotExist(err) {
		logging.Log.Fatal("failed to load .env file", zap.Error(err))
	}

	configPath, err := getConfigPath(AppMode(getAppMode()))
	if err != nil {
		panic(err)
	}

	file, err := configs.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		panic(err)
	}

	if err := env.Parse(&cfg); err != nil {
		log.Println(err.Error())
		panic("unmarshal from environment error")
	}

	TimeoutDuration = cfg.Project.Timeout

	cfg.MakeGrpcURL()
	cfg.MakeHttpURL()

	return &cfg
}

func getAppMode() AppMode {
	mode := AppMode(os.Getenv("APPLICATION_MODE"))

	if mode != DEVELOPMENT {
		mode = PRODUCTION
	}

	return mode
}

func getConfigPath(appMode AppMode) (string, error) {
	suffix := "dev"
	if appMode == PRODUCTION {
		suffix = "prod"
	}

	return fmt.Sprintf("configs/%s.yaml", suffix), nil
}

func (c *Config) MakeGrpcURL() {
	c.Grpc.URL = fmt.Sprintf("%s:%d", c.Grpc.Host, c.Grpc.Port)
}

func (c *Config) MakeHttpURL() {
	c.Http.URL = fmt.Sprintf("%s:%d", c.Http.Host, c.Http.Port)
}

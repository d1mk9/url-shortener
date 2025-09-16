package config

import (
	"log"
	"net/url"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	PostgresHost     string `mapstructure:"postgres_host" validate:"required,hostname|ip"`
	PostgresPort     string `mapstructure:"postgres_port" validate:"required,numeric"`
	PostgresDB       string `mapstructure:"postgres_db" validate:"required"`
	PostgresUser     string `mapstructure:"postgres_user" validate:"required"`
	PostgresPassword string `mapstructure:"postgres_password" validate:"required"`
	BaseURL          string `mapstructure:"base_url" validate:"required,url,endsnotwith=/"`
}

func (c Config) PostgresDSN() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.PostgresUser, c.PostgresPassword),
		Host:   c.PostgresHost + ":" + c.PostgresPort,
		Path:   "/" + c.PostgresDB,
	}
	q := url.Values{}
	q.Set("sslmode", "disable")
	u.RawQuery = q.Encode()
	return u.String()
}

func MustLoad() *Config {
	v := viper.New()

	cf := os.Getenv("CONFIG_FILE")
	if cf == "" {
		log.Fatal("CONFIG_FILE is required")
	}
	v.SetConfigFile(cf)

	_ = v.BindEnv("postgres_user", "POSTGRES_USER")
	_ = v.BindEnv("postgres_password", "POSTGRES_PASSWORD")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("config: cannot read %s: %v", cf, err)
	}

	var cfg Config
	if err := v.UnmarshalExact(&cfg); err != nil {
		log.Fatalf("config: unmarshal: %v", err)
	}

	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		log.Fatalf("config: validation failed: %v", err)
	}

	return &cfg
}

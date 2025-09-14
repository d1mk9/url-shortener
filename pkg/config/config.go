package config

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	PostgresHost     string `mapstructure:"postgres_host"`
	PostgresPort     string `mapstructure:"postgres_port"`
	PostgresDB       string `mapstructure:"postgres_db"`
	PostgresUser     string `mapstructure:"postgres_user"`
	PostgresPassword string `mapstructure:"postgres_password"`
	BaseURL          string `mapstructure:"base_url"`
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
	validate(&cfg)

	return &cfg
}

func validate(cfg *Config) {
	var missing []string

	if cfg.PostgresHost == "" {
		missing = append(missing, "postgres_host")
	}
	if cfg.PostgresPort == "" {
		missing = append(missing, "postgres_port")
	}
	if cfg.PostgresDB == "" {
		missing = append(missing, "postgres_db")
	}
	if cfg.PostgresUser == "" {
		missing = append(missing, "postgres_user / POSTGRES_USER")
	}
	if cfg.PostgresPassword == "" {
		missing = append(missing, "postgres_password / POSTGRES_PASSWORD")
	}
	if cfg.BaseURL == "" {
		missing = append(missing, "base_url")
	} else if strings.HasSuffix(cfg.BaseURL, "/") {
		missing = append(missing, "base_url must not end with '/'")
	}

	if len(missing) > 0 {
		log.Fatalf("config: missing/invalid keys:\n  - %s", strings.Join(missing, "\n  - "))
	}
}

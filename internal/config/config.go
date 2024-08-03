package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/viper"
)

const (
	defaultPath = "./local.yaml"
)

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-default:"local" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
	PostgresEnv
}

type PostgresEnv struct {
	PGHost     string
	PGUser     string
	PGPassword string
	PGName     string
	PGPort     string
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	TimeOut     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"100s"`
	User        string        `yaml:"user"`
	Password    string        `yaml:"password" env:"HTTP_SERVER_PASSWORD"`
}

func MustLoad() Config {

	if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
		log.Fatalf("config file %s doesn't exist", defaultPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(defaultPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	
	getPGEnv(&cfg)
	
	getAdminEnv(&cfg)
	
	return cfg
}

func getPGEnv(cfg *Config) {
	viper.AutomaticEnv()

	viper.SetEnvPrefix("POSTGRES")

	viper.SetDefault("USER", "admin")
	viper.SetDefault("PASSWORD", "admin")
	viper.SetDefault("DB", "url-shortener")
	viper.SetDefault("PORT", "5432")
	viper.SetDefault("HOST", "localhost")

	cfg.PGUser = viper.GetString("USER")
	cfg.PGPassword = viper.GetString("PASSWORD")
	cfg.PGName = viper.GetString("DB")
	cfg.PGPort = viper.GetString("PORT")
	cfg.PGHost = viper.GetString("HOST")

	viper.Reset()

}

func getAdminEnv(cfg *Config) {
	viper.AutomaticEnv()

	viper.SetDefault("USER_SHORTENER", "user")
	viper.SetDefault("PASSWORD_SHORTENER", "password")

	cfg.User = viper.GetString("USER_SHORTENER")
	cfg.Password = viper.GetString("PASSWORD_SHORTENER")
}

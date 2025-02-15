package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env    string `yaml:"env" env-default:"local"`
	DB     DB     `yaml:"db"`
	Server Server `yaml:"server"`
}

type Server struct {
	Port string `yaml:"port" env-required:"true"`
}

type DB struct {
	Username string `yaml:"username" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Host     string `yaml:"host" env-required:"true"`
	Port     string `yaml:"port" env-required:"true"`
	DBname   string `yaml:"dbname" env-required:"true"`
	SSLmode  string `yaml:"sslmode" env-default:"disable"`
}

func MustLoad() *Config {
	path := fetchConfig()
	if path == "" {
		panic("config path is empty")
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(cfgPath string) *Config {
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		panic("config file does not exists: " + cfgPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfig() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

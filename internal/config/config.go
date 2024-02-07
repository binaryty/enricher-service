package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	DefaultEnv            = "local"
	DefaultDBHost         = "localhost"
	DefaultDBPort         = "5432"
	DefaultDBUser         = "postgres"
	DefaultDBPass         = "postgres"
	DefaultDBName         = "people_db"
	DefaultDBSSLMode      = "disable"
	DefaultAgeAPI         = "https://api.agify.io"
	DefaultGenderAPI      = "https://api.genderize.io"
	DefaultNationalityAPI = "https://api.nationalize.io"
	DefaultHTTPAddress    = "localhost:8082"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"local"`
	DBHost     string     `yaml:"db_host" env-default:"localhost"`
	DBPort     string     `yaml:"db_port" env-default:"5432"`
	DBUser     string     `yaml:"db_user" env-required:"true"`
	DBPass     string     `yaml:"db_pass" env-required:"true"`
	DBName     string     `yaml:"db_name" env-default:"people_db"`
	DBSSLMode  string     `yaml:"db_ssl_mode" env-default:"disable"`
	API        API        `yaml:"api"`
	HTTPServer HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address string `yaml:"address" env-default:"localhost:8080"`
}

type API struct {
	Age         string `yaml:"age"`
	Gender      string `yaml:"gender"`
	Nationality string `yaml:"nationality"`
}

// MakePGURL make a connection for postgres.
func (c *Config) MakePGURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser,
		c.DBPass,
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.DBSSLMode,
	)
}

// MustLoad load config from config file and panic if not set.
func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("config path is empty")
	}

	return MustLoadFromEnv()
}

// MustLoadByPath load config from path and panic if not set.
func MustLoadByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file doesn't exist:" + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config:" + err.Error())
	}

	return &cfg
}

// fetchConfigPath returns path to config file read from flag or env var.
func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}

func MustLoadFromEnv() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf(".env file not founf")
	}

	return &Config{
		Env:       getEnv("ENV", DefaultEnv),
		DBHost:    getEnv("DB_HOST", DefaultDBHost),
		DBPort:    getEnv("DB_PORT", DefaultDBPort),
		DBUser:    getEnv("DB_USER", DefaultDBUser),
		DBPass:    getEnv("DB_PASS", DefaultDBPass),
		DBName:    getEnv("DB_NAME", DefaultDBName),
		DBSSLMode: getEnv("DB_SSL_MODE", DefaultDBSSLMode),
		API: API{
			Age:         getEnv("AGE_API", DefaultAgeAPI),
			Gender:      getEnv("GENDER_API", DefaultGenderAPI),
			Nationality: getEnv("NATIONALITY_API", DefaultNationalityAPI),
		},
		HTTPServer: HTTPServer{
			Address: getEnv("HTTP_HOST", DefaultHTTPAddress),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if val, exist := os.LookupEnv(key); exist {
		return val
	}

	return defaultVal
}

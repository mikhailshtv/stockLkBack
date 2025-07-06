package config

type Config struct {
	Host       string
	Port       int
	Secret     string
	Group      string
	DBName     string
	DBUsername string
	DBPassword string
}

var Cfg Config

func NewConfig() *Config {
	Cfg := Config{
		Host:       "localhost",
		Port:       8080,
		Group:      "/api/v1",
		DBName:     "stockLk",
		DBUsername: "mongo_user",
		DBPassword: "mongo_pass",
	}
	return &Cfg
}

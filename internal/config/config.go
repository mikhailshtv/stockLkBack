package config

type Config struct {
	Host       string
	Port       int
	Secret     string
	Group      string
	DbName     string
	DbUsername string
	DbPassword string
}

var Cfg Config

func NewConfig() *Config {
	Cfg := Config{
		Host:       "localhost",
		Port:       8080,
		Group:      "/api/v1",
		DbName:     "stockLk",
		DbUsername: "mongo_user",
		DbPassword: "mongo_pass",
	}
	return &Cfg
}

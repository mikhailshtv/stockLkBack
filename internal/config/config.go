package config

type Config struct {
	Host   string
	Port   int
	Secret string
	Group  string
}

var Cfg Config

func NewConfig() *Config {
	Cfg := Config{
		Host:  "localhost",
		Port:  8080,
		Group: "/api/v1",
	}
	return &Cfg
}

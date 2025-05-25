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
		Host:   "localhost",
		Port:   8080,
		Group:  "/api/v1",
		Secret: "2k0935h84j39k2ks9df8h4fj3dk2s02kj9f8h4g5",
	}
	return &Cfg
}

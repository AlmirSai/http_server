package config

type Config struct {
	Port     string
	CertFile string
	KeyFile  string
}

func NewConfig() *Config {
	return &Config{
		Port:     "8443",
		CertFile: "certs/server.crt",
		KeyFile:  "certs/server.key",
	}
}

package pkg

type Config struct {
	Port int `yaml:"port"`
}

func NewDefaultConfig() *Config {
	return &Config{
		Port: 36000,
	}
}

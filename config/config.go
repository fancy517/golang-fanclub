package config

type Config struct {
	Server struct {
		Network string `yaml:"network"`
		Port    string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"database"`
	Smtp struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"smtp"`
}

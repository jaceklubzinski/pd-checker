package config

//Config application configuration
type Config struct {
	Port         string
	DatabasePath string
}

//NewConfig constructtor for Config struct
func NewConfig(port, databasePath string) *Config {
	return &Config{
		Port:         port,
		DatabasePath: databasePath,
	}
}

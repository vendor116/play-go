package config

type App struct {
	APIPServer APIServer `mapstructure:"api_server"`
	LogLevel   string    `mapstructure:"log_level"`
}

type APIServer struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

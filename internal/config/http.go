package config

type HTTP struct {
	Host string `mapstructure:"host" default:""`
	Port string `mapstructure:"port" default:"8080"`
}

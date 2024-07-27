package config

type Cache struct {
	Size string `mapstructure:"size"`
	Path string `mapstructure:"path"`
}

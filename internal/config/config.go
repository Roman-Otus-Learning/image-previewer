package config

import "net"

type Config struct {
	Logger Logger
	HTTP   HTTP
	Client Client
	Cache  Cache
}

func (c *Config) HTTPAddr() string {
	return net.JoinHostPort(c.HTTP.Host, c.HTTP.Port)
}

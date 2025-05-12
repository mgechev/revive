package fixtures

import "strconv"

type options struct{}

func NewOptions() *options {
	return &options{}
}

type port uint16

func ToPort(s string) (port, bool) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return port(i), true
}

type Config struct {
	p port
}

func (c *Config) Port() port {
	return c.p
}

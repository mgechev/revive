package fixtures

import "strconv"

type options struct{}

func NewOptions() *options { // MATCH /exported func NewOptions returns unexported type *fixtures.options, which can be annoying to use/
	return &options{}
}

type port uint16

func ToPort(s string) (port, bool) { // MATCH /exported func ToPort returns unexported type fixtures.port, which can be annoying to use/
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return port(i), true
}

type Config struct {
	p port
}

func (c *Config) Port() port { // MATCH /exported method Port returns unexported type fixtures.port, which can be annoying to use/
	return c.p
}

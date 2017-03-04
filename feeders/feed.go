package drivers

import (
	"fmt"
	"strings"
)

type Config map[string]string
type StatusF func(mssg string) (bool, error)
type MssgF func(map[string]string) error

type Feed interface {
	Run(statusF StatusF, mssgF MssgF) error
}

func NewConfig(s string) (*Config, error) {

	result := strings.Split(s, "::")
	if len(result)%2 != 0 {
		return nil, fmt.Errorf("bad config string")
	}

	m := Config{}
	for i := 0; i < len(result); i += 2 {
		m[result[i]] = result[i+1]
	}

	return &m, nil
}

func (c *Config) String() string {
	ss := []string{}

	for k, v := range *c {
		ss = append(ss, k, ":", v)
	}

	return strings.Join(ss, "::")
}

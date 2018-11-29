package log

import (
	"strconv"
)

const (
	CFG_ENABLE = "enable"
	CFG_LEVEL  = "level"
	CFG_FORMAT = "format"
)

// Config 通用的Channel配置
type Config struct {
	Enable  bool
	Level   int
	Format  string
	pattern *Pattern
}

func (c *Config) Init() {
	c.Enable = true
	c.Level = LEVEL_TRACE
	c.SetFormat(DefaultPattern)
}

func (c *Config) SetFormat(format string) {
	if c.pattern == nil {
		c.pattern = NewPattern()
	}

	c.Format = format
	c.pattern.Parse(format)
}

func (c *Config) SetConfig(key string, val string) bool {
	switch key {
	case "enable":
		c.Enable, _ = strconv.ParseBool(val)
		return true
	case "level":
		level, ok := ParseLevel(val)
		if ok {
			c.Level = level
		}
		return true
	case "format":
		c.SetFormat(val)
		return true
	}

	return false
}

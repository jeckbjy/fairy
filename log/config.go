package log

import (
	"strconv"
)

const (
	CFG_ENABLE = "enable"
	CFG_LEVEL  = "level"
	CFG_FORMAT = "format"
)

type Config struct {
	Enable  bool
	Level   int
	Format  string
	pattern *Pattern
}

func (self *Config) Init() {
	self.Enable = true
	self.Level = LEVEL_TRACE
}

func (self *Config) SetFormat(format string) {
	if self.pattern == nil {
		self.pattern = NewPattern()
	}

	self.Format = format
	self.pattern.Parse(format)
}

func (self *Config) SetConfig(key string, val string) bool {
	switch key {
	case "enable":
		self.Enable, _ = strconv.ParseBool(val)
		return true
	case "level":
		level, ok := ParseLevel(val)
		if ok {
			self.Level = level
		}
		return true
	case "format":
		self.SetFormat(val)
		return true
	}

	return false
}

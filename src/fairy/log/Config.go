package log

import (
	"strconv"
	"strings"
)

const (
	LAYOUT_FLAT = iota // 平坦格式
	LAYOUT_JSON        // Json格式
)

const (
	CFG_ENABLE = "enable"
	CFG_LEVEL  = "level"
	CFG_LAYOUT = "layout"
	CFG_FORMAT = "format"
)

type Config struct {
	Enable  bool
	Level   int
	Layout  int // json or flat
	Format  string
	pattern *Pattern
}

func (self *Config) Init() {
	self.Enable = true
	self.Level = LEVEL_TRACE
	self.Layout = LAYOUT_FLAT
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
		self.Level = ParseLevel(val)
		return true
	case "layout":
		self.Layout = ParseLayout(val)
		return true
	case "format":
		self.Format = val
		return true
	}

	return false
}

func ParseLayout(str string) int {
	switch strings.ToLower(str) {
	case "flat":
		return LAYOUT_FLAT
	case "json":
		return LAYOUT_JSON
	default:
		return LAYOUT_FLAT
	}
}

func ParseLevel(str string) int {
	switch strings.ToLower(str) {
	// case "all":
	// 	return LEVEL_ALL
	case "trace":
		return LEVEL_TRACE
	case "debug":
		return LEVEL_DEBUG
	case "info":
		return LEVEL_INFO
	case "warn":
		return LEVEL_WARN
	case "error":
		return LEVEL_ERROR
	case "fatal":
		return LEVEL_FATAL
	case "off":
		return LEVEL_OFF
	default:
		return LEVEL_DEBUG
	}
}

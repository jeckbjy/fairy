package fairy

// Option 配置
type Option interface {
}

func WithTagOption(tag string) *TagOption {
	return &TagOption{Tag: tag}
}

// TagOption 传递Tag使用
type TagOption struct {
	Tag string
}

func WithSyncOption(flag bool) *SyncOption {
	return &SyncOption{Flag: flag}
}

// SyncOption 阻塞调用
type SyncOption struct {
	Flag bool
}

func WithReconnectOption(count, interval int) *ReconnectOption {
	return &ReconnectOption{Count: count, Interval: interval}
}

// WithCloseReconnectOption 关闭断线重现
func WithCloseReconnectOption() *ReconnectOption {
	return &ReconnectOption{Count: 0}
}

// ReconnectOption 重连配置
type ReconnectOption struct {
	Count    int // 重连次数,0标识无需断线重连
	Interval int // 重连间隔
}

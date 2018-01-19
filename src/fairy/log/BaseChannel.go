package log

type BaseChannel struct {
	Config
}

func (self *BaseChannel) Open() {

}

func (self *BaseChannel) Close() {

}

func (self *BaseChannel) SetProperty(key string, val interface{}) bool {
	return false
}

func (self *BaseChannel) GetConfig() *Config {
	return &self.Config
}

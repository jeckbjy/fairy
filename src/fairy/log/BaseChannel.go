package log

type BaseChannel struct {
	enable bool
	level  int
}

func (self *BaseChannel) Open() {

}

func (self *BaseChannel) Close() {

}

func (self *BaseChannel) SetEnable(enable bool) {
	self.enable = enable
}

func (self *BaseChannel) GetEnable() bool {
	return self.enable
}

func (self *BaseChannel) SetProperty(key string, val interface{}) {

}

func (self *BaseChannel) CanWrite(level int) bool {
	return level >= self.level
}

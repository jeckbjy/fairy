package fairy

import (
	"io/ioutil"
)

const TABLE_HEAD_LEN_MAX = 3

func ParseTable(data string, meta interface{}) interface{} {
	return nil
}

func ReadTable(path string, meta interface{}) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	ParseTable(string(data), meta)
}
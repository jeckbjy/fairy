package fairy

import (
	"encoding/csv"
	"fairy/util"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

/*
usage:
type Address struct {
	Id   int  `id:"id"`
	Name string
	Desc string
}

var gAddressVec []*Address
var gAddressMap map[int]*Address

func GetAddressVec() []*Address {
	return gAddressVec
}

func GetAddressById(id int) *Address {
	return gAddressMap[id]
}

func LoadAddress(path string) {
	gAddressVec = ReadTable(path, Address{}).([]*Address)
	for _, record := range gAddressVec {
		gAddressMap[record.Id] = record
	}
}
*/

const TABLE_HEAD_LEN_MAX = 3

type tabHead struct {
	mode  int    // 类型
	name  string // 名字
	field int    // field索引
}

func setField(field *reflect.Value, str string) error {
	// set value
	switch field.Kind() {
	case reflect.String:
		field.SetString(str)
	case reflect.Bool:
		if val, err := strconv.ParseBool(str); err == nil {
			field.SetBool(val)
		} else {
			return err
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val, err := util.ParseInt64(str); err == nil {
			field.SetInt(val)
		} else {
			return err
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val, err := util.ParseUint64(str); err == nil {
			field.SetUint(val)
		} else {
			return err
		}
	case reflect.Float32, reflect.Float64:
		if val, err := strconv.ParseFloat(str, 0); err == nil {
			field.SetFloat(val)
		} else {
			return err
		}
	case reflect.Slice:
		// split，必须是整数或者浮点数?
		tokens := util.SplitNum(str)
		if len(tokens) == 0 {
			return nil
		}
		// check
	case reflect.Map:
		tokens := util.SplitNum(str)
		if len(tokens) == 0 {
			return nil
		}

		if len(tokens)%2 == 1 {
			return fmt.Errorf("table map cell count must be even!%v", str)
		}
		//
	}
	return nil
}

func ParseTable(reader *csv.Reader, meta interface{}) (interface{}, error) {
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(lines) < TABLE_HEAD_LEN_MAX {
		return nil, nil
	}

	rtype := util.GetRealType(meta)
	ptype := reflect.PtrTo(rtype)

	// fileds reflect map
	fieldsMap := make(map[string]int)
	for i := 0; i < rtype.NumField(); i++ {
		field := rtype.Field(i)
		name := strings.ToLower(field.Name)
		fieldsMap[name] = i
	}

	// 解析头信息
	colNum := len(lines[0])
	heads := make([]tabHead, colNum, colNum)
	for i := 0; i < colNum; i++ {
		head := &heads[i]
		head.name = lines[1][i]
		// parse name:
		// array rule:Id*?? Id_1,Name_1,Id_2,Name_2
		// enum rule:Id[key1:1,key2:2]???
		// find field
		lowName := strings.ToLower(head.name)
		if index, ok := fieldsMap[lowName]; ok {
			head.field = index
		} else {
			head.field = -1
		}
	}

	// 读取数据
	recordList := reflect.MakeSlice(reflect.SliceOf(ptype), 0, len(lines)-3)

	// foreach lines
	for i := 3; i < len(lines); i++ {
		line := lines[i]
		record := reflect.New(rtype)

		col := util.MinInt(len(line), colNum)
		// foreach column
		for j := 0; j < col; j++ {
			fieldIndex := heads[j].field
			if fieldIndex == -1 {
				continue
			}
			// create field
			field := record.Elem().Field(fieldIndex)
			err := setField(&field, line[j])
			if err != nil {
				return nil, err
			}
		}

		recordList = reflect.Append(recordList, record)
	}

	return recordList.Interface(), nil
}

func ReadTable(path string, meta interface{}) interface{} {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(file)
	record, err := ParseTable(reader, meta)
	if err != nil {
		panic(err)
	}

	return record
}

func ReadTableFromString(data string, meta interface{}) interface{} {
	strReader := strings.NewReader(data)
	reader := csv.NewReader(strReader)
	record, err := ParseTable(reader, meta)
	if err != nil {
		panic(err)
	}

	return record
}

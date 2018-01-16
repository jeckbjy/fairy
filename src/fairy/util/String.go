package util

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func ParseInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 0)
}

func ParseUInt64(str string) (uint64, error) {
	return strconv.ParseUint(str, 10, 0)
}

// 使用特殊字符分隔成数字，分隔符可以是任何非数字字符，例如|,()等
// 例如:-1|2|3 => [1,2,3]  4,5,6 => [4,5,6]
// 支持负数,"-"开始
func SplitInt64(str string) []int64 {
	result := []int64{}
	last := 0
	for {
		start := -1
		// find start
		for i := last; i < len(str); {
			r, size := utf8.DecodeRuneInString(str[i:])
			if unicode.IsDigit(r) || r == '-' {
				start = i
				break
			}
			i += size
		}

		if start == -1 {
			break
		}

		// find end
		end := -1
		for i := start + 1; i < len(str); {
			r, size := utf8.DecodeRuneInString(str[i:])
			if !unicode.IsDigit(r) {
				end = i
				last = i + size
				break
			}
			i += size
		}

		if end == -1 {
			end = len(str)
			last = end
		}

		subStr := str[start:end]
		val, _ := ParseInt64(subStr)
		result = append(result, val)
	}

	return result
}

// -.不能做分隔符，其他符号可以做
func SplitNum(str string) []string {
	result := []string{}
	last := 0
	for {
		// find start
		start := -1
		for i := last; i < len(str); i++ {
			r, size := utf8.DecodeRuneInString(str[i:])
			if unicode.IsDigit(r) || r == '-' {
				start = i
				break
			}
			i += size
		}

		if start == -1 {
			break
		}

		// find end
		end := -1
		for i := start + 1; i < len(str); {
			r, size := utf8.DecodeRuneInString(str[i:])
			if !unicode.IsDigit(r) && r != '.' {
				end = i
				last = i + size
				break
			}

			i += size
		}

		if end == -1 {
			end = len(str)
			last = end
		}
		sub := str[start:end]
		result = append(result, sub)
	}

	return result
}

// SplitAny:split by any char
func SplitAny(str string, sep string) []string {
	result := []string{}
	last := 0
	for {
		// find start
		start := -1
		for i := last; i < len(str); i++ {
			if strings.IndexByte(sep, str[i]) == -1 {
				start = i
				break
			}
		}

		if start == -1 {
			break
		}

		// find end
		end := -1
		for i := start + 1; i < len(str); i++ {
			if strings.IndexByte(sep, str[i]) != -1 {
				end = i
				last = i + 1
				break
			}
		}

		if end == -1 {
			end = len(str)
			last = end
		}

		sub := str[start:end]
		result = append(result, sub)
	}

	return result
}

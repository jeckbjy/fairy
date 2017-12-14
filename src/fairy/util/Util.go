package util

import (
	"errors"
	"fmt"
)

func NewError(format string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(format, args))
}

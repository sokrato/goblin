package utils

import (
	"errors"
	"strconv"
)

var (
	ErrNoSuchKey = errors.New("no such key")
	ErrBadValue  = errors.New("bad value")
)

type Dict map[string]interface{}

func (d Dict) String(key string) (string, error) {
	val, ok := d[key]
	if !ok {
		return "", ErrNoSuchKey
	}
	valStr, ok := val.(string)
	if !ok {
		return "", ErrBadValue
	}
	return valStr, nil
}

func (d Dict) MustString(key string) string {
	val, err := d.String(key)
	if err != nil {
		panic(err)
	}
	return val
}

func (d Dict) Int(key string) (int, error) {
	val, ok := d[key]
	if !ok {
		return 0, ErrNoSuchKey
	}
	valInt, ok := val.(int)
	if !ok {
		return 0, ErrBadValue
	}
	return valInt, nil
}

func (d Dict) ParseInt(key string) (int, error) {
	val, ok := d[key]
	if !ok {
		return 0, ErrNoSuchKey
	}
	switch v := val.(type) {
	case string:
		return strconv.Atoi(v)
	case int:
		return v, nil
	default:
		return 0, ErrBadValue
	}
}

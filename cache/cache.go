package cache

import (
	"errors"
	"github.com/dlutxx/goblin/utils"
)

var (
	cacheFactories              = map[string]func(utils.Dict) Cache{}
	ErrFactoryAlreadyRegistered = errors.New("cache factory already registered")
)

type Cache interface {
	Get(string) interface{}
	Set(string, interface{})
	Del(string)
	Incr(string)
	Decr(string)
}

func Register(key string, factory func(utils.Dict) Cache) error {
	if _, ok := cacheFactories[key]; ok {
		return ErrFactoryAlreadyRegistered
	}
	cacheFactories[key] = factory
	return nil
}

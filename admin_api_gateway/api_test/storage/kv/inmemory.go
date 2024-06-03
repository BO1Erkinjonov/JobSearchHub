package kv

import (
	"admin_api_gateway/api_test/storage"
	"errors"
	"sync"
)

type InMemory struct {
	inst sync.Map
}

func NewInMemory() *InMemory {
	return &InMemory{}
}

func (im *InMemory) Set(key string, value *storage.User) error {
	im.inst.Store(key, value)
	return nil
}

func (im *InMemory) Get(key string) (*storage.User, error) {
	value, exists := im.inst.Load(key)
	if !exists {
		return nil, errors.New("not found error")
	}

	return value.(*storage.User), nil
}

func (im *InMemory) Delete(key string) error {
	im.inst.Delete(key)
	return nil
}

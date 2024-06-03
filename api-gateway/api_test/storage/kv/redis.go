package kv

import (
	"api-gateway/api_test/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Redigo struct {
	rdb *redis.Client
}

func NewRedisInit(client *redis.Client) *Redigo {
	return &Redigo{rdb: client}
}

func (r *Redigo) Set(key string, value *storage.User) error {
	str, err := json.Marshal(value)
	if err != nil {
		return err
	}
	fmt.Println(string(str))
	if err := r.rdb.Set(context.Background(), key, str, 0); err != nil {
		return err.Err()
	}
	return nil
}

func (r *Redigo) Delete(key string) error {
	if err := r.rdb.Del(context.Background(), key); err != nil {
		return err.Err()
	}
	return nil
}

func (r *Redigo) Get(key string) (*storage.User, error) {

	str, err := r.rdb.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	var user *storage.User
	err = json.Unmarshal([]byte(str), &user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

package sessionstores

import (
	"errors"
	"github.com/fzzy/radix/redis"
	"log"
)

// Redis session store. See http://redis.io
type Redis struct {
	address, password string
	client            *redis.Client
}

// NewRedis creates a new pointer to Redis instance.
// First argument is *Redis address.
// Second argument is optional Redis password.
func NewRedis(args ...string) *Redis {
	if len(args) == 0 {
		log.Panicln(errors.New("Address is required"))
	}
	address, password := args[0], ""
	if len(args) > 1 {
		password = args[1]
	}
	return &Redis{address: address, password: password}
}

func (s *Redis) Connect() error {
	client, err := redis.Dial("tcp", s.address)
	if err != nil {
		return err
	}
	if s.password != "" {
		err = client.Cmd("AUTH", s.password).Err
		if err != nil {
			return err
		}
	}
	s.client = client
	return err
}

func (s *Redis) expire(key string) {
	s.client.Cmd("EXPIRE", key, 300)
}

func (s *Redis) SetValue(key, value string) error {
	err := s.client.Cmd("SET", key, value).Err
	if err != nil {
		return err
	}
	s.expire(key)
	return nil
}

func (s *Redis) GetValue(key string) (string, error) {
	str, err := s.client.Cmd("GET", key).Str()
	if err != nil {
		return str, err
	}
	s.expire(key)
	return str, err
}

func (s *Redis) ValueExists(key string) (bool, error) {
	return s.client.Cmd("EXISTS", key).Bool()
}

func (s *Redis) DeleteValue(key string) error {
	return s.client.Cmd("DEL", key).Err
}

func (s *Redis) HashSetValue(name, key, value string) error {
	err := s.client.Cmd("HSET", name, key, value).Err
	if err != nil {
		return err
	}
	s.expire(name)
	return nil
}

func (s *Redis) HashGetValue(name, key string) (string, error) {
	str, err := s.client.Cmd("HGET", name, key).Str()
	if err != nil {
		return str, err
	}
	s.expire(name)
	return str, err
}

func (s *Redis) HashValueExists(name, key string) (bool, error) {
	return s.client.Cmd("HEXISTS", name, key).Bool()
}

func (s *Redis) HashDeleteValue(name, key string) error {
	return s.client.Cmd("HDEL", name, key).Err
}

func (s *Redis) HashExists(name string) (bool, error) {
	return s.ValueExists(name)
}

func (s *Redis) HashDelete(name string) error {
	return s.DeleteValue(name)
}

func (s *Redis) Close() error {
	return s.client.Close()
}

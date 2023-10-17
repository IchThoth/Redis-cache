package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

type redisConfig struct {
	host     string
	password string
	prefix   string
}

type Cache interface {
	Has(string) (bool, error)
	Get(string) (interface{}, error)
	Set(string, interface{}, ...int) error
	Forget(string) error
	EmptyByMatch(string) error
	Empty() error
}

type RedisCache struct {
	Conn      *redis.Pool
	Prefix    string
	HasPrefix bool
}

type Entry map[string]interface{}

var redconfig *redisConfig

func CreateRedisPool() *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redconfig.host, redis.DialPassword(redconfig.password))
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
		MaxIdle:         50,
		MaxActive:       10000,
		IdleTimeout:     240,
		Wait:            false,
		MaxConnLifetime: 0,
	}
}

func createClientRedis() *RedisCache {
	redconfig = &redisConfig{
		host:     os.Getenv("REDIS_HOST"),
		password: os.Getenv("REDIS_PASSWORD"),
		prefix:   os.Getenv("REDIS_PREFIX"),
	}
	cache := RedisCache{
		Conn:   CreateRedisPool(),
		Prefix: redconfig.prefix,
	}
	return &cache
}

func encode(item Entry) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(item)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func decode(val string) (Entry, error) {
	item := Entry{}
	b := bytes.Buffer{}
	b.Write([]byte(val))
	d := gob.NewDecoder(&b)
	err := d.Decode(&item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (c *RedisCache) Has(val string) (bool, error) {
	key := fmt.Sprintf("%s:%s", c.Prefix, val)
	connect := c.Conn.Get()

	defer connect.Close()

	ok, err := redis.Bool(connect.Do("EXISTS", key))
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (c *RedisCache) Get(val string) (interface{}, error) {
	key := fmt.Sprintf("%s:%s", c.Prefix, val)
	connect := c.Conn.Get()

	defer connect.Close()

	cachedKey, err := redis.Bytes(connect.Do("GET", key))
	if err != nil {
		return nil, err
	}

	decode, err := decode(string(cachedKey))
	if err != nil {
		return nil, err
	}

	item := decode[key]

	return item, nil
}

func (c *RedisCache) Set(val string, data interface{}, exp ...int) error {
	key := fmt.Sprintf("%s:%s", c.Prefix, val)
	connect := c.Conn.Get()

	defer connect.Close()

	entry := Entry{}
	entry[key] = data

	encode, err := encode(entry)
	if err != nil {
		return err
	}

	if len(exp) > 0 {
		_, err = connect.Do("SETEX", key, exp[0], string(encode))
		if err != nil {
			return err
		}
	} else {
		_, err = connect.Do("SETEX", key, exp[0])
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *RedisCache) Forget(val string) error {
	key := fmt.Sprintf("%s:%s", c.Prefix, val)
	connect := c.Conn.Get()

	defer connect.Close()

	_, err := connect.Do("DEL", key)
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisCache) getKeys(pattern string) ([]string, error) {
	connect := c.Conn.Get()

	defer connect.Close()

	iterator := 0
	keys := []string{}

	for {
		arr, err := redis.Values(connect.Do("SCAN", iterator, "MATCH", fmt.Sprintf("%s*", pattern)))
		if err != nil {
			return keys, err
		}
		iterator, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iterator == 0 {
			break
		}
	}

	return keys, nil
}

func (c *RedisCache) EmptyByMatch(val string) error {
	key := fmt.Sprintf("%s:%s", c.Prefix, val)
	connect := c.Conn.Get()

	defer connect.Close()

	keys, err := c.getKeys(key)
	if err != nil {
		return err
	}

	for _, x := range keys {
		err := c.Forget(x)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *RedisCache) Empty() error {
	key := fmt.Sprintf("%s:", c.Prefix)
	connect := c.Conn.Get()

	defer connect.Close()

	keys, err := c.getKeys(key)
	if err != nil {
		return err
	}

	for _, x := range keys {
		err = c.Forget(x)
		if err != nil {
			return err
		}
	}
	return nil
}
func main() {
	createClientRedis()
	CreateRedisPool()
}

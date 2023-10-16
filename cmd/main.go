package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/ichthoth/Redis-cache/cache"
	"github.com/ichthoth/Redis-cache/routes"
)

type Application struct {
	name   string
	env    string
	port   string
	config config
}

type config struct {
	redis redisConfig
}

type redisConfig struct {
	host     string
	password string
	prefix   string
}

var redconfig *config

func (c *Application) CreateRedisPool() *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redconfig.redis.host, redis.DialPassword(redconfig.redis.password))
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

func (c *Application) createClientRedis() *cache.RedisCache {
	redconfig = &config{
		redis: redisConfig{
			host:     os.Getenv("REDIS_HOST"),
			password: os.Getenv("REDIS_PASSWORD"),
			prefix:   os.Getenv("REDIS_PREFIX"),
		},
	}
	cache := cache.RedisCache{
		Conn:      c.CreateRedisPool(),
		Prefix:    c.config.redis.prefix,
		HasPrefix: false,
	}
	return &cache
}

func Run() error {
	r := gin.Default()

	routes.UserRoutes(r)

	mount := r.Run()

	return mount
}

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}

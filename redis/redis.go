package redis

import (
	"errors"
	"log"
	"os"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

type Redis struct {
	RedisClient redis.Client
}

func NewRedis() Redis {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error readng env file")
	}

	RedisPassword := os.Getenv("REDIS_PASSWORD")

	var client = redis.NewClient(&redis.Options{
		// Container name + port since we are using docker
		Addr:     "redis:6379",
		Password: RedisPassword,
	})

	if client == nil {
		errors.New("Cannot run redis")
	}

	return Redis{
		RedisClient: *client,
	}
}

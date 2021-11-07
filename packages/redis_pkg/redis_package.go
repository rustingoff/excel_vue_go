package redis_pkg

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

func NewRedisConnection() *redis.Client {

	redisDB := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
	})

	ctx := context.Background()
	err := redisDB.Ping(ctx).Err()

	if err != nil {
		log.Println("Cant connect to REDIS")
		panic(err.Error())
	}

	log.Println("Successfully connected to REDIS.")
	return redisDB

}

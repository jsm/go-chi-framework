package connector

import (
	"log"

	"github.com/go-redis/redis"
)

// ConnectRedis initializes a Redis connection
func ConnectRedis(address string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: address,
	})

	_, err := client.Ping().Result()

	if err != nil {
		log.Println("Failed to connect to Redis")
		log.Fatalln(err)
	}

	return client
}

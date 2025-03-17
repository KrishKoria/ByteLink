package store

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type StoreService struct {
	client *redis.Client
}

var (
	ctx     = context.Background()
	service = &StoreService{}
)

const CacheDuration = 1 * time.Hour

func InitializeStoreService() *StoreService {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}
	fmt.Printf("Redis Pong: %s\n", pong)
	service.client = redisClient
	return service
}

func SaveMapping(shortUrl string, longUrl string) {
	err := service.client.Set(ctx, shortUrl, longUrl, CacheDuration).Err()
	if err != nil {
		panic(fmt.Sprintf("Failed to save mapping: %v", err))
	}
}

func GetLongUrl(shortUrl string) string {
	longUrl, err := service.client.Get(ctx, shortUrl).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to get long URL: %v", err))
	}
	return longUrl
}

package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/KrishKoria/ByteLink/internal/database"
	"github.com/KrishKoria/ByteLink/miscellaneous"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
	"time"
)

type StoreService struct {
	client          *redis.Client
	db              *database.Queries
	rawDB           *sql.DB
	cleanupStopChan chan bool
}

type URLMapping struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}
type URLStats struct {
	ShortURL   string `json:"short_url"`
	LongURL    string `json:"long_url"`
	ClickCount int64  `json:"click_count"`
}

var (
	ctx     = context.Background()
	service = &StoreService{}
)

const CacheDuration = 1 * time.Hour

func InitializeStoreService() *StoreService {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}
	fmt.Printf("Redis Pong: %s\n", pong)
	service.client = redisClient

	sqlDB, err := sql.Open("sqlite3", "./bytelink.sqlite")
	if err != nil {
		panic(fmt.Sprintf("Failed to open database: %v", err))
	}
	queries := database.New(sqlDB)
	service.db = queries
	service.rawDB = sqlDB
	fmt.Println("Connected to SQLite database")

	stopCleanupChan := miscellaneous.StartCleanupJob(service, CacheDuration)
	service.cleanupStopChan = stopCleanupChan
	return service
}

func GetDBConn() *sql.DB {
	if service == nil || service.rawDB == nil {
		panic("Store service not initialized")
	}
	return service.rawDB
}

func SaveMapping(shortUrl string, longUrl string, userId string) error {
	var urlId uuid.UUID
	existingId, err := service.db.GetURLIdByLongURL(ctx, longUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			urlId = uuid.New()
			err = service.db.SaveURL(ctx, database.SaveURLParams{
				ID:      urlId,
				LongUrl: longUrl,
			})
			if err != nil {
				panic(fmt.Sprintf("Failed to save URL: %v", err))
			}
		} else {
			panic(fmt.Sprintf("Failed to query URL: %v", err))
		}
	} else {
		if uuidStr, ok := existingId.(string); ok {
			parsedUUID, parseErr := uuid.Parse(uuidStr)
			if parseErr != nil {
				panic(fmt.Sprintf("Failed to parse UUID: %v", parseErr))
			}
			urlId = parsedUUID
			_, err = service.db.GetMappingByUserIDAndUrlID(ctx, database.GetMappingByUserIDAndUrlIDParams{
				UserID: userId,
				UrlID:  urlId,
			})
			if err == nil {
				fmt.Printf("Mapping already exists for user %s and URL %s\n", userId, longUrl)
				return nil
			}
		} else if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to check existing mapping: %w", err)
		}
	}
	mappingId := uuid.New()
	err = service.db.SaveMapping(ctx, database.SaveMappingParams{
		ID:       mappingId,
		ShortUrl: shortUrl,
		UrlID:    urlId,
		UserID:   userId,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to save mapping: %v", err))
	}
	userCacheKey := fmt.Sprintf("%s:%s", shortUrl, userId)
	service.client.Set(ctx, userCacheKey, longUrl, CacheDuration)
	service.client.Set(ctx, shortUrl, longUrl, CacheDuration)

	fmt.Printf("Saved mapping: %s -> %s\n", shortUrl, longUrl)
	return nil
}

func GetLongUrl(shortUrl string, userId string) string {
	var longUrl string

	cacheKey := fmt.Sprintf("%s:%s", shortUrl, userId)
	longUrl, err := service.client.Get(ctx, cacheKey).Result()
	if err == nil {
		fmt.Printf("Cache hit: Retrieved URL from Redis: %s -> %s\n", shortUrl, longUrl)
		return longUrl
	}

	dbLongUrl, err := service.db.GetLongURLByShortURLAndUserID(ctx, database.GetLongURLByShortURLAndUserIDParams{
		ShortUrl: shortUrl,
		UserID:   userId,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("No URL mapping found for %s and user %s\n", shortUrl, userId)
			return ""
		}
		panic(fmt.Sprintf("Error querying database: %v", err))
	}
	err = service.client.Set(ctx, cacheKey, dbLongUrl, CacheDuration).Err()
	if err != nil {
		fmt.Printf("Warning: Failed to cache result in Redis: %v\n", err)
	}

	fmt.Printf("Cache miss: Retrieved URL from SQLite and cached in Redis: %s -> %s\n", shortUrl, dbLongUrl)
	return dbLongUrl
}

func GetLongUrlPublic(shortUrl string) string {
	longUrl, err := service.client.Get(ctx, shortUrl).Result()
	if err == nil {
		return longUrl
	}

	dbLongUrl, err := service.db.GetLongURLByShortURL(ctx, shortUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ""
		}
		panic(fmt.Sprintf("Error querying database: %v", err))
	}

	err = service.client.Set(ctx, shortUrl, dbLongUrl, CacheDuration).Err()
	if err != nil {
		fmt.Printf("Warning: Failed to cache result in Redis: %v\n", err)
	}

	return dbLongUrl
}

func GetMappingsByUserID(userId string) []URLMapping {
	if userId == "" {
		return []URLMapping{}
	}

	mappings, err := service.db.GetMappingsByUserID(ctx, userId)
	if err != nil {
		fmt.Printf("Error retrieving mappings for user %s: %v\n", userId, err)
		return []URLMapping{}
	}

	result := make([]URLMapping, len(mappings))
	for i, mapping := range mappings {
		result[i] = URLMapping{
			ShortURL: mapping.ShortUrl,
			LongURL:  mapping.LongUrl,
		}
	}

	return result
}

func DeleteMapping(shortURL string, userID string) error {
	_, err := service.db.GetMappingByShortURLAndUserID(ctx, database.GetMappingByShortURLAndUserIDParams{
		ShortUrl: shortURL,
		UserID:   userID,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no mapping found for short URL %s and user %s", shortURL, userID)
		}
		return fmt.Errorf("failed to check mapping existence: %w", err)
	}

	err = service.db.DeleteMappingByShortURLAndUserID(ctx, database.DeleteMappingByShortURLAndUserIDParams{
		ShortUrl: shortURL,
		UserID:   userID,
	})

	if err != nil {
		return fmt.Errorf("failed to delete mapping: %w", err)
	}

	userCacheKey := fmt.Sprintf("%s:%s", shortURL, userID)
	service.client.Del(ctx, userCacheKey)
	service.client.Del(ctx, shortURL)

	return nil
}

func (s *StoreService) GetOrphanedURLs(ctx context.Context) ([]string, error) {
	orphanedURLsRaw, err := s.db.GetOrphanedURLs(ctx)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	result := make([]string, len(orphanedURLsRaw))
	for i, urlID := range orphanedURLsRaw {
		switch id := urlID.(type) {
		case string:
			result[i] = id
		case uuid.UUID:
			result[i] = id.String()
		default:
			result[i] = fmt.Sprintf("%v", id)
		}
	}

	return result, nil
}

func (s *StoreService) DeleteURLByID(ctx context.Context, id string) error {
	return s.db.DeleteURLByID(ctx, id)
}

func GetURLStats(ctx context.Context, shortURL string, userID string) (URLStats, error) {
	stats, err := service.db.GetURLStatsForUser(ctx, database.GetURLStatsForUserParams{
		ShortUrl: shortURL,
		UserID:   userID,
	})

	if err != nil {
		return URLStats{}, fmt.Errorf("failed to get URL stats: %w", err)
	}
	var clickCount int64
	if stats.ClickCount.Valid {
		clickCount = stats.ClickCount.Int64
	}

	return URLStats{
		ShortURL:   stats.ShortUrl,
		LongURL:    stats.LongUrl,
		ClickCount: clickCount,
	}, nil
}

func IncrementClickCount(ctx context.Context, shortURL string) error {
	err := service.db.IncrementClickCount(ctx, shortURL)
	if err != nil {
		return fmt.Errorf("failed to increment click count: %w", err)
	}
	return nil
}

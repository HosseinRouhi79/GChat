package cache

import (
	"chat/config"
	"chat/src/api/pkg/logging"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

type single struct{}

var cfg = config.GetConfig()

var logger = logging.NewLogger(cfg)

var singleton *single

var redisClient *redis.Client

func InitRedis(cfg *config.Config) *single {
	if singleton == nil {
		fmt.Println("creating redis connection...")
		redisClient = redis.NewClient(&redis.Options{
			Addr:               fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
			Password:           cfg.Redis.Password,
			DB:                 0,
			DialTimeout:        cfg.Redis.DialTimeout * time.Second,
			ReadTimeout:        cfg.Redis.ReadTimeout * time.Second,
			WriteTimeout:       cfg.Redis.WriteTimeout * time.Second,
			PoolSize:           cfg.Redis.PoolSize,
			PoolTimeout:        cfg.Redis.PoolTimeout,
			IdleTimeout:        500 * time.Millisecond,
			IdleCheckFrequency: cfg.Redis.IdleCheckFrequency * time.Millisecond,
		})
		singleton = &single{}
		logger.Info(logging.Redis, logging.Migration, "Redis connected successfully", nil)
		return singleton
	} else {
		logger.Info(logging.Redis, logging.Migration, "Redis is already configured", nil)
		return singleton
	}
}

func GetRedis() *redis.Client {
	return redisClient
}

func CloseRedis() {
	redisClient.Close()
}

func Set[T any](redisClient *redis.Client, key string, value T, durationTime time.Duration) error {
	logger := logging.NewLogger(config.GetConfig())
	v, err := json.Marshal(value)
	if err != nil {
		logger.Infof("failed to marshal:%v", err)
		return err
	}
	if err := redisClient.Set(key, v, durationTime).Err(); err != nil {
		logger.Infof("failed to set data in Redis: %v", err)
		return err
	}
	logger.Info(logging.Internal, logging.Insert, "data set successfully", nil)
	return nil
}

func Get[T any](redisClient *redis.Client, key string) (dest T, err error) {
	val, err := redisClient.Get(key).Result()
	if err != nil {
		logger.Infof("failed to get from redis:%v", err)
		return dest, err
	}
	err = json.Unmarshal([]byte(val), &dest)
	if err != nil {
		logger.Infof("failed to unmarshal:%v", err)
		return dest, err
	}
	logger.Infof("result is:%v", dest)
	return dest, nil
}

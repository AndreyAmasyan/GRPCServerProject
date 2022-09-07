package redisStorage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"test/internal/data"
	"time"

	"github.com/go-redis/redis/v9"
)

const (
	hashKey string = "usershash"
)

// RedisStorage ...
type RedisStorage struct {
	rdb     *redis.Client
	hashKey string
}

func New() (*RedisStorage, error) {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		fmt.Println("RedisStorage New strconv.Atoi", err)
		db = 0
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	return &RedisStorage{
		rdb:     client,
		hashKey: hashKey,
	}, nil
}

func (s *RedisStorage) SetAllUsers(ctx context.Context, list []data.User) error {
	var keyValuePairs []interface{}

	for _, v := range list {
		key := strconv.Itoa(v.ID)
		b, err := json.Marshal(v)
		if err != nil {
			fmt.Println("RedisStorage SetAllUser json.Marshal", err)
			return err
		}
		value := string(b)

		keyValuePairs = append(keyValuePairs, key, value)
	}

	if err := s.rdb.HSet(ctx, s.hashKey, keyValuePairs...).Err(); err != nil {
		fmt.Println("RedisStorage SetAllUsers s.rdb.HSet", err)
		return err
	}

	_, err := s.rdb.Expire(ctx, s.hashKey, 60*time.Second).Result()
	if err != nil {
		fmt.Println("RedisStorage SetAllUsers s.rdb.Expire", err)
		return err
	}

	return nil
}

func (s *RedisStorage) IsExists(ctx context.Context) (bool, error) {
	val, err := s.rdb.Exists(ctx, s.hashKey).Result()
	if err != nil {
		fmt.Println("RedisStorage IsExists s.rdb.Exists", err)
		return false, err
	}

	if val == 1 {
		fmt.Println("Cache is exists")
		return true, nil
	}

	fmt.Println("Cache is not exists")

	return false, nil
}

func (s *RedisStorage) GetAllUsers(ctx context.Context) ([]data.User, error) {
	var users []data.User

	m, err := s.rdb.HGetAll(ctx, s.hashKey).Result()
	if err != nil {
		fmt.Println("RedisStorage GetAllUsers s.rdb.HGetAll", err)
		return nil, err
	}

	for _, v := range m {
		var user data.User

		if err := json.Unmarshal([]byte(v), &user); err != nil {
			fmt.Println("RedisStorage GetAllUsers json.Unmarshal", err)
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil

}

func (s *RedisStorage) Close() error {
	if err := s.rdb.Close(); err != nil {
		fmt.Println("RedisStorage Close s.rdb.Close", err)
		return err
	}

	return nil
}

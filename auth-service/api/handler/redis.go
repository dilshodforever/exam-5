package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type InMemoryStorageI interface {
	Set(key, value string, exp time.Duration) error
	Get(key string) (string, error)
	SaveEmailCode(email, code string, exp time.Duration)error
	SaveToken(id, token string, exp time.Duration) error
}

type storageRedis struct {
	client *redis.Client
}

// NewInMemoryStorage creates a new instance of Redis storage.
func NewInMemoryStorage(rdb *redis.Client) InMemoryStorageI {
	return &storageRedis{
		client: rdb,
	}
}

// Set stores a key-value pair with expiration in Redis.
func (r *storageRedis) Set(key, value string, exp time.Duration) error {
	fmt.Println("Setting key:", key, "with value:", value, "and expiration:", exp)
	err := r.client.Set(context.Background(), key, value, exp).Err()
	if err != nil {
		fmt.Println("Error in Set:", err)
		return err
	}
	return nil
}

// Get retrieves the value associated with a key from Redis.
func (r *storageRedis) Get(key string) (string, error) {
	val, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// SaveEmailCode stores an email verification code in Redis with expiration.
func (r *storageRedis) SaveEmailCode(email, code string, exp time.Duration) error {
	key := "email_code:" + email
	return r.Set(key, code, exp)
}
func (r *storageRedis) SaveToken(id, token string, exp time.Duration) error {
	key := "token:" + id
	fmt.Println("Saving token:", token, "with key:", key, "and expiration:", exp)
	err := r.Set(key, token, exp)
	if err != nil {
		fmt.Println("Error in SaveToken:", err)
		return err
	}
	fmt.Println("Token saved successfully")
	return nil
}
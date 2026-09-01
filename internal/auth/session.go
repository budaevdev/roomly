package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const sessionTTL = 24 * time.Hour

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
func CreateSession(ctx context.Context, rdb *redis.Client, userID int64) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	key := "session:" + token
	err = rdb.Set(ctx, key, userID, sessionTTL).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}

var ErrSessionNotFound = errors.New("session not found")

func GetSession(ctx context.Context, rdb *redis.Client, token string) (int64, error) {
	key := "session:" + token
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, ErrSessionNotFound
	}
	if err != nil {
		return 0, err
	}

	userID, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func DeleteSession(ctx context.Context, rdb *redis.Client, token string) error {
	key := "session:" + token
	return rdb.Del(ctx, key).Err()
}

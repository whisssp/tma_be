package service

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/goccy/go-json"
	redis2 "github.com/redis/go-redis/v9"
	"onboarding_test/internal/delivery/http/payload"
	"onboarding_test/pkg/redis"
	"time"
)

var redisDriver *redis2.Client
var ctx = context.Background()

func NewRedisService(redisClient *redis.RedisClient) {
	if redisClient != nil {
		redisDriver = redisClient.GetDriver()
	}
}

func GetHash(key string, property string, src interface{}) {
	val, err := redisDriver.HGet(ctx, key, property).Result()
	if err != nil {
		return
	}
	errU := json.Unmarshal([]byte(val), &src)
	if errU != nil {
		fmt.Println("error!!!!")
	}
}

func GetAllHashGeneric[M any](key string, src *map[string]M) {
	if redisDriver == nil {
		return
	}

	val, err := redisDriver.HGetAll(ctx, key).Result()
	if err != nil {
		return
	}

	*src = make(map[string]M)
	for k, v := range val {
		var item M
		if err := json.Unmarshal([]byte(v), &item); err != nil {
			fmt.Printf("Error unmarshaling value for key %s: %v\n", k, err)
			continue
		}
		(*src)[k] = item
	}

	return
}

func GetHashGenericWithPagination(key string, pagination *payload.Pagination) []string {
	if redisDriver == nil {
		return nil
	}
	start := int64(pagination.GetOffset())
	end := start + int64(pagination.GetLimit()) - 1
	val, err := redisDriver.LRange(ctx, key, start, end).Result()
	totalItems, _ := redisDriver.LLen(ctx, key).Result()
	pagination.TotalPages = int(totalItems)
	pagination.TotalPages = (int(totalItems) + pagination.GetLimit() - 1) / pagination.GetLimit()

	if err != nil {
		return nil
	}

	return val
}

func DeleteHash(key string, property string) error {
	if redisDriver == nil {
		return fmt.Errorf("RedisDriver not found")
	}
	_, err := redisDriver.HDel(ctx, key, property).Result()
	if err != nil {
		return err
	}
	return nil
}

func RedisGetHashGenericKey[M any](path string, key string, object *M) {
	if redisDriver == nil {
		return
	}
	GetHash(path, key, object)
}

func RedisSetHashGenericKeySlice[M any](path string, objects []M, getID func(M) int64, exp time.Duration) error {
	if redisDriver == nil {
		return fmt.Errorf("RedisDriver not found")
	}
	mapErrs := make(map[interface{}]string)
	for _, obj := range objects {
		key := getID(obj)
		marshalObject, _ := sonic.Marshal(obj)
		errR := SetHashObject(path, fmt.Sprintf("%v", key), marshalObject)
		if errR != nil {
			errSet := fmt.Errorf("RedisSetGenericHashSlice", "Error to save key:", key, "Object", string(marshalObject))
			mapErrs[key] = errSet.Error()
			continue
		}
	}

	if len(mapErrs) > 0 {
		return fmt.Errorf("error set: %v", mapErrs)
	}

	// Set expiration if needed (only once)
	if exp > 0 {
		if err := SetExpireKey(path, exp); err != nil {
			return fmt.Errorf("RedisSetGenericHashSlice", "Error setting expiration:", path, "err", err)
		}
	}
	return nil
}

func RedisSetHashGenericKey[M any](path string, key string, object M, exp time.Duration) error {
	if redisDriver == nil {
		return fmt.Errorf("RedisDriver not found")
	}
	marshalObject, _ := sonic.Marshal(object)
	errR := SetHashObject(path, key, marshalObject)
	if errR != nil {
		return fmt.Errorf("RedisSet Generic", "Error to save key:", key, "Object", string(marshalObject))
	}
	// Set expiration if needed
	if exp > 0 {
		if err := SetExpireKey(path, exp); err != nil {
			return fmt.Errorf("RedisSetGenericHash", "Error setting expiration:", key, "err", err.Error())
		}
	}
	return nil
}

func RedisRemoveHashGenericKey(path string, key string) error {
	if redisDriver == nil {
		return fmt.Errorf("RedisDriver not found")
	}

	if err := DeleteHash(path, key); err != nil {
		return fmt.Errorf("RedisRemoveGenericHash", "Error to delete key:", key, "err", err.Error())
	}
	return nil
}

//func RedisSetSetObjects(path string, key string)

func SetHashObject(key string, property string, object interface{}) error {

	return redisDriver.HSet(ctx, key, property, object).Err()
}

func SetExpireKey(key string, expiration time.Duration) error {
	return redisDriver.Expire(ctx, key, expiration).Err()
}
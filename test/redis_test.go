package test

import (
	"GeekCoding/models"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func TestRedisSet(t *testing.T) {
	err := rdb.Set(ctx, "name", "mmc", time.Second*10).Err()
	if err != nil {
		t.Fatal(err)
	}

}

func TestRedisGet(t *testing.T) {
	v, err := rdb.Get(ctx, "name").Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(v)
}

func TestRedisGetByModels(t *testing.T) {
	v, err := models.RDB.Get(ctx, "name").Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(v)
}

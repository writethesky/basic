package internal

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func init() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Config.Redis.Host, Config.Redis.Port),
		Username: Config.Redis.User,
		Password: Config.Redis.Password, // no password set
		DB:       Config.Redis.Database, // use default DB
	})
}

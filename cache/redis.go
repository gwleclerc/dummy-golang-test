package cache

import (
	"github.com/alicebob/miniredis"
)

var Redis *miniredis.Miniredis

func InitRedis() (err error) {
	Redis, err = miniredis.Run()
	return
}

func CloseRedis() {
	Redis.Close()
}

package cache

import (
	"github.com/alicebob/miniredis"
)

// Cache interface
type Cache interface {
	Get(string) (string, error)
	Set(string, string) error
	Close()
}

// Redis cache implementation
type Redis struct {
	mr *miniredis.Miniredis
}

// New cache
func New() (Cache, error) {
	redis, err := miniredis.Run()
	if err != nil {
		return nil, err
	}
	return &Redis{mr: redis}, nil
}

// Get value from redis
func (r Redis) Get(key string) (string, error) {
	return r.mr.Get(key)
}

// Set value into redis
func (r Redis) Set(key, value string) error {
	return r.mr.Set(key, value)
}

// Close redis instance
func (r Redis) Close() {
	r.mr.Close()
}

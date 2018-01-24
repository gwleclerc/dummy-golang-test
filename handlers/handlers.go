package handlers

import (
	"net/http"

	"github.com/gwleclerc/dummy-golang-test/cache"
	"github.com/labstack/echo"
)

// Cache handler
type Cache struct{}

// NewCacheHandler instance
func NewCacheHandler() Cache {
	return Cache{}
}

// Get value from cache
func (ch Cache) Get(c echo.Context) error {
	key := c.Param("key")
	value, err := cache.Redis.Get(key)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, value)
}

// Set value into cache
func (ch Cache) Set(c echo.Context) error {
	key := c.Param("key")
	value := c.FormValue("value")
	err := cache.Redis.Set(key, value)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "Value stored successfully")
}

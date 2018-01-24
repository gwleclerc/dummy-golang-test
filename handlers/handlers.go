package handlers

import (
	"net/http"

	"github.com/gwleclerc/dummy-golang-test/cache"
	"github.com/gwleclerc/dummy-golang-test/errs"
	"github.com/labstack/echo"
)

// CacheHandler structure
type CacheHandler struct{}

// NewCacheHandler instance
func NewCacheHandler() CacheHandler {
	return CacheHandler{}
}

// Get value from cache
func (ch CacheHandler) Get(c echo.Context) error {
	key := c.Param("key")
	store := c.Get("cache").(cache.Cache)
	value, err := store.Get(key)
	if err != nil {
		return c.String(http.StatusBadRequest, errs.Errorf("Can't get value for key '%v'", err, key))
	}

	return c.String(http.StatusOK, value)
}

// Set value into cache
func (ch CacheHandler) Set(c echo.Context) error {
	key := c.Param("key")
	value := c.FormValue("value")
	store := c.Get("cache").(cache.Cache)
	err := store.Set(key, value)
	if err != nil {
		return c.String(http.StatusBadRequest, errs.Errorf("Can't set value for key '%v'", err, key))
	}
	return c.String(http.StatusOK, "Value stored successfully")
}

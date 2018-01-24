package main

import (
	"log"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/gwleclerc/dummy-golang-test/cache"
	"github.com/gwleclerc/dummy-golang-test/handlers"
	"github.com/labstack/echo"
)

// ServerHeader middleware adds a `Server` header to the response.
func cacheMiddleware(store cache.Cache) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("cache", store)
			return next(c)
		}
	}
}

func main() {
	store, err := cache.New()
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	cacheHandler := handlers.NewCacheHandler()
	e := echo.New()

	redis := e.Group("/redis")
	{
		redis.Use(cacheMiddleware(store))
		redis.GET("/:key", cacheHandler.Get)
		redis.POST("/:key", cacheHandler.Set)
	}
	e.Server.Addr = ":8080"
	e.Logger.Fatal(gracehttp.Serve(e.Server))
}

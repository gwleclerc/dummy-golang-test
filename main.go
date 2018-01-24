package main

import (
	"github.com/gwleclerc/dummy-golang-test/cache"
	"github.com/gwleclerc/dummy-golang-test/handlers"
	"github.com/labstack/echo"
)

func main() {
	cache.InitRedis()
	defer cache.CloseRedis()
	cacheHandler := handlers.NewCacheHandler()
	e := echo.New()

	redis := e.Group("/redis")
	{
		redis.GET("/:key", cacheHandler.Get)
		redis.POST("/:key", cacheHandler.Set)
	}
	e.Logger.Fatal(e.Start(":8080"))
}

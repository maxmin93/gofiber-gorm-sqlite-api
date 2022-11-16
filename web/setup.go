package web

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupFiber func to setup Go Fiber
func SetupFiber(app *fiber.App) {
	setupMiddlewares(app)
	setupRoutes(app)
}

func setupMiddlewares(app *fiber.App) {
	// Logger middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))

	// CORS middleware
	app.Use(cors.New(cors.Config{
		// AllowOrigins: "https://gofiber.io, https://gofiber.net",
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Mount HTML view
	app.Mount("/home", SetupView())

	// Cache middleware
	// app.Use(cache.New())
	// ==> refresh 쿼리 파라미터까지 캐시가 먹어버려서 갱신이 안됨 (사용금지)
	testCache(app)
}

func testCache(app *fiber.App) {
	var requests = new(uint64)
	var requestsWithoutCache = new(uint64)

	cacheGroup := app.Group("/click")
	cacheGroup.Use(func(c *fiber.Ctx) error {
		atomic.AddUint64(requests, 1) // requests++
		return c.Next()
	})
	// Cache middleware
	cacheGroup.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},
		Expiration:   30 * time.Minute,
		CacheControl: true,
	}))
	cacheGroup.Get("/", func(c *fiber.Ctx) error {
		atomic.AddUint64(requestsWithoutCache, 1) // requestsWithoutCache++
		log.Printf("refreshed requests=%d", requestsWithoutCache)
		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/cacheHits", func(c *fiber.Ctx) error {
		// TODO: fix div by zero
		sum := atomic.LoadUint64(requests)
		cacheHits := (sum - atomic.LoadUint64(requestsWithoutCache))

		return c.JSON(fiber.Map{
			"requests":            sum,
			"cacheHits":           cacheHits,
			"cacheHitsPercentage": cacheHits * 100 / sum,
		})
	})
}

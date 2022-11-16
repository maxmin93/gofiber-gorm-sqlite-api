package web

import (
	"fmt"

	"example.com/db"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	// Create a new route group '/api'
	api := app.Group("/api")

	// Route Constraints
	// https://docs.gofiber.io/guide/routing#constraints
	api.Get("/dogs", db.GetDogs)
	api.Get("/dogs/:id<int>", db.GetDog)
	api.Post("/dogs", db.AddDog)
	api.Put("/dogs/:id<int>", db.UpdateDog)
	api.Patch("/dogs/:id<int>", db.UpdateDogPartial)
	api.Delete("/dogs/:id<int>", db.RemoveDog)

	setupRoutesOthers(app.Group("/test"))
}

func setupRoutesOthers(router fiber.Router) {
	// http://localhost:3000/test/greedy/?name=Bilbo
	// http://localhost:3000/test/greedy/?name=Bilbo&family=Baggins&city=Shire
	router.Get("/greedy/+", func(c *fiber.Ctx) error {
		return c.SendString(c.Params("+"))
	})

	// Limitations for characters in the path
	router.Get("/resource/key\\:value", func(c *fiber.Ctx) error {
		return c.SendString("escaped key:value")
	})

	// http://localhost:3000/test/hello/register
	router.Get("/hello/*", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("✋ %s", c.Params("*"))
		return c.SendString(msg) // => ✋ register
	})

	// http://localhost:3000/test/flights/LAX-SFO
	router.Get("/flights/:from-:to", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("💸 From: %s, To: %s", c.Params("from"), c.Params("to"))
		return c.SendString(msg) // => 💸 From: LAX, To: SFO
	})

	// http://localhost:3000/test/file/dictionary.txt
	router.Get("/file/:file.:ext", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("📃 %s.%s", c.Params("file"), c.Params("ext"))
		return c.SendString(msg) // => 📃 dictionary.txt
	})

	// http://localhost:3000/test/john/75/male
	router.Get("/:name/:age/:gender?", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("👴 %s is %s years old", c.Params("name"), c.Params("age"))
		if c.Params("gender") != "" {
			msg += fmt.Sprintf(" and is %s", c.Params("gender"))
		}
		return c.SendString(msg) // => 👴 john is 75 years old
	})

	// http://localhost:3000/test/john
	router.Get("/:name", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
		return c.SendString(msg) // => Hello john 👋!
	})
}

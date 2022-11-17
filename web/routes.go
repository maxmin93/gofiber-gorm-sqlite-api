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

	setupRoutesSpecial(app.Group("/path"))
	setupRoutesOthers(app.Group("/test"))
}

func setupRoutesSpecial(router fiber.Router) {
	// greedy/{Any} 가 아니면 404 Not Found 처리
	// 된다 http://localhost:3000/path/greedy/ABC ==> ABC
	// 된다 http://localhost:3000/path/greedy/ABC/DEF/G ==> ABC/DEF/G
	// 안됨 http://localhost:3000/path/greedy/
	// 안됨 http://localhost:3000/path/greedy/?name=Bilbo&family=Baggins&city=Shire
	router.Get("/greedy/+", func(c *fiber.Ctx) error {
		return c.SendString(c.Params("+"))
	})

	// Limitations for characters in the path
	// http://localhost:3000/path/resource/key:value
	router.Get("/resource/key\\:value", func(c *fiber.Ctx) error {
		return c.SendString("escaped key:value")
	})
}

func setupRoutesOthers(router fiber.Router) {
	// '*' 은 후속 Path Param 이 없어도 된다
	// http://localhost:3000/test/hello/Golang
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

	// {:gender} 는 optional 이다 (?표시)
	// http://localhost:3000/test/john/75
	// http://localhost:3000/test/john/75/male
	router.Get("/:name/:age/:gender?", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("👴 %s is %s years old", c.Params("name"), c.Params("age"))
		if c.Params("gender") != "" {
			msg += fmt.Sprintf(" and is %s", c.Params("gender"))
		}
		return c.SendString(msg) // => 👴 john is 75 years old
	})

	// 이 외의 모든 패턴을 처리한다
	// http://localhost:3000/test/john
	router.Get("/:name", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
		return c.SendString(msg) // => Hello john 👋!
	})
}

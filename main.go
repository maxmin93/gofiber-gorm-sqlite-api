package main

import (
	"fmt"
	"log"

	"example.com/db"
	m "example.com/models"
	u "example.com/utils"
	"example.com/web"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// NOTE: `대문자`로 시작하는 함수는 외부에서 접근 가능 (public)
	// 		 `소문자`로 시작하는 함수는 외부에서 접근 불가능 (private)
	u.Crypto()

	// 대소문자 구분
	book := m.Book{BookID: 1, BookName: "Go"}
	fmt.Print(book, " ")
	Book := new(m.Book)
	fmt.Print(book == *Book, "\n\n")

	//////////////////////////////////////////

	db.Connect()

	app := fiber.New()
	web.SetupFiber(app)

	var port = u.Config("PORT")
	log.Fatal(app.Listen(":" + port))
}

package models

// Book can be referenced from outside (public)
type Book struct {
	BookID   int    `json:"book_id"`
	BookName string `json:"book_name"`
}

type book struct { // private
	BookAuthor string `json:"book_author"`
	BookPrice  int    `json:"book_price"`
}

// HelloWorld can be referenced from outside (public)
func HelloWorld() string {
	return "Hello, world!"
}

func helloWorld() string { // private
	return "hello, world?"
}

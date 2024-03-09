package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Book struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	ReleaseYear string `json:"releaseyear"`
	Pages       int    `json:"pages"`
}

type user struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var books = []Book{
	{Id: 1, Title: "Harry Potter and the Philosopher's Stone", Author: "J. K. Rowling", ReleaseYear: "6/26/1997", Pages: 223},
	{Id: 2, Title: "Harry Potter and the Chamber of Secrets", Author: "J. K. Rowling", ReleaseYear: "7/2/1998", Pages: 251},
	{Id: 3, Title: "Harry Potter and the Prisoner of Azkaban", Author: "J. K. Rowling", ReleaseYear: "7/8/1999", Pages: 317},
	{Id: 4, Title: "Harry Potter and the Goblet of Fire", Author: "J. K. Rowling", ReleaseYear: "7/8/2000", Pages: 636},
	{Id: 5, Title: "Harry Potter and the Order of the Phoenix", Author: "J. K. Rowling", ReleaseYear: "6/21/2003", Pages: 766},
	{Id: 6, Title: "Harry Potter and the Half-Blood Prince", Author: "J. K. Rowling", ReleaseYear: "7/16/2005", Pages: 607},
	{Id: 7, Title: "Harry Potter and the Deathly Hallows", Author: "J. K. Rowling", ReleaseYear: "7/21/2007", Pages: 610},
}

var users = []user{
	{Id: 1, Name: "Fahri", Age: 23},
	{Id: 2, Name: "Desrizal", Age: 19},
}

func main() {
	router := gin.Default()
	router.Use(LoggerMiddleware)

	apiGroup := router.Group("/api")
	{
		booksGroup := apiGroup.Group("/books")
		{
			booksGroup.GET("/", getAllBook)
			booksGroup.POST("/", createBook)
			booksGroup.GET("/:id", getBookById)
			booksGroup.PUT("/:id", updateBookById)
		}

		usersGroup := router.Group("/users")
		{
			usersGroup.GET("/", getAllUser)
		}
	}

	/*
		// Middleware
		router.Use(customMiddleware)

		// Menampilkan Buku
		router.GET("/books", getAllBook)

		// Menambahkan Buku
		router.POST("/books", createBook)

		// Menampilkan detail buku berdasarkan Id
		router.GET("/books/:id", getBookById)

		// Mengupdate Buku
		router.PUT("/books/:id", updateBookById)

		// Menghapus Buku
		// router.DELETE("/books/:id", deleteBook)
	*/

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}

}

func getAllUser(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}

func getAllBook(c *gin.Context) {
	title := c.Query("title")

	if title == "" {
		c.JSON(http.StatusOK, books)
		return
	}

	var matchedBooks []Book
	for _, book := range books {
		if strings.Contains(strings.ToLower(book.Title), strings.ToLower(title)) {
			matchedBooks = append(matchedBooks, book)
		}
	}

	if len(matchedBooks) > 0 {
		c.JSON(http.StatusOK, matchedBooks)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book Not Found"})
	}

}

func createBook(c *gin.Context) {
	var newBook Book
	err := c.ShouldBind(&newBook)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	books = append(books, newBook)
	c.JSON(http.StatusCreated, newBook)
}

func getBookById(c *gin.Context) {
	id := c.Param("id")

	bookId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Book Id"})
		return
	}

	for _, book := range books {
		if book.Id == bookId {
			c.JSON(http.StatusOK, book)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Book Not Found"})
}

func updateBookById(c *gin.Context) {
	id := c.Param("id")

	bookId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Book Id"})
		return
	}

	var updatedBook Book

	if err := c.ShouldBind(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, book := range books {
		if book.Id == bookId {
			if strings.TrimSpace(updatedBook.Title) != "" {
				books[i].Title = updatedBook.Title
			}
			if strings.TrimSpace(updatedBook.Author) != "" {
				books[i].Author = updatedBook.Author
			}
			if strings.TrimSpace(updatedBook.ReleaseYear) != "" {
				books[i].ReleaseYear = updatedBook.ReleaseYear
			}
			if updatedBook.Pages != 0 {
				books[i].Pages = updatedBook.Pages
			}
			c.JSON(http.StatusOK, gin.H{"message": "Book Updated", "data": books[i]})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Book Not Found"})
}

func customMiddleware(c *gin.Context) {
	fmt.Println("Request Lewat CustomMiddleware")
	c.Next()
	fmt.Println("Response Lewat CustomMiddleware")
}

func LoggerMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()
	elapsed := time.Since(start).Microseconds()
	fmt.Println("Request memakan waktu sekitar", elapsed, "mikro detik")
}

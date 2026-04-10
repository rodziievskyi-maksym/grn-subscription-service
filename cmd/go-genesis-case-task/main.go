package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

// Навмисно неекспортована функція з коментарем, що вимагає лінтер revive
func doSomething() error {
	return errors.New("something went wrong")
}

func main() {
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "GitHub Release Notifier API (Gin) v0.0.1 on Go 1.26.2")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Gin server on port %s...", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

	// Помилка 1: Ігнорування помилки (trigger errcheck)
	doSomething()

	// Помилка 2: Марне присвоювання (trigger ineffassign)
	x := 10
	x = 20

	// Помилка 3: Unreachable code (trigger staticcheck)
	return
	fmt.Println("This will never run", x)
}

package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// Глобальные переменные и мьютекс
var (
	storedName string
	mu         sync.Mutex
)

func main() {
	r := gin.Default()

	// POST /set-name - записываем имя
	r.POST("/set-name", func(c *gin.Context) {
		// Заберём имя из формы или JSON (зависит от Content-Type)
		// Проще всего — из формы (PostForm).
		name := c.PostForm("name")
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "name is required",
			})
			return
		}

		// Запишем имя в глобальную переменную
		mu.Lock()
		storedName = name
		mu.Unlock()

		c.JSON(http.StatusOK, gin.H{
			"message": "Name set to " + name,
		})
	})

	// GET /hello - выводим приветствие
	r.GET("/hello", func(c *gin.Context) {
		mu.Lock()
		name := storedName
		mu.Unlock()

		if name == "" {
			c.JSON(http.StatusOK, gin.H{
				"message": "No name stored. Please set with /set-name first.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, " + name + "!",
		})
	})

	r.Run(":8080")
}

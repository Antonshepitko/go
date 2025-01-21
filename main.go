package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.GET("/health", func(c *gin.Context) {
        c.String(200, "OK")
    })
    r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
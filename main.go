package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("Please enter a valide openai api key!")
	}

	client := openai.NewClient(apiKey)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	r.GET("/gpt", func(c *gin.Context) {
		question := c.Query("q")
		if question == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Please provide a question via ?q=YourQuestion",
			})
			return
		}

		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT4o,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: question,
					},
				},
			},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		answer := resp.Choices[0].Message.Content
		c.JSON(http.StatusOK, gin.H{
			"question": question,
			"answer":   answer,
		})
	})

	r.Run(":8080")
}

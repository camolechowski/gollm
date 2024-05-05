package main

import (
	// "bufio"
	"context"
	"net/http"
	"os"

	// "strings"

	openai "gollm/openai"
	tui "gollm/tui"

	"github.com/charmbracelet/log"

	gin "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	oa "github.com/sashabaranov/go-openai"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file", "error", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("failed to load OPENAI_API_KEY from .env file")
	}

	// log.Debug("openai api key initialized", "apiKey", apiKey)
	// CLI MODE
	if len(os.Args) > 1 && os.Args[1] == "cli" {
		tui.RunTUI(apiKey)
		// reader := bufio.NewReader(os.Stdin)
		// for {
		// 	log.Printf("Enter a prompt: ")
		// 	input, err := reader.ReadString('\n')
		// 	if err != nil {
		// 		log.Fatal("failed to read input", "error", err)
		// 	}
		// 	input = strings.TrimSpace(input)

		// 	completion, err := openaiClient.CreateChatCompletion(context.Background(), oa.ChatCompletionRequest{
		// 		Model: oa.GPT4,
		// 		Messages: []oa.ChatCompletionMessage{
		// 			{
		// 				Role: oa.ChatMessageRoleUser,
		// 				Content: input,
		// 			},
		// 		},
		// 	})

		// 	if err != nil {
		// 		log.Error("failed to create completion.", "error", err)
		// 		continue
		// 	}
			
		// 	log.Debug("chat completion request context", "context", context.Background())
		// 	log.Debug("chat completion request input", "input", input)

		// 	log.Print(completion)
		// }
    
	} else {
		openaiClient := openai.NewClient(apiKey)

		router := gin.Default()

		router.POST("/api/completion", func(c *gin.Context) {
			var input struct {
				Text string `json:"text"`
			}
			if err := c.BindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
				log.Error("invalid request error", "error", err)
				return
			}

			completion, err := openaiClient.CreateChatCompletion(context.Background(), oa.ChatCompletionRequest{
				Model: oa.GPT4,
				Messages: []oa.ChatCompletionMessage{
					{
						Role: oa.ChatMessageRoleUser,
						Content: input.Text,
					},
				},
			})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "OpenAI API Error"})
				log.Error("openai api error", "error", err)
				return
			}

			c.JSON(http.StatusOK, gin.H{"completion": completion})
			// log.Debug("completion generated", "completion", completion)


		})

		router.Run(":8080")
		// log.Info("running http server on port 8080")
	}
}

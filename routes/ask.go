package routes

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

func StreamChat(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	prompt := ctx.GetHeader("Prompt")
	fileName := ctx.GetHeader("FileName")

	if token == "" {
		ctx.String(http.StatusUnauthorized, "You are not currently authorized to use this service!")
		return
	}
	if prompt == "" {
		ctx.String(http.StatusBadRequest, "Prompt is required")
		return
	}
	if fileName == "" {
		ctx.String(http.StatusBadRequest, "File is required")
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/plain")
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		ctx.Stream(func(w io.Writer) bool {
			c := openai.NewClient(token)
			ctx := context.Background()

			transcriptionReq := openai.AudioRequest{
				Model:    openai.Whisper1,
				FilePath: "temp/" + fileName,
			}
			resp, _ := c.CreateTranscription(context.Background(), transcriptionReq)

			req := openai.ChatCompletionRequest{
				Model:     openai.GPT3Dot5Turbo0301,
				MaxTokens: 500,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    "system",
						Content: "Answer the following question by using information from this transcribed audio file. If you do not know the answer with the information based on this audio file. Do not make up anything. If the answer if common knowledge, then go ahead and tell the user, however if it isn't, just simply say that you don't know the answer to the question. Use full complete, gramatically correct sentences for your responses!\n" + resp.Text,
					},
					{
						Role:    "user",
						Content: prompt,
					},
				},
				Temperature: 0.8,
				Stream: true,
			}
			stream, err := c.CreateChatCompletionStream(ctx, req)
			if err != nil {
				fmt.Fprintf(w, "ChatCompletionStream error: %v\n", err)
				return false
			}
			defer stream.Close()

			for {
				response, err := stream.Recv()
				if errors.Is(err, io.EOF) {
					fmt.Println("Stream finished")
					return false
				}

				if err != nil {
					fmt.Fprintf(w, "\nStream error: %v\n", err)
					return false
				}

				fmt.Fprint(w, response.Choices[0].Delta.Content)
				fmt.Println(response.Choices[0].Delta.Content)
				// Flush the response writer to ensure the content is sent immediately
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			}
		})
	}()

	wg.Wait()
}

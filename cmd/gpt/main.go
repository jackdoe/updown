package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	model := flag.String("m", "gpt-3.5-turbo", "options: gpt-4-32k, gpt-4, gpt-3.5-turbo")
	temp := flag.Float64("t", 0.7, "temperature")
	initial := flag.String("s", "summarize the following text", "initial system prompt, use @file.txt to load from a file named file.txt")
	flag.Parse()

	ai := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	systemPrompt := *initial
	if strings.HasPrefix(systemPrompt, "@") {
		b, err := ioutil.ReadFile(systemPrompt[1:])
		if err != nil {
			panic(err)
		}
		systemPrompt = string(b)
	}

	messages := []openai.ChatCompletionMessage{
		openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
	}

	resp, err := ai.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Temperature: float32(*temp),
			Model:       *model,
			Messages:    messages,
		},
	)
	if err != nil {
		panic(err)
	}

	text := resp.Choices[0].Message.Content
	fmt.Printf("%s\n", text)
}

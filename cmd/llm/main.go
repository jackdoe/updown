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
	fname := flag.String("f", "", "load the system prompt from a file")
	flag.Parse()

	ai := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	systemPrompt := ""
	if len(*fname) != 0 {
		b, err := ioutil.ReadFile(*fname)
		if err != nil {
			panic(err)
		}
		systemPrompt = string(b)
	}

	if len(systemPrompt) == 0 {
		systemPrompt = strings.Join(flag.Args(), " ")
	}

	messages := []openai.ChatCompletionMessage{
		openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
	}

	stdinBytes, _ := ioutil.ReadAll(os.Stdin)

	if len(stdinBytes) > 0 {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: string(stdinBytes),
		})
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
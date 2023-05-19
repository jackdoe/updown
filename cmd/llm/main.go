package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s [options] system promopt\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  Example:\n  cat file.txt | %s format as csv\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  Options:\n")
	flag.PrintDefaults()

}

func main() {
	flag.Usage = Usage

	model := flag.String("m", "gpt-3.5-turbo", "options: gpt-4-32k, gpt-4, gpt-3.5-turbo, etc")
	temp := flag.Float64("t", 0.7, "temperature")
	dostream := flag.Bool("s", true, "stream the output")
	fname := flag.String("f", "", "load the system prompt from a file")
	flag.Parse()
	key := os.Getenv("OPENAI_API_KEY")
	if len(key) == 0 {
		fmt.Fprintf(os.Stderr, "Define environment variable OPENAI_API_KEY, to get a key go to https://platform.openai.com/account/api-keys\n")
		os.Exit(1)
	}

	ai := openai.NewClient(key)
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

	req := openai.ChatCompletionRequest{
		Temperature: float32(*temp),
		Model:       *model,
		Messages:    messages,
	}
	if *dostream {
		req.Stream = true
		stream, err := ai.CreateChatCompletionStream(context.Background(), req)
		if err != nil {
			panic(err)
		}
		defer stream.Close()
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}

			if err != nil {
				panic(err)
			}

			os.Stdout.Write([]byte(response.Choices[0].Delta.Content))
		}
	} else {
		req.Stream = false
		resp, err := ai.CreateChatCompletion(context.Background(), req)

		if err != nil {
			panic(err)
		}
		text := resp.Choices[0].Message.Content
		os.Stdout.Write([]byte(text))
	}
	fmt.Printf("\n")
}

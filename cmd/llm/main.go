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
	"time"

	openai "github.com/sashabaranov/go-openai"
)

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s [options] system promopt\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  Example:\n  cat file.txt | %s format as csv\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  Options:\n")
	flag.PrintDefaults()

}

type promptFlags []string

func (i *promptFlags) String() string {
	return "array of promptFlags"
}

func (p *promptFlags) Set(value string) error {
	*p = append(*p, value)
	return nil
}

func shouldRetry(err error) bool {
	e := &openai.APIError{}
	if errors.As(err, &e) {
		switch e.HTTPStatusCode {
		case 401:
			return false
		// invalid auth or key (do not retry)
		case 429:
			return true
		// rate limiting or engine overload (wait and retry)
		case 500:
			return true
		default:
			return false
		}
	}
	return false
}

func main() {
	var prompts promptFlags
	flag.Usage = Usage

	model := flag.String("m", "gpt-3.5-turbo", "options: gpt-4-32k, gpt-4, gpt-3.5-turbo, etc")
	temp := flag.Float64("t", 0.7, "temperature")
	dostream := flag.Bool("s", true, "stream the output")
	debug := flag.Bool("d", false, "debug print request")
	readStdin := flag.Bool("stdin", true, "read the standard input")
	flag.Var(&prompts, "p", "set prompts, can be set multiple times, e.g -p @a.txt -p 'you are the best go developer' -p @b.txt")
	flag.Parse()
	key := os.Getenv("OPENAI_API_KEY")
	if len(key) == 0 {
		fmt.Fprintf(os.Stderr, "Define environment variable OPENAI_API_KEY, to get a key go to https://platform.openai.com/account/api-keys\n")
		os.Exit(1)
	}

	ai := openai.NewClient(key)
	messages := []openai.ChatCompletionMessage{}

	if len(flag.Args()) != 0 {
		messages = append(messages,
			openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: strings.Join(flag.Args(), " "),
			},
		)
	}

	for _, p := range prompts {
		if strings.HasPrefix(p, "@") {
			b, err := ioutil.ReadFile(p[1:])
			if err != nil {
				panic(err)
			}
			messages = append(messages,
				openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleSystem,
					Content: string(b),
				},
			)
		} else {
			messages = append(messages,
				openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleSystem,
					Content: p,
				},
			)

		}
	}

	if *readStdin {
		stdinBytes, _ := ioutil.ReadAll(os.Stdin)

		if len(stdinBytes) > 0 {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: string(stdinBytes),
			})
		}
	}

	if *debug {
		for _, m := range messages {
			os.Stdout.Write([]byte(m.Role + ":" + m.Content + "\n"))
		}
		os.Exit(0)
	}

	req := openai.ChatCompletionRequest{
		Temperature: float32(*temp),
		Model:       *model,
		Messages:    messages,
	}

	if *dostream {
		req.Stream = true
		for {
			stream, err := ai.CreateChatCompletionStream(context.Background(), req)
			if err != nil {
				if shouldRetry(err) {
					time.Sleep(5 * time.Second)
					continue
				}
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
			break
		}
	} else {
		req.Stream = false
		for {
			resp, err := ai.CreateChatCompletion(context.Background(), req)
			if err != nil {
				if shouldRetry(err) {
					time.Sleep(5 * time.Second)
					continue
				}
				panic(err)
			}
			text := resp.Choices[0].Message.Content
			os.Stdout.Write([]byte(text))
			break
		}
	}
	fmt.Printf("\n")
}

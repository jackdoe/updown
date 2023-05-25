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

func split(input string, chunkSize int) []string {
	if chunkSize == 0 {
		return []string{input}
	}

	lines := strings.Split(input, "\n")
	chunks := make([]string, 0)

	current := strings.Builder{}
	wordsInChunk := 0
	for _, line := range lines {
		wordsInLine := len(strings.Fields(line))

		if wordsInChunk < chunkSize {
			current.WriteString(line)
			current.WriteRune('\n')
			wordsInChunk += wordsInLine
		} else {
			chunks = append(chunks, current.String())
			current.Reset()

			current.WriteString(line)
			current.WriteRune('\n')

			wordsInChunk = wordsInLine
		}
	}
	if current.Len() > 0 {
		chunks = append(chunks, current.String())
	}

	return chunks
}

func main() {
	var prompts promptFlags
	flag.Usage = Usage

	model := flag.String("m", "gpt-3.5-turbo", "options: gpt-4-32k, gpt-4, gpt-3.5-turbo, etc")
	temp := flag.Float64("t", 0.7, "temperature")
	dostream := flag.Bool("s", true, "stream the output")
	expert := flag.String("e", "", "what are you an expert in?")
	debug := flag.Bool("d", false, "debug print request")
	autosplit := flag.Int("a", 0, "split the stdin input every N words (up to closest line break) and create separate request for each chunk")
	readStdin := flag.Bool("stdin", true, "read the standard input")
	flag.Var(&prompts, "p", "set prompts, can be set multiple times, e.g -p a.txt -p <('echo you are the best go developer') -p b.txt")
	flag.Parse()
	key := os.Getenv("OPENAI_API_KEY")
	if len(key) == 0 {
		fmt.Fprintf(os.Stderr, "Define environment variable OPENAI_API_KEY, to get a key go to https://platform.openai.com/account/api-keys\n")
		os.Exit(1)
	}

	ai := openai.NewClient(key)
	messagesSystemPrompt := []openai.ChatCompletionMessage{}

	if *expert != "" {
		messagesSystemPrompt = append(messagesSystemPrompt,
			openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: fmt.Sprintf("You are an expert in %s", *expert),
			},
		)
	}

	if len(flag.Args()) != 0 {
		messagesSystemPrompt = append(messagesSystemPrompt,
			openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: strings.Join(flag.Args(), " "),
			},
		)
	}

	for _, p := range prompts {
		b, err := ioutil.ReadFile(p)
		if err != nil {
			panic(err)
		}
		messagesSystemPrompt = append(messagesSystemPrompt,
			openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: string(b),
			},
		)
	}

	splitted := []openai.ChatCompletionRequest{}

	if *readStdin {
		stdinBytes, _ := ioutil.ReadAll(os.Stdin)
		stdinString := string(stdinBytes)

		for _, chunk := range split(stdinString, *autosplit) {
			messagesChunk := append([]openai.ChatCompletionMessage(nil), messagesSystemPrompt...)

			messagesChunk = append(messagesChunk, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: chunk,
			})

			splitted = append(splitted, openai.ChatCompletionRequest{
				Temperature: float32(*temp),
				Model:       *model,
				Messages:    messagesChunk,
			})
		}
	} else {
		splitted = append(splitted, openai.ChatCompletionRequest{
			Temperature: float32(*temp),
			Model:       *model,
			Messages:    messagesSystemPrompt,
		})
	}

	if *debug {
		for _, req := range splitted {
			os.Stdout.Write([]byte("----\n"))
			for _, m := range req.Messages {
				os.Stdout.Write([]byte(m.Role + ":" + m.Content + "\n"))
			}
		}
		os.Exit(0)
	}

	for _, req := range splitted {
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
				stream.Close()
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
}

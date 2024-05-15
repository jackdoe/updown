package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

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
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	defaultRoot := filepath.Join(home, ".updown-llm")

	pmodel := flag.String("m", "gpt-4o", "options: gpt-4-32k, gpt-4, gpt-3.5-turbo, gpt-3.5-turbo-16k, 3 means gpt-3.5-turbo")
	temp := flag.Float64("temperature", 0.7, "temperature")
	lastNfiles := flag.Uint64("n", 40, "only read last N question/answers")
	root := flag.String("root", defaultRoot, "root")
	pmultiline := flag.Bool("multiline", false, "multiline")
	topic := flag.String("t", "programming", "topic")
	debug := flag.Bool("d", false, "debug print request")
	onlyInput := flag.Bool("i", false, "only stdin")
	flag.Parse()
	hasCmdLinePrompt := len(flag.Args()) != 0 || *onlyInput
	key := os.Getenv("OPENAI_API_KEY")
	if len(key) == 0 {
		fmt.Fprintf(os.Stderr, "Define environment variable OPENAI_API_KEY, to get a key go to https://platform.openai.com/account/api-keys\n")
		os.Exit(1)
	}
	multiline := *pmultiline
	model := *pmodel
	if model == "3.5" || model == "3" {
		model = "gpt-3.5-turbo"
	} else if model == "4" {
		model = "gpt-4o"
	}
	scanner := bufio.NewScanner(os.Stdin)
	ai := openai.NewClient(key)
	dir := filepath.Join(*root, *topic)
	os.MkdirAll(dir, 0700)
	prepromptFn := filepath.Join(dir, "preprompt")
	preprompt := ""
	if !hasCmdLinePrompt {
		if _, err := os.Stat(prepromptFn); err != nil {
			fmt.Printf("preprompt file not found:\n%s\nwrite the preprompt (empty for default).\n> ", prepromptFn)
			for scanner.Scan() {
				preprompt = strings.TrimSpace(scanner.Text())
				break
			}
			if preprompt == "" {
				preprompt = `I want you to act as a very experienced and versatile developer with experience in C, Go, the linux kernel, netbsd. You also have experience with vim and tmux.`
			}

			if err := ioutil.WriteFile(prepromptFn, []byte(preprompt), 0600); err != nil {
				panic(err)
			}
		} else {
			b, err := ioutil.ReadFile(prepromptFn)
			if err != nil {
				panic(err)
			}

			preprompt = string(b)
		}
	}

	messagesSystemPrompt := []openai.ChatCompletionMessage{}
	if preprompt != "" {
		messagesSystemPrompt = append(messagesSystemPrompt,
			openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: fmt.Sprintf(preprompt),
			},
		)
	}

	// load the old question/answer files
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	if len(files) > int(*lastNfiles) {
		files = files[len(files)-int(*lastNfiles):]
	}
	id := int64(0)
	for _, p := range files {
		if p.IsDir() {
			continue
		}

		splitted := strings.Split(p.Name(), ".")
		if len(splitted) != 2 {
			continue
		}
		suffix := splitted[1]
		id, err = strconv.ParseInt(splitted[0], 10, 64)
		if err != nil {
			panic(err)
		}

		b, err := ioutil.ReadFile(filepath.Join(dir, p.Name()))
		if err != nil {
			panic(err)
		}
		role := openai.ChatMessageRoleUser
		if suffix == "1a" {
			role = openai.ChatMessageRoleAssistant
			if !hasCmdLinePrompt {
				os.Stdout.WriteString("A:")
			}
		} else {
			if !hasCmdLinePrompt {
				os.Stdout.WriteString("Q:")
			}
		}
		if !hasCmdLinePrompt {
			os.Stdout.Write(b)
		}
		messagesSystemPrompt = append(messagesSystemPrompt,
			openai.ChatCompletionMessage{
				Role:    role,
				Content: string(b),
			},
		)
		if !hasCmdLinePrompt {
			os.Stdout.WriteString("\n")
		}

	}
	id++

	question := strings.Builder{}
	answer := strings.Builder{}

	inputPrompt := func() {
		if multiline {
			os.Stdout.WriteString(fmt.Sprintf("multiline %s> ", model))
		} else {
			os.Stdout.WriteString(fmt.Sprintf("%s> ", model))
		}

	}
	qa := func() {
		if question.Len() == 0 {
			return
		}
		questionFile, err := os.OpenFile(filepath.Join(dir, fmt.Sprintf("%06d.0q", id)), os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
		questionFile.WriteString(question.String())
		questionFile.Close()

		messagesSystemPrompt = append(messagesSystemPrompt,
			openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: question.String(),
			},
		)

		answerFile, err := os.OpenFile(filepath.Join(dir, fmt.Sprintf("%06d.1a", id)), os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}

		req := openai.ChatCompletionRequest{
			Stream:      true,
			Temperature: float32(*temp),
			Model:       model,
			Messages:    messagesSystemPrompt,
		}

		if *debug {
			for _, m := range req.Messages {
				os.Stdout.Write([]byte(m.Role + ":" + m.Content + "\n"))
			}
		}

		stream, err := ai.CreateChatCompletionStream(context.Background(), req)
		if err != nil {
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
			delta := []byte(response.Choices[0].Delta.Content)
			answer.Write(delta)

			os.Stdout.Write(delta)
			answerFile.Write(delta)
		}

		messagesSystemPrompt = append(messagesSystemPrompt,
			openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: answer.String(),
			},
		)

		stream.Close()
		answerFile.Close()

		answer.Reset()
		question.Reset()

		os.Stdout.WriteString("\n")
		id++

	}

	if hasCmdLinePrompt {
		if len(flag.Args()) > 0 {
			question.WriteString(strings.Join(flag.Args(), " "))
			question.WriteRune('\n')
		}
		if *onlyInput {
			b, _ := ioutil.ReadAll(os.Stdin)
			question.Write(b)
		}
		qa()
		os.Exit(0)
	}

	inputPrompt()
	for scanner.Scan() {
		text := scanner.Text()
		if text == ":3" {
			model = "gpt-3.5-turbo-16k"
			inputPrompt()
			continue
		}
		if text == ":4" {
			model = "gpt-4"
			inputPrompt()
			continue
		}
		if text == ":." {
			multiline = !multiline
			inputPrompt()
			continue
		}

		if multiline && text == "." {
			qa()
			inputPrompt()
		} else {
			if multiline {
				question.WriteString(text)
				question.WriteString("\n")
			} else {
				question.WriteString(text)
				qa()
				inputPrompt()
			}
		}
	}
}

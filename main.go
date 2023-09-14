package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		name          string
		apiKey        string
		content       string
		invalidName   bool = true
		invalidApiKey bool = true
	)

	//set api key
	for invalidApiKey {
		fmt.Println(`
			Openai api key is required (ex: sk-eM0aaaWkgUIrmRwlUJLToBT3BlbkFJysHpaj8e4x36Qux8), 
			if you have no idea about api key, you must open your openai account first and generate an api key on this page https://platform.openai.com/account/api-keys
		`)
		apiKey = os.Getenv("OPEN_AI_API_KEY")

		if apiKey == "" {
			apiKey = StringPrompt("Supply your API key, please: ")
		}

		responseCode := ApiKeyValidation(apiKey)
		if responseCode == 200 {
			invalidApiKey = false
		}
	}

	// set name
	for invalidName {

		name = StringPrompt("Hi, what is your name?")

		if len(name) >= 3 {
			invalidName = false
		}

	}

	client := openai.NewClient(apiKey)

	for {

		content = StringPrompt(name + " : ")

		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: content,
					},
				},
			},
		)

		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			return
		}

		fmt.Println(resp.Choices[0].Message.Content)
	}

}

func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func ApiKeyValidation(apiKey string) int {
	apiUrl := "https://api.openai.com/v1/engines"

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Openai API key is invalid!")
	}

	return resp.StatusCode
}

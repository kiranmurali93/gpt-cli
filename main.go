package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "chatgpt-cli",
	Short: "Interact with ChatGPT via CLI",
	Run:   startChat,
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is ./config.yaml)")
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile("configFile")
	} else {
		home, err := os.Getwd()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Cant read config file", err)
	}
}

func startChat(cmd *cobra.Command, args []string) {
	apiKey := viper.GetString("api_key")

	if apiKey == "" {
		log.Fatal("API key not found in the config file")
	}

	fmt.Println("ChatGPT CLI - Type 'exit' to quit")
	for {
		fmt.Println("You: ")
		userInput, err := getUserInput()

		if err != nil {
			fmt.Println(err)
		}

		if strings.ToLower(userInput) == "exit\n" {
			fmt.Println("Ending chat....")
			os.Exit(0)
		}

		res, err := getResFromGpt(apiKey, userInput)

		if err != nil {
			fmt.Println(err)
			break
			os.Exit(1)
		}
    fmt.Println("assistant: ", res)

	}
}

func getUserInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')

	if err != nil {
		return "", fmt.Errorf("error in getting input: %s", err)
	}

	return line, nil
}

func getResFromGpt(apiKey string, input string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	payload := map[string]interface{}{
		"model":       "gpt-3.5-turbo",
		"messages":    []map[string]string{{"role": "system", "content": "You are a helpful assistant."}, {"role": "user", "content": input}},
		"temperature": 0.7,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON payload: %s", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return "", fmt.Errorf("error making HTTP request: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %s", err)
	}

	var completion ChatCompletion

	resErr := json.Unmarshal([]byte(body), &completion)
	if resErr != nil {
		return "", fmt.Errorf("Error unmarshaling JSON:", resErr)
	}

	messageContent := completion.Choices[0].Message.Content

	return string(messageContent), nil
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Index        int         `json:"index"`
	Message      Message     `json:"message"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatCompletion struct {
	ID                string      `json:"id"`
	Object            string      `json:"object"`
	Created           int64       `json:"created"`
	Model             string      `json:"model"`
	Choices           []Choice    `json:"choices"`
	Usage             Usage       `json:"usage"`
	SystemFingerprint interface{} `json:"system_fingerprint"`
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

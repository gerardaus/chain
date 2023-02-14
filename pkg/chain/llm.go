package chain

import (
	"context"
	"fmt"
	"math"
	"os"
	"time"

	gogpt "github.com/gerardaus/go-gpt3"
)

// llmresult

type LLMResult struct {
	Text       string
	TokensUsed int
}

type Generation struct {
	Text string
}

// llm

type LLM interface {
	Generate(ctx context.Context, prompt string) (*LLMResult, error)
}

type openAILLM struct {
	Client    *gogpt.Client
	MaxTokens int
}

func NewOpenAI(maxTokens int) LLM {
	return &openAILLM{
		Client:    gogpt.NewClient(os.Getenv("OPENAI_API_KEY")),
		MaxTokens: maxTokens,
	}
}
func (o *openAILLM) Generate(ctx context.Context, prompt string) (*LLMResult, error) {
	resp, err := Retry(
		func() (gogpt.CompletionResponse, error) {
			return o.Client.CreateCompletion(ctx,
				gogpt.CompletionRequest{
					Model:            "text-davinci-003",
					Prompt:           prompt,
					Temperature:      0.7,
					MaxTokens:        o.MaxTokens,
					TopP:             1,
					BestOf:           1,
					FrequencyPenalty: 0.5,
					PresencePenalty:  0,
					Echo:             false,
				},
			)
		},
		3,
	)

	if err != nil {
		return nil, err
	}

	return &LLMResult{
		Text:       resp.Choices[0].Text,
		TokensUsed: resp.Usage.TotalTokens,
	}, nil
}

func Retry(fn func() (gogpt.CompletionResponse, error), retries int) (gogpt.CompletionResponse, error) {
	var err error
	for i := 0; i < retries; i++ {
		resp, err := fn()
		if err == nil {
			return resp, nil
		}
		time.Sleep(time.Duration(math.Pow(3, float64(i))) * time.Second)
	}
	return gogpt.CompletionResponse{}, fmt.Errorf("failed after %v attempts with err=%v", retries, err.Error())
}

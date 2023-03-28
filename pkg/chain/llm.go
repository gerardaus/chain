package chain

import (
	"context"
	"fmt"
	"math"
	"os"
	"strings"
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
	Config    *OpenAIConfig
	MaxTokens int
}

type OpenAIConfig struct {
	Model            string
	MaxTokens        int
	Temperature      float32
	TopP             float32
	N                int
	Stream           bool
	BestOf           int
	LogProbs         int
	Echo             bool
	Stop             []string
	PresencePenalty  float32
	FrequencyPenalty float32
}

func NewOpenAI(config *OpenAIConfig) LLM {
	return &openAILLM{
		Client: gogpt.NewClient(os.Getenv("OPENAI_API_KEY")),
		Config: config,
	}
}

func (o *openAILLM) Generate(ctx context.Context, prompt string) (*LLMResult, error) {
	resp, err := Retry(
		func() (gogpt.CompletionResponse, error) {
			return o.Client.CreateCompletion(ctx,
				gogpt.CompletionRequest{
					Model:            o.Config.Model,
					Prompt:           prompt,
					Temperature:      o.Config.Temperature,
					MaxTokens:        o.Config.MaxTokens,
					TopP:             o.Config.TopP,
					BestOf:           o.Config.BestOf,
					FrequencyPenalty: o.Config.FrequencyPenalty,
					PresencePenalty:  o.Config.PresencePenalty,
					Echo:             o.Config.Echo,
				},
			)
		},
		3,
	)

	if err != nil {
		return nil, err
	}

	return &LLMResult{
		Text:       strings.Trim(resp.Choices[0].Text, "\n"),
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
	if err != nil {
		return gogpt.CompletionResponse{}, fmt.Errorf("failed after %v attempts with err=%v", retries, err.Error())
	}
	return gogpt.CompletionResponse{}, fmt.Errorf("failed after %v attempts with err=%v", retries, "timed out")
}

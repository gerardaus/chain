package chain

import (
	"context"
	"fmt"

	"golang.org/x/exp/maps"
)

// llm chain
type LLMChain interface {
	Predict(map[string]string) (string, error)
}

type LLMChainConfig struct {
	LLM            LLM
	PromptTemplate *PromptTemplate
	Verbose        bool
	Memory         Memory
}

type basicLLMChain struct {
	PromptTemplate *PromptTemplate
	LLM            LLM
	Verbose        bool
	Memory         Memory
}

func NewLLMChain(config *LLMChainConfig) LLMChain {
	return &basicLLMChain{
		PromptTemplate: config.PromptTemplate,
		LLM:            config.LLM,
		Verbose:        config.Verbose,
		Memory:         config.Memory,
	}
}

func (b *basicLLMChain) Predict(args map[string]string) (string, error) {
	if b.Memory != nil {
		maps.Copy(args, b.Memory.LoadMemoryVariables())
	}

	prompt, err := b.PromptTemplate.Format(args)
	if err != nil {
		return "", err
	}

	if b.Verbose {
		fmt.Println(prompt)
	}

	result, err := b.LLM.Generate(context.Background(), prompt)
	if err != nil {
		return "", err
	}

	if b.Memory != nil {
		b.Memory.SaveContext(
			args["input"],
			result.Text,
		)
	}
	return result.Text, nil
}

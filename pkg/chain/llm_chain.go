package chain

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/exp/maps"
)

// llm chain
type LLMChain interface {
	Predict(map[string]string) (string, error)
	Generate(map[string]string) (*LLMResult, error)
	GetMemory() Memory
}

type LLMChainConfig struct {
	LLM                LLM
	PromptTemplate     *PromptTemplate
	Verbose            bool
	Memory             Memory
	IdleExpiryDuration time.Duration
}

type basicLLMChain struct {
	PromptTemplate     *PromptTemplate
	LLM                LLM
	Verbose            bool
	Memory             Memory
	IdleExpiryDuration time.Duration
}

func NewLLMChain(config *LLMChainConfig) LLMChain {
	return &basicLLMChain{
		PromptTemplate:     config.PromptTemplate,
		LLM:                config.LLM,
		Verbose:            config.Verbose,
		Memory:             config.Memory,
		IdleExpiryDuration: config.IdleExpiryDuration,
	}
}

func (b *basicLLMChain) Predict(args map[string]string) (string, error) {
	if result, err := b.Generate(args); err != nil {
		return "", err
	} else {
		return result.Text, nil
	}
}

func (b *basicLLMChain) Generate(args map[string]string) (*LLMResult, error) {
	if b.Memory != nil {
		maps.Copy(args, b.Memory.LoadMemoryVariables())
	}

	prompt, err := b.PromptTemplate.Format(args)
	if err != nil {
		return nil, err
	}

	if b.Verbose {
		fmt.Println(prompt)
	}

	result, err := b.LLM.Generate(context.Background(), prompt)
	if err != nil {
		return nil, err
	}

	if b.Memory != nil {
		// XXX TODO fix this input here
		b.Memory.SaveContext(args["input"], result.Text)
	}
	return result, nil
}

func (b *basicLLMChain) GetMemory() Memory {
	return b.Memory
}

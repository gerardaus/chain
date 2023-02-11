package chain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DummyLLM struct {
	Output string
}

func (d *DummyLLM) Generate(ctx context.Context, prompt string) (*LLMResult, error) {
	return &LLMResult{
		Text:       d.Output,
		TokensUsed: 5,
	}, nil
}

func TestLLMChain(t *testing.T) {
	prompt := &PromptTemplate{
		Template: "Tell me a {{.Adjective}} joke",
		Inputs:   []string{"Adjective"},
	}

	chain := NewLLMChain(
		&LLMChainConfig{
			LLM: &DummyLLM{
				Output: "Why did the chicken cross the road? To get to the other side",
			},
			PromptTemplate: prompt,
			Verbose:        true,
			//Memory:         NewConversationalBufferWindowMemory(),
		},
	)

	input := map[string]string{
		"Adjective": "funny",
	}

	output, err := chain.Predict(input)
	assert.Nil(t, err)
	assert.Equal(t, "Why did the chicken cross the road? To get to the other side", output)
	assert.Nil(t, nil)
}

package main

import (
	"bufio"
	"fmt"
	"os"

	"magneticlabs.com/link/pkg/chain"
)

func main() {
	prompt := &chain.PromptTemplate{
		Inputs: []string{"input", "history"},
		Template: `
{{.history}}

Assume the persona of an AI robot. Your task is to answer any questions the human would like answered.

Human: {{.input}}
AI:
`}

	llmchain := chain.NewLLMChain(
		&chain.LLMChainConfig{
			LLM:            chain.NewOpenAI(),
			PromptTemplate: prompt,
			Verbose:        false,
			Memory: chain.NewConversationBufferWindowMemory(
				&chain.BufferConfig{
					HumanPrefix: "Human:",
					AIPrefix:    "AI:",
					K:           20,
				},
			),
		},
	)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()
		if input == "quit" {
			break
		}

		output, err := llmchain.Predict(map[string]string{"input": input})
		if err != nil {
			fmt.Println("there was an error, try again")
			continue
		}
		fmt.Println(output)
	}
}

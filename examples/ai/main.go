package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gerardaus/chain/pkg/chain"
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
			LLM: chain.NewOpenAI(
				&chain.OpenAIConfig{
					Model:            "text-davinci-003",
					Temperature:      0.9,
					MaxTokens:        140,
					TopP:             1,
					BestOf:           1,
					FrequencyPenalty: 0,
					PresencePenalty:  0,
					Echo:             false,
					Stop:             []string{"\n"},
				},
			),
			PromptTemplate: prompt,
			Verbose:        false,
			Memory: chain.NewConversationBufferWindowMemory(
				&chain.ConversationBufferWindowMemoryConfig{
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

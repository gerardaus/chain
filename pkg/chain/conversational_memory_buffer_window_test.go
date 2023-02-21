package chain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConversationBufferWindowMemory(t *testing.T) {
	tests := []struct {
		name     string
		contexts [][]string
		history  string
		k        int
	}{
		{
			name:     "test 1",
			contexts: [][]string{{"in1", "out1"}, {"in2", "out2"}, {"in3", "out3"}},
			history:  "Human: in3\nAI: out3",
			k:        1,
		},
		{
			name:     "test 1",
			contexts: [][]string{{"in1", "out1"}, {"in2", "out2"}, {"in3", "out3"}},
			history:  "Human: in2\nAI: out2\nHuman: in3\nAI: out3",
			k:        2,
		},
		{
			name:     "test 2",
			contexts: [][]string{{"in1", "out1"}, {"in2", "out2"}, {"in3", "out3"}},
			history:  "Human: in1\nAI: out1\nHuman: in2\nAI: out2\nHuman: in3\nAI: out3",
			k:        3,
		},
	}

	for _, tt := range tests {
		memory := NewConversationBufferWindowMemory(
			&ConversationBufferWindowMemoryConfig{
				HumanPrefix: "Human:",
				AIPrefix:    "AI:",
				K:           tt.k,
			},
		)
		for _, context := range tt.contexts {
			memory.SaveContext(context[0], context[1])
		}

		assert.Equal(t,
			map[string]string{"history": tt.history},
			memory.LoadMemoryVariables(),
		)
	}
}

func TestConversationBufferWindowMemory_Clear(t *testing.T) {
	memory := NewConversationBufferWindowMemory(
		&ConversationBufferWindowMemoryConfig{
			HumanPrefix: "Human:",
			AIPrefix:    "AI:",
			K:           2,
		},
	)
	memory.SaveContext("hey how are you", "im good thanks!")
	memory.Clear()
	vars := memory.LoadMemoryVariables()
	assert.Equal(t, map[string]string{"history": ""}, vars)
}

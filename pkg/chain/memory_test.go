package chain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConversationBufferWindowMemory(t *testing.T) {
	memory := NewConversationBufferWindowMemory(
		&BufferConfig{
			HumanPrefix: "Human:",
			AIPrefix:    "AI:",
			K:           2,
		},
	)
	memory.SaveContext(
		"hey how are you",
		"im good thanks!",
	)
	vars := memory.LoadMemoryVariables()
	assert.Equal(t,
		map[string]string{
			"history": "Human: hey how are you\nAI: im good thanks!",
		},
		vars,
	)
}

func TestConversationBufferWindowMemory_k_limit(t *testing.T) {
	memory := NewConversationBufferWindowMemory(
		&BufferConfig{
			HumanPrefix: "Human:",
			AIPrefix:    "AI:",
			K:           2,
		},
	)

	memory.SaveContext(
		"hey how are you",
		"im good thanks!",
	)
	memory.SaveContext(
		"hello",
		"hi",
	)
	memory.SaveContext(
		"bye",
		"bye now!",
	)

	vars := memory.LoadMemoryVariables()
	assert.Equal(t,
		map[string]string{
			"history": "Human: hello\nAI: hi\nHuman: bye\nAI: bye now!",
		},
		vars,
	)
}

func TestConversationBufferWindowMemory_Clear(t *testing.T) {
	memory := NewConversationBufferWindowMemory(
		&BufferConfig{
			HumanPrefix: "Human:",
			AIPrefix:    "AI:",
			K:           2,
		},
	)
	memory.SaveContext(
		"hey how are you",
		"im good thanks!",
	)

	memory.Clear()
	vars := memory.LoadMemoryVariables()
	assert.Equal(t, map[string]string{"history": ""}, vars)
}

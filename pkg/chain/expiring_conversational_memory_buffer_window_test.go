package chain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExpiringConversationBufferWindowMemory(t *testing.T) {
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
		memory := NewExpiringConversationBufferWindowMemory(
			&ExpiringConversationBufferConfig{
				HumanPrefix:        "Human:",
				AIPrefix:           "AI:",
				K:                  tt.k,
				IdleExpiryDuration: time.Second * 10,
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

func TestExpiringConversationBufferWindowMemory_Clear(t *testing.T) {
	memory := NewExpiringConversationBufferWindowMemory(
		&ExpiringConversationBufferConfig{
			HumanPrefix:        "Human:",
			AIPrefix:           "AI:",
			K:                  2,
			IdleExpiryDuration: time.Second * 10,
		},
	)
	memory.SaveContext("hey how are you", "im good thanks!")
	memory.Clear()
	vars := memory.LoadMemoryVariables()
	assert.Equal(t, map[string]string{"history": ""}, vars)
}

func TestExpiringConversationBufferWindowMemory_IdleExpiryDuration(t *testing.T) {
	tests := []struct {
		name     string
		expiry   time.Duration
		contexts [][]string
		sleep    time.Duration
		history  map[string]string
	}{
		{
			name:     "expires memory after 50 milliseconds",
			expiry:   time.Millisecond * 50,
			contexts: [][]string{{"in1", "out1"}},
			sleep:    time.Millisecond * 100,
			history:  map[string]string{"history": ""},
		},
		{
			name:     "does not expire memory after 150 milliseconds",
			expiry:   time.Millisecond * 150,
			contexts: [][]string{{"in1", "out1"}},
			sleep:    time.Millisecond * 100,
			history:  map[string]string{"history": "Human: in1\nAI: out1"},
		},
	}

	for _, tt := range tests {
		memory := NewExpiringConversationBufferWindowMemory(
			&ExpiringConversationBufferConfig{
				HumanPrefix:        "Human:",
				AIPrefix:           "AI:",
				K:                  2,
				IdleExpiryDuration: tt.expiry,
			},
		)
		for _, context := range tt.contexts {
			memory.SaveContext(context[0], context[1])
		}
		time.Sleep(tt.sleep)
		assert.Equal(t, tt.history, memory.LoadMemoryVariables())
	}
}

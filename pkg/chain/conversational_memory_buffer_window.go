package chain

import (
	"strings"
	"sync"
)

type memBuffer struct {
	sync.Mutex
	Buffer      []string
	HumanPrefix string
	AIPrefix    string
	MemoryKey   string
	K           int
}

type ConversationBufferWindowMemoryConfig struct {
	K           int
	HumanPrefix string
	AIPrefix    string
}

func NewConversationBufferWindowMemory(config *ConversationBufferWindowMemoryConfig) Memory {
	return &memBuffer{
		Buffer:      []string{},
		HumanPrefix: config.HumanPrefix,
		AIPrefix:    config.AIPrefix,
		MemoryKey:   "history",
		K:           config.K,
	}
}

func (m *memBuffer) Clear() {
	m.Lock()
	defer m.Unlock()
	m.Buffer = []string{}
}

func (m *memBuffer) SaveContext(input, output string) {
	m.Lock()
	defer m.Unlock()

	m.Buffer = append(m.Buffer,
		strings.Join(
			[]string{m.HumanPrefix + " " + input, m.AIPrefix + " " + output}, "\n"),
	)
}

func (m *memBuffer) LoadMemoryVariables() map[string]string {
	m.Lock()
	defer m.Unlock()

	history := m.Buffer
	if m.K < len(m.Buffer) {
		history = m.Buffer[len(m.Buffer)-m.K:]
	}
	return map[string]string{
		m.MemoryKey: strings.Join(history, "\n"),
	}
}

package chain

import "strings"

type Memory interface {
	Clear()
	SaveContext(input string, output string)
	LoadMemoryVariables() map[string]string
}
type memBuffer struct {
	Buffer      []string
	HumanPrefix string
	AIPrefix    string
	MemoryKey   string
	K           int
}

type BufferConfig struct {
	K           int
	HumanPrefix string
	AIPrefix    string
}

func NewConversationBufferWindowMemory(config *BufferConfig) Memory {
	return &memBuffer{
		Buffer:      []string{},
		HumanPrefix: config.HumanPrefix,
		AIPrefix:    config.AIPrefix,
		MemoryKey:   "history",
		K:           config.K,
	}
}

func (m *memBuffer) Clear() {
	m.Buffer = []string{}
}

func (m *memBuffer) SaveContext(input, output string) {
	m.Buffer = append(m.Buffer,
		strings.Join(
			[]string{m.HumanPrefix + " " + input, m.AIPrefix + " " + output}, "\n"),
	)
}

func (m *memBuffer) LoadMemoryVariables() map[string]string {
	history := m.Buffer
	if m.K < len(m.Buffer) {
		history = m.Buffer[len(m.Buffer)-m.K:]
	}
	return map[string]string{
		m.MemoryKey: strings.Join(history, "\n"),
	}
}

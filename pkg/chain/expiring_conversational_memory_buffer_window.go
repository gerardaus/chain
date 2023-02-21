package chain

import (
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type expiringMemBuffer struct {
	Buffer             []string
	HumanPrefix        string
	AIPrefix           string
	MemoryKey          string
	K                  int
	IdleExpiryDuration time.Duration
	Ticker             *time.Ticker
}

type ExpiringConversationBufferConfig struct {
	K                  int
	HumanPrefix        string
	AIPrefix           string
	IdleExpiryDuration time.Duration
}

func NewExpiringConversationBufferWindowMemory(config *ExpiringConversationBufferConfig) Memory {
	buffer := &expiringMemBuffer{
		Buffer:             []string{},
		HumanPrefix:        config.HumanPrefix,
		AIPrefix:           config.AIPrefix,
		MemoryKey:          "history",
		K:                  config.K,
		IdleExpiryDuration: config.IdleExpiryDuration,
		Ticker:             time.NewTicker(config.IdleExpiryDuration),
	}

	go func(buffer *expiringMemBuffer) {
		for {
			select {
			case <-buffer.Ticker.C:
				log.Debugln("clearing buffer")
				buffer.Clear()
			}
		}
	}(buffer)

	return buffer
}

func (m *expiringMemBuffer) ResetTicker() {
	m.Ticker.Reset(m.IdleExpiryDuration)
}

func (m *expiringMemBuffer) Clear() {
	m.ResetTicker()
	m.Buffer = []string{}
}

func (m *expiringMemBuffer) SaveContext(input, output string) {
	m.ResetTicker()
	m.Buffer = append(m.Buffer,
		strings.Join(
			[]string{m.HumanPrefix + " " + input, m.AIPrefix + " " + output}, "\n"),
	)
}

func (m *expiringMemBuffer) LoadMemoryVariables() map[string]string {
	history := m.Buffer
	if m.K < len(m.Buffer) {
		history = m.Buffer[len(m.Buffer)-m.K:]
	}
	return map[string]string{
		m.MemoryKey: strings.Join(history, "\n"),
	}
}

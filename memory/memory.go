package memory

import (
	"container/list"
	"sync"

	// "github.com/charmbracelet/log"

	oa "github.com/sashabaranov/go-openai"
)

type MemoryManager struct {
	history *list.List
	mutex   sync.Mutex
}

func NewMemoryManager() *MemoryManager {
	return &MemoryManager{
		history: list.New(),
	}
}

func (m *MemoryManager) AddMessage(message oa.ChatCompletionMessage) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// log.Debug("[memory] message added", "message", message)
	m.history.PushBack(message)
}

func (m *MemoryManager) GetHistory() []oa.ChatCompletionMessage {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	history := make([]oa.ChatCompletionMessage, 0, m.history.Len())
	for e := m.history.Front(); e != nil; e = e.Next() {
		history = append(history, e.Value.(oa.ChatCompletionMessage))
	}

	// log.Debug("[memory] get history called", "history", history)

	return history
}
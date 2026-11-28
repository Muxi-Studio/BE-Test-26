package store

import "sync"

// Todo 是待办条目。
type Todo struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Finished bool   `json:"finished"`
}

// Memory 是进程内的待办存储。
type Memory struct {
	mu     sync.Mutex
	nextID int
	items  []Todo
}

func NewMemory() *Memory {
	return &Memory{nextID: 1}
}

func (m *Memory) List() []Todo {
	m.mu.Lock()
	defer m.mu.Unlock()

	out := make([]Todo, len(m.items))
	copy(out, m.items)
	return out
}

func (m *Memory) Add(title string) Todo {
	m.mu.Lock()
	defer m.mu.Unlock()

	item := Todo{ID: m.nextID, Title: title}
	m.nextID++
	m.items = append(m.items, item)
	return item
}

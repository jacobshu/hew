package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Message struct {
	Role    string    `json:"role"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
}

type Chat struct {
	ID       string    `json:"id"`
	Messages []Message `json:"messages"`
}

type ChatManager struct {
	StorageDir string
}

func New(storageDir string) (*ChatManager, error) {
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &ChatManager{StorageDir: storageDir}, nil
}

func (cm *ChatManager) CreateChat() (*Chat, error) {
	id := fmt.Sprintf("chat_%d", time.Now().UnixNano())
	chat := &Chat{ID: id, Messages: []Message{}}
	if err := cm.SaveChat(chat); err != nil {
		return nil, err
	}
	return chat, nil
}

func (cm *ChatManager) LoadChat(id string) (*Chat, error) {
	data, err := os.ReadFile(cm.chatFilePath(id))
	if err != nil {
		return nil, fmt.Errorf("failed to read chat file: %w", err)
	}

	var chat Chat
	if err := json.Unmarshal(data, &chat); err != nil {
		return nil, fmt.Errorf("failed to unmarshal chat data: %w", err)
	}

	return &chat, nil
}

func (cm *ChatManager) SaveChat(chat *Chat) error {
	data, err := json.Marshal(chat)
	if err != nil {
		return fmt.Errorf("failed to marshal chat data: %w", err)
	}

	if err := os.WriteFile(cm.chatFilePath(chat.ID), data, 0644); err != nil {
		return fmt.Errorf("failed to write chat file: %w", err)
	}

	return nil
}

func (cm *ChatManager) AddMessage(chat *Chat, role, content string) error {
	message := Message{
		Role:    role,
		Content: content,
		Time:    time.Now(),
	}
	chat.Messages = append(chat.Messages, message)
	return cm.SaveChat(chat)
}

func (cm *ChatManager) chatFilePath(id string) string {
	return filepath.Join(cm.StorageDir, id+".json")
}

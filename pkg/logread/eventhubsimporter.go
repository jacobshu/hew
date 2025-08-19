package logread

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type EHILogEntry struct {
	Logger         string `json:"logger"`
	Level          string `json:"level"`
	TimestampLocal string `json:"timestampLocal"`
	Message        string `json:"message"`
	Timestamp      string `json:"timestamp"`
	Thread         string `json:"thread"`
	MachineName    string `json:"machineName"`
}

type EHIProcessedMessage struct {
	UID                   string             `json:"uid,omitempty"`
	TS                    string             `json:"ts,omitempty"`
	TVS                   map[string]float64 `json:"tvs,omitempty"`
	EventProcessedUtcTime string             `json:"EventProcessedUtcTime,omitempty"`
	PartitionId           int                `json:"PartitionId,omitempty"`
	EventEnqueuedUtcTime  string             `json:"EventEnqueuedUtcTime,omitempty"`
}

func ProcessEHILogFile(filename string) ([]EHILogEntry, error) {
	fullPath := filepath.Join("C:\\ProgramData\\zdScada\\Logs\\zdUserSpace\\EventHubsImporter", filename)

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var entries []EHILogEntry
	scanner := bufio.NewScanner(file)

	const maxCapacity = 1024 * 1024 // 1MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "" {
			continue
		}

		var entry EHILogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			fmt.Printf("Error parsing line: %v\n", err)
			continue
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return entries, fmt.Errorf("error reading file: %v", err)
	}

	return entries, nil
}

func ParseEHIMessageJSON(entry EHILogEntry) (*EHIProcessedMessage, error) {
	if !strings.HasPrefix(entry.Message, "{") {
		return nil, nil // Not JSON, no error
	}

	jsonStr := strings.ReplaceAll(entry.Message, "\\u0022", "\"")

	var processed EHIProcessedMessage
	err := json.Unmarshal([]byte(jsonStr), &processed)
	if err != nil {
		return nil, fmt.Errorf("error parsing message JSON: %v", err)
	}

	return &processed, nil
}

func getEntriesWithDuplicateTimestamps(entries []EHILogEntry) (map[string]map[string][]*EHIProcessedMessage, error) {
	duplicates := make(map[string]map[string][]*EHIProcessedMessage)
	for _, entry := range entries {
		message, err := ParseEHIMessageJSON(entry)
		if err != nil {
			return nil, fmt.Errorf("error parsing message: %w", err)
		}

		if message == nil || message.TS == "" {
			continue
		}

		if duplicates[message.UID] == nil {
			duplicates[message.UID] = map[string][]*EHIProcessedMessage{}
		}
		duplicates[message.UID][message.TS] = append(duplicates[message.UID][message.TS], message)
	}

	for _, msgs := range duplicates {
		for ts, msgArr := range msgs {
			if len(msgArr) <= 1 {
				delete(msgs, ts)
			}
		}
	}

	return duplicates, nil
}

func FilterEHILogsByDevice(entries []EHILogEntry, deviceID string) []EHILogEntry {
	var filtered []EHILogEntry

	for _, entry := range entries {
		if strings.Contains(entry.Message, deviceID) {
			filtered = append(filtered, entry)
		}
	}

	return filtered
}

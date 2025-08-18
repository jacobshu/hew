package logread

import (
	"bufio"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"
)

type LogEntry struct {
	Logger         string `json:"logger"`
	Level          string `json:"level"`
	TimestampLocal string `json:"timestampLocal"`
	Message        string `json:"message"`
	Timestamp      string `json:"timestamp"`
	Thread         string `json:"thread"`
	MachineName    string `json:"machineName"`
}

type ProcessedMessage struct {
	UID                   string             `json:"uid,omitempty"`
	TS                    string             `json:"ts,omitempty"`
	TVS                   map[string]float64 `json:"tvs,omitempty"`
	EventProcessedUtcTime string             `json:"EventProcessedUtcTime,omitempty"`
	PartitionId           int                `json:"PartitionId,omitempty"`
	EventEnqueuedUtcTime  string             `json:"EventEnqueuedUtcTime,omitempty"`
}

func ProcessLogFile(filename string) ([]LogEntry, error) {
	// Construct the full path to the log file
	fullPath := filepath.Join("C:\\ProgramData\\zdScada\\Logs\\zdUserSpace\\EventHubsImporter", filename)

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)

	const maxCapacity = 1024 * 1024 // 1MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "" {
			continue
		}

		var entry LogEntry
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

func ParseMessageJSON(entry LogEntry) (*ProcessedMessage, error) {
	if !strings.HasPrefix(entry.Message, "{") {
		return nil, nil // Not JSON, no error
	}

	jsonStr := strings.ReplaceAll(entry.Message, "\\u0022", "\"")

	var processed ProcessedMessage
	err := json.Unmarshal([]byte(jsonStr), &processed)
	if err != nil {
		return nil, fmt.Errorf("error parsing message JSON: %v", err)
	}

	return &processed, nil
}

func getEntriesWithDuplicateTimestamps(entries []LogEntry) (map[string][]*ProcessedMessage, error) {
	duplicates := make(map[string][]*ProcessedMessage)
	for _, entry := range entries {
		message, err := ParseMessageJSON(entry)
		if err != nil {
			return nil, fmt.Errorf("error parsing message: %w", err)
		}

		if message == nil || message.TS == "" {
			continue
		}

		if _, ok := duplicates[message.TS]; ok {
			duplicates[message.TS] = append(duplicates[message.TS], message)
		} else {
			duplicates[message.TS] = []*ProcessedMessage{message}
		}
	}

	for k, msgs := range duplicates {
		if len(msgs) <= 1 {
			delete(duplicates, k)
		}
	}

	return duplicates, nil
}

func FilterLogsByDevice(entries []LogEntry, deviceID string) []LogEntry {
	var filtered []LogEntry

	for _, entry := range entries {
		if strings.Contains(entry.Message, deviceID) {
			filtered = append(filtered, entry)
		}
	}

	return filtered
}

func LogRead(path string) {
	logEntries, err := ProcessLogFile(path)
	if err != nil {
		fmt.Printf("Error processing log file: %v\n", err)
		return
	}

	fmt.Printf("Processed %d log entries\n", len(logEntries))

	duplicates, err := getEntriesWithDuplicateTimestamps(logEntries)
	if err != nil {
		fmt.Printf("error in getEntriesWithDuplicateTimestamps: %v", err)
	}

	for ts, entries := range duplicates {
		str := ""
		hasDuplicate := false
		for i := range entries {
			for j := len(entries) - 1; j > len(entries)-i; j-- {
				if entries[i].UID == entries[j].UID {
					if maps.Equal(entries[i].TVS, entries[j].TVS) {
						hasDuplicate = true
						str += fmt.Sprintf("%s,%s,", ts, entries[i].UID)
						for k, v := range entries[i].TVS {
							str += fmt.Sprintf("%v:%v,", k, v)
						}
					}
				}
			}
		}
		if hasDuplicate {
			fmt.Println(str)
		}
	}
	// filtered := FilterLogs(logEntries, "410776", "1")
	//
	// if len(filtered) > 0 {
	// 	// Example of processing a specific entry
	// 	processed, err := ParseMessageJSON(filtered[0])
	// 	if err != nil {
	// 		fmt.Printf("Error parsing message: %v\n", err)
	// 	} else if processed != nil {
	// 		// Print or further process the parsed message
	// 		fmt.Printf("Device UID: %s\n", processed.UID)
	// 		fmt.Printf("Timestamp: %s\n", processed.TS)
	//
	// 		// Example: Print all TVS values
	// 		if len(processed.TVS) > 0 {
	// 			fmt.Println("TVS Values:")
	// 			for key, value := range processed.TVS {
	// 				fmt.Printf("  %s: %.1f\n", key, value)
	// 			}
	// 		}
	// 	}
	//}
}

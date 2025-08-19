package logread

import (
	"fmt"
)

func LogRead(path string) {
	logEntries, err := ProcessEHILogFile(path)
	if err != nil {
		fmt.Printf("Error processing log file: %v\n", err)
		return
	}

	fmt.Printf("Processed %d log entries\n", len(logEntries))

	duplicates, err := getEntriesWithDuplicateTimestamps(logEntries)
	if err != nil {
		fmt.Printf("error in getEntriesWithDuplicateTimestamps: %v", err)
	}

	for uid, entries := range duplicates {
		fmt.Printf("%v:\n", uid)
		if len(entries) == 0 {
			fmt.Println("\tno duplicate timestamps")
		}
		for ts, entry := range entries {
			fmt.Printf("\t%v\n", ts)
			for _, msg := range entry {
				for _, f := range msg.TVS {
					fmt.Printf("\t\t%v\n", f)
				}
			}
		}
	}
}

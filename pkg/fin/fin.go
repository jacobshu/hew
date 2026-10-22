package fin

import (
	"bufio"
	"fmt"
	"os"
)

// length; the header row, if any, is returned as the first record.
func ReadCSV(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewScanner(f)
	lines := []string{}
	for r.Scan() {
		line := r.Text()
		lines = append(lines, line)
	}
	if err := r.Err(); err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	return lines, nil
}

package fin

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
)

//go:embed finance.json
var finConfig []byte

type LedgerConfig struct {
	RootParts        []string `json:"root"`
	SourcesDir       []string `json:"sourcesDir"`
	RulesDir         []string `json:"rulesDir"`
	MainJournalParts []string `json:"main"`
	MainJournal      string
	Rules            map[string]string
	Root             string
	Sources          []string
}

type RulesMap map[string]string

func resolve(parts []string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return path.Join(append([]string{home}, parts...)...)
}

func GetLedgerConfig() (LedgerConfig, error) {
	var cfg LedgerConfig
	if err := json.Unmarshal(finConfig, &cfg); err != nil {
		return LedgerConfig{}, err
	}

	cfg.Root = resolve(cfg.RootParts)
	cfg.MainJournal = resolve(append(cfg.RootParts, cfg.MainJournalParts...))

	srcs := resolve(append(cfg.RootParts, cfg.SourcesDir...))
	entries, err := os.ReadDir(srcs)
	if err != nil {
		return LedgerConfig{}, fmt.Errorf("error getting source files: %w", err)
	}

	var files []string
	for _, e := range entries {
		if e.Type().IsRegular() {
			files = append(files, path.Join(srcs, e.Name()))
		}
	}

	var rules map[string][]string
	for k, r := range cfg.Rules {
		fmt.Println(k, r)
		rules[k] = append(rules[k], resolve(append(cfg.RootParts, r)))
	}

	cfg.Sources = files
	return cfg, nil
}

func ExecHledger(args []string) error {
	cmd := exec.Command("hledger", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running hledger: %w", err)
	}
	return nil
}

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

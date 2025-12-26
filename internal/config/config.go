package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	DSN        string
	OutputDir  string
	Package    string
	Tables     []string
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := &Config{
		OutputDir: ".",
		Package:   "model",
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch key {
		case "dsn":
			cfg.DSN = val
		case "output":
			cfg.OutputDir = val
		case "package":
			cfg.Package = val
		case "tables":
			for _, t := range strings.Split(val, ",") {
				if t = strings.TrimSpace(t); t != "" {
					cfg.Tables = append(cfg.Tables, t)
				}
			}
		}
	}
	return cfg, scanner.Err()
}

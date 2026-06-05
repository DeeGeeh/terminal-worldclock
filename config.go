package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type config struct {
	Zones [2]string `json:"zones"`
}

func configPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = "."
	}
	return filepath.Join(dir, "worldclock", "config.json")
}

func loadConfig() config {
	c := config{Zones: [2]string{"Europe/Helsinki", "America/Los_Angeles"}}
	b, err := os.ReadFile(configPath())
	if err == nil {
		json.Unmarshal(b, &c)
	}
	return c
}

func (c config) save() {
	p := configPath()
	os.MkdirAll(filepath.Dir(p), 0o755)
	if b, err := json.MarshalIndent(c, "", "  "); err == nil {
		os.WriteFile(p, b, 0o644)
	}
}

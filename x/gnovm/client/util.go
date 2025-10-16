package client

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

const gnoModName = "gnomod.toml"

// GnoMod represents the gnomod.toml for a Gno package.
type GnoMod struct {
	Module string `toml:"module"`
	Gno    string `toml:"gno"`
}

// parseGnoMod parses a gnomod.toml file and returns the configuration.
func parseGnoMod(configPath string) (*GnoMod, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config GnoMod
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse TOML config: %w", err)
	}

	return &config, nil
}

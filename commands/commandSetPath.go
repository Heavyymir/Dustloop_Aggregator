package commands

import (
	"fmt"
	"os"
	"github.com/Heavyymir/Dustloop_Aggregator/config"
)

func CommandSetPath(cfg *config.Config, args ...string) error {
	if len(args) < 2 || args[0] != "path" {
		return fmt.Errorf("usage: set path <directory>")
	}

	targetDir := args[1]
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("invalid path '%s': %w, targetDir, err")
	}

	cfg.DataDir = targetDir
	fmt.Printf("Data directory updated to: %s\n", cfg.DataDir)
	return nil
}


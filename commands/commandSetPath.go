package commands

import (
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/Heavyymir/Dustloop_Aggregator/config"
)



func commandSet(cfg *config.Config, args ...string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: set path <directory> or set db <sqlite_file>")
	}

	subcommand := args[0]
	target := args[1]

	switch subcommand{
	case "path":
		return handleSetPath(cfg, target)
	case "db":
		return handleSetDB(cfg, target)
	default:
		return fmt.Errorf("unkown set target '%s'. Available: path, db", subcommand)
	}
}

func handleSetPath(cfg *config.Config, rawPath string) error {
	expanded, err := config.ExpandTilde(rawPath)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	cleanDir := filepath.Clean(expanded)

	// Ensure directory exists
	if err := os.MkdirAll(cleanDir, 0755); err != nil {
		return fmt.Errorf("invalid directory '%s': %w", cleanDir, err)
	}

	cfg.DataDir = cleanDir
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Data directory updated to: %s\n", cfg.DataDir)
	return nil
}

func handleSetDB(cfg *config.Config, rawPath string) error {
	expanded, err := config.ExpandTilde(rawPath)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	cleanPath := filepath.Clean(expanded)

	// Ensure parent dir exists
	parentDir := filepath.Dir(cleanPath)
	if parentDir != "." && parentDir != "" {
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Switch DB
	newDB, _, err := config.SwitchDatabase(cfg.DB, cleanPath)
	if err != nil {
		return fmt.Errorf("failed to switch database: %w", err)
	}

	cfg.DB = newDB
	cfg.DBPath = cleanPath

	fmt.Printf("Active database switched to: %s\n", cleanPath)
	return nil
}

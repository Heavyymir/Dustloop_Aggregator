package config

import(
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"io"
	"path/filepath"
	"strings"
	"database/sql"

	"github.com/chzyer/readline"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/storage/sqlite"
)

type PathCompleter struct{}

// Get config path returns the path to the config file in a users config directory
func getConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir, err = os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not find user director: %w", err)
		}
	}

	appDir := filepath.Join(configDir, "dustloop_aggregator")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", fmt.Errorf("could not create config directory: %w", err)
	}
	return filepath.Join(appDir, "config.json"), nil
}

// Load reads the config file from disk
func Load() (*Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}

// Save writes the config file to disk
func Save(cfg *Config) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

// PromptNewDBPath asks the user to input a save location for the SQLite Database
func PromptNewDBPath(promptMessage string) (string, error) {
	if promptMessage != "" {
		fmt.Println(promptMessage)
	}

	defaultPath, _ := getDefaultPath()
	promptText := fmt.Sprintf("Enter path for SQLite database [default: %s] >", defaultPath)

	// Configure readline instance
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          promptText,
		AutoComplete:	 &PathCompleter{},
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})

	if err != nil {
		return "", fmt.Errorf("initialise readline: %w", err)
	}

	defer rl.Close()

	for {
		line, err := rl.Readline()
		if errors.Is(err, readline.ErrInterrupt) {
			return "", fmt.Errorf("operation canceled by user (Ctrl + C)")
		} else if errors.Is(err, io.EOF) {
			return "", fmt.Errorf("operation cancelled by user (EOF)")
		} else if err != nil {
			return "", fmt.Errorf("read line error: %w", err)
		}

		cleanPath := strings.TrimSpace(line)

		// Handle empty inputs with defaultPath
		if cleanPath == "" {
			if defaultPath == "" {
				fmt.Println("Path cannot be empty. Please try again")
				continue
			}
			cleanPath = defaultPath
			fmt.Printf("Using default location: %s\n", cleanPath)
		} else {
			// Expand tilde if users entered path starts with (~)
			expanded, err := ExpandTilde(cleanPath)
			if err != nil {
				fmt.Printf("error expanding path: %v\n", err)
				continue
			}
			cleanPath = expanded
		}

		// Clean and resolve path
		cleanPath = filepath.Clean(cleanPath)

		//Check if user has entered an existing directory
		info, err := os.Stat(cleanPath)
		if err == nil && info.IsDir() {
			// If it is an existing directory, append a default base filename
			cleanPath = filepath.Join(cleanPath, "dustloop.db")
			fmt.Printf("Directory provided. Using file: %s\n", cleanPath)
		} else if strings.HasSuffix(line, "/") || strings.HasSuffix(line, "\\") {
			// Or if the user typed the path with a trailing slash for a dir that doesn't exist
			cleanPath = filepath.Join(cleanPath, "dustloop.db")
			fmt.Printf("Directory path detected. Using file: %s\n", cleanPath)
		}

		// Ensure parent directories exist
		parentDir := filepath.Dir(cleanPath)
		if parentDir != "." && parentDir != "" {
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				fmt.Printf("Error: Unable to create directory %s: %v\n", parentDir, err)
				continue
			}
		}

		return cleanPath, nil
	}
}


// Setup Database handles loading of an existing DB or the creation of a new one
func SetupDatabase() (*sql.DB, *Config, error) {
	cfg, err := Load()
	needPrompt := false
	var promptReason string
	
	if errors.Is(err, os.ErrNotExist) {
		needPrompt = true
		promptReason = "No configuration file found"	
	} else if err != nil {
		return nil, nil, fmt.Errorf("read config: %w", err)
	} else {
		// Configuration file exists, check to see if DB file exists
		info, statErr := os.Stat(cfg.DBPath)
		if errors.Is(statErr, os.ErrNotExist) {
			needPrompt = true
			promptReason = fmt.Sprintf("Database file not found at configured path: %s", cfg.DBPath)
		} else if statErr == nil && info.IsDir() {
			needPrompt = true
			promptReason = fmt.Sprintf("Configured path '%s' is a directory, not a database file.", cfg.DBPath)
		}
	}

	if needPrompt {
		newPath, err := PromptNewDBPath(promptReason)
		if err != nil {
			return nil, nil, err
		}

		cfg = &Config{DBPath: newPath}
		if err := Save(cfg); err != nil {
			return nil, nil, fmt.Errorf("save config: %w", err)
		}

		// Open creates the file if it does not already exist
		db, err := sqlite.Open(cfg.DBPath)
		if err != nil {
			return nil, nil, fmt.Errorf("open sqlite db: %w", err)
		}

		// Apply the required DB tables and schema
		if err := sqlite.InitialiseSchema(db); err != nil {
			db.Close()
			return nil, nil, fmt.Errorf("initialise schema: %w", err)
		}

		fmt.Printf("Initialised new database and schema at: %s", cfg.DBPath)
		return db, cfg, nil
	}

	// Normal Action: Database exists at cfg.DBPath
	db, err := sqlite.Open(cfg.DBPath)
	if err != nil {
		return nil, nil, fmt.Errorf("open sqlite db: %w", err)
	}

	return db, cfg, nil
}

func (p *PathCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	str := string(line[:pos])

	// Determine the directory to look in and the prefix to match
	var dir, prefix string
	if idx := strings.LastIndexAny(str, "/\\"); idx != -1 {
		dir = str[:idx+1]
		prefix = str[idx+1:]
	} else {
		dir = "."
		prefix = str
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, 0
	}

	var candidates [][]rune
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			suffix := name[len(prefix):]
			if entry.IsDir() {
				suffix += "/"
			}
			candidates = append(candidates, []rune(suffix))
		}
	}

	return candidates, len(prefix)
}

//Switch DataBase closes the current DB, saves a new path and opens a new DB at that location
func SwitchDatabase(currentDB *sql.DB, newPath string) (*sql.DB, *Config, error) {
	expandedPath, err := ExpandTilde(strings.TrimSpace(newPath))
	if err != nil {
		return nil, nil, err
	}
	expandedPath = filepath.Clean(expandedPath)

	// Ensure parent Dir exists
	parentDir := filepath.Dir(expandedPath)
	if parentDir != "." && parentDir != "" {
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return nil, nil, fmt.Errorf("open database: %w", err)
		}
	}

	// Open the new DB
	newDB, err := sqlite.Open(expandedPath)
	if err != nil {
		return nil, nil, fmt.Errorf("open new database: %w", err)
	}

	// Ensure the schema is present
	if err := sqlite.InitialiseSchema(newDB); err != nil {
		newDB.Close()
		return nil, nil, fmt.Errorf("initialise schema: %w", err)
	}

	// Save to config
	newCfg := &Config{DBPath: expandedPath}
	if err := Save(newCfg); err != nil {
		newDB.Close()
		return nil, nil, fmt.Errorf("save config: %w", err)
	}

	// Close old database connection only after the new one has succeeded
	if currentDB != nil {
		currentDB.Close()
	}

	return newDB, newCfg, nil
}

// Helper to help handle tilde (~) inputs for paths
func ExpandTilde (path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine user home directory: %w", err)
	}

	// Handles "~" or ~/path/to/db
	if path == "~" {
		return homeDir, nil
	}

	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "\\~") {
		return filepath.Join(homeDir, path[2:]), nil
	}

	return path, nil
}

// Sets a default path for DB
func getDefaultPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".dustloop", ".dustloop.db"), nil
}

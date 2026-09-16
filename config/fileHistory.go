package config

import(
	"os"
	"path/filepath"
)

func GetHistoryFilePath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir, err = os.UserHomeDir()
		if err != nil {
			return "" // if both fail, disable file history gracefully
		}
	}

	appDir := filepath.Join(configDir, "dustloop_aggregator")
	_ = os.MkdirAll(appDir, 0755)
	return filepath.Join(appDir, "history")
}

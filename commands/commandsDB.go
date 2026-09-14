package commands

import(
	"fmt"
	
	"github.com/Heavyymir/Dustloop_Aggregator/config"
	
)

func commandDBInfo(cfg *config.Config, args ...string) error {
	savedCfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("could not read config: %w", err)
	}

	fmt.Printf("Active SQLite database: %s\n", savedCfg.DBPath)
	return nil
}


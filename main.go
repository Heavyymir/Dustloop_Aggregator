package main

import (
	"log"
	
	"github.com/Heavyymir/Dustloop_Aggregator/config"
)

func main() {
	// Call to open local sqlite DB to hold data
	db, cfg, err := config.SetupDatabase()
	if err != nil {
		log.Fatalf("initialisation failed: %v", err)
	}
	defer db.Close()
	
	log.Printf("Connected to database at: %s", cfg.DBPath)
	
	// Call to start CLI loop
	startRepl(db)
}

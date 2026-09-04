package main

import(
	"log"
	
	"github.com/Heavyymir/Dustloop_Aggregator/internal/storage/sqlite"
)
func main() {
	// Call to open local sqlite DB to hold data
	db, err := sqlite.Open("chardata.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Initialise the require SQL table schema to store scraped data locally
	if err := sqlite.InitialiseSchema(db); err != nil {
		log.Fatal(err)
	}

	// Call to start CLI loop
	startRepl(db)
}

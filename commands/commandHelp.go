package commands

import (
	"fmt"

	"github.com/Heavyymir/Dustloop_Aggregator/config"
)

func commandHelp(cfg *config.Config, args ...string) error {
	fmt.Println(`Welcome to the Character Data Aggregator

Usage:
help: Displays a list of commands usable in the Character Data Aggregator
exit: Exits the program
select: Allows user to select a wiki and game. Usage: select * | select <wiki> * | select <wiki> <game>
fetch: Fetches a character page for a game. Saves fetched character data to database. Usage: fetch <character>
discover: Discovers character names for a game, and saves to a local json
list: lists characters present in the database.
frames: Displays framedata for a saved character from the database in the console. Usage frames <character>
set path: allows a user to set the path for json file saves from the discover command. Usage: set path <directory>`)

	return nil
}

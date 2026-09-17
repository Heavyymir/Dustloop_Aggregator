package commands

import (
	"fmt"

	"github.com/Heavyymir/Dustloop_Aggregator/config"
)

func commandHelp(cfg *config.Config, args ...string) error {
	fmt.Println(`Welcome to the Character Data Aggregator

Usage:
help: Displays a list of usable commands.
exit: Exits the program.
select: Allows user to select a wiki and game. Usage: select * | select <wiki> * | select <wiki> <game>.
fetch: Used after select, fetches a character page and saves to SQLite DB. Usage: fetch <character> [-v--|verbose] [-d|--details] [-f|--force]
• [-v|--verbose] Use to print move frame data tables.
• [-d|--details] Use to print move frame data tables, description and notes.
• [-f|--force] Use to bypass the character ID check that aims to prevent frivolous scraping of wikis.
discover: Use after select, discovers character names , and saves to a local JSON. Usage discover [-n|--no-json]
• [-n|--no-json] Use to prevent a local JSON save 
list: Used after selecting a game, lists characters present in the database.
frames: Used after selecting a game, Displays framedata for a saved character. Usage frames <character> [-d|--details]
• [-d|--details] Use to print move description and notes.
set path: Allows a user to set the path for json file saves from the discover command. Usage: set path <directory>
set db: Allows user to switch to a new directory for the internal SQlite DB file. This creates a fresh DB. Usage set db <directory>
dbinfo: Dispalys the location of the internal SQlite DB`)

return nil
}

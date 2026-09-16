package main

import (
	"database/sql"
	"fmt"
	"github.com/Heavyymir/Dustloop_Aggregator/catalog"
	"github.com/Heavyymir/Dustloop_Aggregator/commands"
	"github.com/Heavyymir/Dustloop_Aggregator/config"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/api"
	"github.com/chzyer/readline"
	"strings"
)

// REPL for command line interaction
func startRepl(db *sql.DB) {
	cfg := config.Config{
		CharDataClient: api.NewClient(),
		DB:             db,
		DataDir:		"./data",
	}

	// Initialise the completer to handle tab completion of internal commands
	completer := readline.NewPrefixCompleter(
		readline.PcItem("help"),
		readline.PcItem("select", wikiCompleter()...),
		readline.PcItem("fetch", gameCompleters()...),
		readline.PcItem("exit"),
		readline.PcItem("discover"),
		readline.PcItem("frames", gameCompleters()...),
		readline.PcItem("list", gameCompleters()...),
		readline.PcItem("db-info"),
		readline.PcItem("set",
			readline.PcItem("path"),
			readline.PcItem("db"),
		),
	)

	// Start the REPL
	fmt.Println("Welcome to the Dustloop Aggregator. Please type 'help' for a list of commands.")
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "Dustloop Aggregator CLI > ",
		AutoComplete: &replCompleter{
			cmdCompleter:	completer,
			pathCompleter:  &config.PathCompleter{},
		},
		HistoryFile:	config.GetHistoryFilePath(),
		HistoryLimit:	1000,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	// Defer close of the initialised REPL loop
	defer rl.Close()

	// Initialise an infinate loop to continue the REPL until the Exit command is called
	for {
		text, err := rl.Readline()
		if err != nil {
			return
		}

		words := cleanInput(text)
		if len(words) == 0 {
			fmt.Println("Empty Command. Please enter a valid command. Use command `help` for a list of valid commands")
			continue
		}

		if words[0] == "exit" {
			return
		}

		command, exists := commands.GetCommands()[words[0]]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}

		if err := command.Callback(&cfg, words[1:]...); err != nil {
			fmt.Println(err)
		}
	}
}

// Cleans user inputs
func cleanInput(text string) []string {
	lowerCase := strings.ToLower(text)
	fields := strings.Fields(lowerCase)
	return fields
}

// Wrapper function to enable tab completion of Games and Wikis
func wikiCompleter() []readline.PrefixCompleterInterface {
	// "*" supports: User input of select *
	items := []readline.PrefixCompleterInterface{
		readline.PcItem("*"),
	}
	// Add each wiki and its valid game choices from the catalog map
	for wikikey, wiki := range catalog.Wikis {
		// "*" supports: user input of select <wiki> *. Example Command: select dustloop *
		children := []readline.PrefixCompleterInterface{
			readline.PcItem("*"),
		}

		// Add each valid game for the selected wiki
		for gamekey := range wiki.Games {
			children = append(children, readline.PcItem(gamekey))
		}

		// Attach game completions underneath wiki completions. Structure: select <wiki> <game>
		items = append(items, readline.PcItem(
			wikikey,
			children...,
		))
	}
	// return the completed select command tree
	return items
}

func gameCompleters() []readline.PrefixCompleterInterface {
	var items []readline.PrefixCompleterInterface
	// Set to avoid duplicates if multiple wikis contain a game key
	seen := make(map[string]bool)

	for _, wiki := range catalog.Wikis {
		for gameKey := range wiki.Games {
			if !seen[gameKey] {
				seen[gameKey] = true
				items = append(items, readline.PcItem(gameKey))
			}
		}
	}
	return items
}

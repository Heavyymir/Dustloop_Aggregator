package commands

import "github.com/Heavyymir/Dustloop_Aggregator/config"

type CliCommand struct {
	Name        string
	Description string
	Callback    func(*config.Config, ...string) error
}

func GetCommands() map[string]CliCommand {
	return map[string]CliCommand{
		"exit": {
			Name:        "exit",
			Description: "command to exit CharData_Aggregator",
			Callback:    commandExit,
		},
		"help": {
			Name:        "help",
			Description: "displays information on commands usable in CharData_Aggregator",
			Callback:    commandHelp,
		},
		"select": {
			Name:        "select",
			Description: "selects a wiki and game",
			Callback:    commandSelect,
		},
		"fetch": {
			Name:        "fetch",
			Description: "fetches character data from a wiki",
			Callback:    commandFetch,
		},
		"discover": {
			Name:        "discover",
			Description: "Discovers playable characters for a selected game, and saves to a json",
			Callback:    commandDiscover,
		},
		"list": {
			Name:		 "list",
			Description: "Lists the Characters stored in the internal SQL database based on a selected game",
			Callback:	 commandList,	
		},
		"frames": {
			Name:		 "frames",
			Description: "pulls framedata for a selected character and prints to console",
			Callback:	 commandFrames,
		},
	}
}

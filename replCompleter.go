package main

import (
	"strings"
	"github.com/chzyer/readline"

	"github.com/Heavyymir/Dustloop_Aggregator/config"
)

type replCompleter struct {
	cmdCompleter	readline.AutoCompleter
	pathCompleter	*config.PathCompleter	
}

func (c *replCompleter) Do(line []rune, pos int) ([][]rune, int) {
	lineStr := string(line[:pos])

	// If the user is typing a path argument for set db or set path
	for _, prefix := range []string{"set db ", "set path "} {
		if strings.HasPrefix(lineStr, prefix) {
			// Pass only the path portion to your existing PathCompleter
			pathRunes := line[len(prefix):pos]
			candidates, length := c.pathCompleter.Do(pathRunes, len(pathRunes))
			return candidates, length
		}	
	}

	// Else, delegate to your existing commands completer
	return c.cmdCompleter.Do(line, pos)
}

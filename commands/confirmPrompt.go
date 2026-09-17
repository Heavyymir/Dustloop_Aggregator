package commands

import(
	"fmt"
	"strings"

	"github.com/chzyer/readline"
)

func confirmPrompt(rl *readline.Instance, promptText string) bool {
	defer rl.SetPrompt("Dustloop Aggregator CLI > ")

	rl.SetPrompt(fmt.Sprintf("%s [Y/N]: ", promptText))

	line, err := rl.Readline()
	if err != nil {
		return false
	}

	normalised := strings.ToLower(strings.TrimSpace(line))
	return normalised == "y" || normalised == "yes"
}

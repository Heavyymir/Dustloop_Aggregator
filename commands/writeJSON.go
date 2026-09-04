package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Heavyymir/Dustloop_Aggregator/internal/models"
)

// The below helper is designed to create a Json output for the program to store outputs
func writeMovesJSON(character string, moves []models.Move) error {
	filename := safeFileName(character) + ".json"

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", " ")

	if err := encoder.Encode(moves); err != nil {
		return err
	}

	fmt.Printf("Wrote  %s\n", filename)
	return nil
}

func safeFileName(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "_"))
}

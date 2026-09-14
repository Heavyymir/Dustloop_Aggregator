package config

import (
	"database/sql"
	"github.com/Heavyymir/Dustloop_Aggregator/catalog"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/api"
	"github.com/chzyer/readline"
)

type Config struct {
	CharDataClient 	*api.Client
	DB				*sql.DB
	Wiki           	catalog.Wiki
	Game           	catalog.Game
	RL				*readline.Instance
	CharacterPage  	string
	DataDir			string
	DBPath			string	`json:"db_path"`
}

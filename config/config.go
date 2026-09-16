package config

import (
	"database/sql"
	"github.com/Heavyymir/Dustloop_Aggregator/catalog"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/api"
	"github.com/chzyer/readline"
)

type Config struct {
	CharDataClient 	*api.Client 		`json:"-"`
	DB				*sql.DB 			`json:"-"`
	Wiki           	catalog.Wiki		`json:"-"`
	Game           	catalog.Game		`json:"-"`
	RL				*readline.Instance	`json:"-"`
	CharacterPage  	string				`json:"-"`
	DataDir			string				`json:"data_dir,omitempty"`
	DBPath			string				`json:"db_path"`
}

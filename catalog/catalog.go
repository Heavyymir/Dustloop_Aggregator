package catalog

// This Catalog is designed to hardcode Wiki Addresses and Games. More will be added over time.

type SourceType string

const (
	SourceHTML SourceType = "html"
	SourceJSON SourceType = "json"
)

// Define Game Elements
type Game struct {
	Name          string
	Slug          string
	SourceType    SourceType
	CharacterPath string
}

// Define Wiki Elements. Nested Map to assist in mapping game URLs.
type Wiki struct {
	Name  		string
	URL   		string
	Slug		string
	SourceType	SourceType
	Games map[string]Game
}

// Literal to hold Wiki Data
var Wikis = map[string]Wiki{
	"dustloop": {
		Name: "Dustloop",
		URL:  DustloopURL,
		Slug: "dustloop",
		Games: map[string]Game{
			"ggst": {
				Name:          "Guilty Gear Strive",
				Slug:          "ggst",
				CharacterPath: "GGST/{character}",
			},
			"bbcf": {
				Name:          "Blazblue Centralfiction",
				Slug:          "bbcf",
				CharacterPath: "BBCF/{character}",
			},
		},
	},
	"mizuumi": {
		Name:	"Mizuumi",
		URL:	MizuumiURL,
		Slug:	"mizuumi",
		Games:	map[string]Game{
			"uni2": {
				Name:			"Under Night IN-BIRTH II Sys:Celes",
				Slug:			"UNI2",
				CharacterPath: 	"Under_Night_In-Birth/UNI2/{character}",
			},
		},
	},
	"fat": {
		Name:	"FAT (Frame Advantage Tool)",
		URL:	FATURL,
		Slug:	"fat",
		Games: map[string]Game{
			"sf6": {
				Name:			"Street Fighter 6",
				Slug:			"sf6",
				CharacterPath:	"SF6/{character}/moves.json",
			},
			"sf5": {
				Name:			"Street Fighter V",
				Slug:			"sf5",
				CharacterPath:	"SFV/{character}/moves.json",
			},
			"usf4": {
				Name:			"Ultra Street Fighter IV",
				Slug:			"usf4",
				CharacterPath:	"USF4/{character}/moves.json",
			},
		},
	},
}	


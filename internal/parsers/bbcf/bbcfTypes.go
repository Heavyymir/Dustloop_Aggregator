package bbcf


import "github.com/Heavyymir/Dustloop_Aggregator/internal/models"


// Define Structs to handle Character data
type BBCFCharData struct {
	Name		string
	Overview	string
	Mechanics	[]Mechanic
	Moves		[]models.Move
	Infobox		map[string]string
}

type Mechanic struct {
	Name		string
	Description	string
}

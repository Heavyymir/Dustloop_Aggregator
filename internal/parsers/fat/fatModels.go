package fat

// FATMove represents a single move object in FAT's JSON feed
type FATMove struct {
	MoveName string `json:"moveName"`
	Input    string `json:"input"`
	Damage   string `json:"damage"`
	Startup  string `json:"startup"`
	Active   string `json:"active"`
	Recovery string `json:"recovery"`
	OnBlock  string `json:"onBlock"`
	OnHit    string `json:"onHit"`
	Cancel   string `json:"cancel"`
	Comments string `json:"comments"`
}

// FATCharacterPayload represents the top-level structure of character data in FAT
// Some FAT feeds use a slice of moves, others use a map or struct wrapper.
type FATCharacterPayload struct {
	Moves []FATMove `json:"moves"`
}

package world

// Optimize per comment: code design: data structure heavy
type T_World struct {
	// Geography and Alien Locations
	Atlas                map[int]map[int]map[int]string // [Key][X][Y]["City Name" or "ROAD"]
	Aliens               map[int]map[int][]T_Alien      // [X][Y][Alien Info]
	GridSizeX, GridSizeY int                            // Grid Size (Added Y For Optics)
}

type T_Alien struct {
	AlienID      int
	AlienSpecies string
}

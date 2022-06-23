package world

func (world *T_World) OnRoad(x, y int) bool {
	_, _, found := world.LookupAtlas(x, y)
	return found
}

func (world *T_World) BreakRoad(x, y int) {
	index, _, found := world.LookupAtlas(x, y)
	if found {
		world.Atlas[index][x][y] = ""
	}
}

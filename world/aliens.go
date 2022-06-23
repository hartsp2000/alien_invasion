package world

func (world *T_World) ListAliens() (xlist []int, ylist []int) {
	for X, _ := range world.Aliens {
		for Y, _ := range world.Aliens[X] {
			if len(world.Aliens[X][Y]) > 0 {
				xlist = append(xlist, X)
				ylist = append(ylist, Y)
			}
		}
	}
	return xlist, ylist
}

func (world *T_World) LookupAliens(x, y int) (invaders []T_Alien, found bool) {
	found = false
	if _, ok := world.Aliens[x][y]; ok {
		if len(world.Aliens[x][y]) > 0 {
			found = true
		}
		return world.Aliens[x][y], found
	}
	return []T_Alien{}, found
}

func (world *T_World) MoveAlien(aid, fromX, fromY, toX, toY int) {
	var inhabitants []T_Alien = world.Aliens[fromX][fromY]
	var newinhabitants []T_Alien
	world.CheckAliensInit(toX, toY)
	for _, alien := range inhabitants {
		if alien.AlienID == aid {
			world.Aliens[toX][toY] = append(world.Aliens[toX][toY], alien)
			// If there is more than one alien here, we want to preserve their location
			for stillhere, _ := range inhabitants {
				if inhabitants[stillhere].AlienID != aid {
					newinhabitants = append(newinhabitants, inhabitants[stillhere])
				}
			}
			delete(world.Aliens[fromX], fromY)
			world.CheckAliensInit(fromX, fromY)
			world.Aliens[fromX][fromY] = newinhabitants
		}
	}
}

func (world *T_World) CheckAliensInit(x, y int) {
	if _, ok := world.Aliens[x]; !ok {
		world.Aliens[x] = make(map[int][]T_Alien)
	}
}

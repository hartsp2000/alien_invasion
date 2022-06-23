package world

import (
	"errors"
)

func (world *T_World) ListAtlas() (xlist []int, ylist []int) {
	for key := 0; key < len(world.Atlas); key++ {
		for X, _ := range world.Atlas[key] {
			for Y, _ := range world.Atlas[key][X] {
				if len(world.Atlas[key][X][Y]) > 0 {
					xlist = append(xlist, X)
					ylist = append(ylist, Y)
				}
			}
		}
	}
	return xlist, ylist
}

func (world *T_World) LookupAtlas(x, y int) (index int, data string, found bool) {
	found = false
	for key := 0; key < len(world.Atlas); key++ {
		if _, ok := world.Atlas[key][x][y]; ok {
			if len(world.Atlas[key][x][y]) > 0 {
				found = true
			}
			return key, world.Atlas[key][x][y], found
		}
	}
	return -1, "", found
}

func (world *T_World) LookupCity(name string) (x int, y int, err error) {
	for key, _ := range world.Atlas {
		for x, _ := range world.Atlas[key] {
			for y, _ := range world.Atlas[key][x] {
				if _, ok := world.Atlas[key][x][y]; ok {
					if world.Atlas[key][x][y] == name {
						return x, y, nil
					}
				}
			}
		}
	}
	return -1, -1, errors.New("Can't find city in world map")
}

func (world *T_World) LookupAtlasIndex(index int) (x, y int, name string) {
	for X, _ := range world.Atlas[index] {
		for Y, _ := range world.Atlas[index][x] {
			if X == x && Y == y {
				return X, Y, world.Atlas[index][x][y]
			}
		}
	}
	return -1, -1, ""
}

func (world *T_World) CheckAtlasInit(index, x, y int) {
	if _, ok := world.Atlas[index]; !ok {
		world.Atlas[index] = make(map[int]map[int]string)
	}

	if _, ok := world.Atlas[index][x]; !ok {
		world.Atlas[index][x] = make(map[int]string)
	}
}

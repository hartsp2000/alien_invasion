package world

import (
	"fmt"
	"github.com/hartsp2000/alien_invasion/config"
	"time"
)

func RunSimulation(conf config.T_Config) {
	var err error
	var battlefield T_World

	if err = battlefield.Genesis(conf); err != nil {
		panic(err)
	}

	refresh, terr := time.ParseDuration(conf.Refresh)
	if terr != nil {
		panic(terr)
	}

	if conf.ShowMap {
		battlefield.DisplayWorld(false)
		fmt.Printf(string(config.Cgreen) + "Invasion Begins... Ships are landing...\n" + string(config.Creset))
		time.Sleep(refresh)
		battlefield.DisplayWorld(true)
	}

	// Let the battle begin
	for loop := 0; loop < conf.AllowedMoves; loop++ {
		battlefield.Fight()
		battlefield.Iterate()
		if conf.ShowMap {
			battlefield.DisplayWorld(true)
		}
		time.Sleep(refresh)
	}

	// The Apocalypse Happened
	if err = battlefield.Apocalypse(conf); err != nil {
		panic(err)
	}

	fmt.Printf("The battle is over and perhaps it's time to find a new planet.  Please see results saved in: %s\n\n", conf.MapOutfile)
}

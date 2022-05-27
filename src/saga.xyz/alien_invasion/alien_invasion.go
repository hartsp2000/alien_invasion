package main

import (
	"flag"
	"fmt"
	"os"
	"saga.xyz/alien_invasion/actions"
	"saga.xyz/alien_invasion/config"
	"time"
	"version"
)

func DisplayVersion() {
	fmt.Printf("Alien Invasion - A Wargames Simlated Invasion\n")
	fmt.Printf("by Sean Hart\n")
	fmt.Printf("Version: %s\n", version.VERSION)
	fmt.Printf("Build: %s\n", version.BUILDID)
	fmt.Printf("\n")
}

func DisplayHelp() {
	DisplayVersion()
	flag.PrintDefaults()
	fmt.Printf("\n")
}

func CommandLineArgs() (conf config.T_Config) {
	var help = flag.Bool("help", false, "Display usage")
	var ver = flag.Bool("version", false, "Display Version")
	var mapinfile = flag.String("inmap", "world.in", "Starting World Map Data File (starting/input)")
	var mapoutfile = flag.String("outmap", "world.out", "Oblierated World Map Data File (ending/outptu)")
	var aliens = flag.Int("aliens", 2, "Size of alien army")
	var allowedmoves = flag.Int("moves", 10000, "Maximum moves allowed by alien invaders")
	var showmap = flag.Bool("gfx", false, "Show graphics of the world (ANSI terminal required / Not recommended for large worlds)")
	var refresh = flag.String("refresh", "500ms", "Delay between moves")
	var distance = flag.Int("distance", 2, "Maximum distance between cities")

	flag.Parse()

	if *help {
		DisplayHelp()
		os.Exit(2)
	}

	if *ver {
		DisplayVersion()
		os.Exit(2)
	}

	conf.MapInfile = *mapinfile
	conf.MapOutfile = *mapoutfile
	conf.Aliens = *aliens
	conf.AllowedMoves = *allowedmoves
	conf.ShowMap = *showmap
	conf.Refresh = *refresh
	conf.Distance = *distance

	return conf
}

func main() {
	conf := CommandLineArgs()
	refresh, terr := time.ParseDuration(conf.Refresh)
	if terr != nil {
		panic(terr)
	}

	var battlefield actions.T_World
	var err error

	if err = battlefield.Genesis(conf); err != nil {
		panic(err)
	}

	fmt.Println(config.TITLE)

	if conf.ShowMap {
		battlefield.DisplayWorld(false)
		fmt.Printf(string(config.Cgreen) + "Invasion Begins... Ships are landing...\n" + string(config.Creset))
		time.Sleep(refresh)
		battlefield.DisplayWorld(true)
	}

	// Let the battle begin
	for loop := 0; loop < conf.AllowedMoves; loop++ {
		for !battlefield.CheckForDestruction() {
		} // Make sure check completes before moving to next iteration
		battlefield.Move()
		if conf.ShowMap {
			battlefield.DisplayWorld(true)
		}
		time.Sleep(refresh)
	}

	// The Apocalypse is OVER
	if err = battlefield.PostApocalypse(conf); err != nil {
		panic(err)
	}

	fmt.Printf("The battle is over and perhaps it's time to find a new planet.  Please see results saved in: %s\n\n", conf.MapOutfile)

}

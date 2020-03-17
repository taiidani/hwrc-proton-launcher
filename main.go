package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	// defaultProtonVersion stores the most stable reported Proton version for Homeworld
	defaultProtonVersion = "4.11"

	// applicationName matches the name of the binary produced that end-users will see
	applicationName = "hwrc-proton-launcher"

	// Steam ID of Homeworld Remastered in the library
	appID = 244160
)

var (
	// defaultSteamPaths contains a list of common locations that steam is found on a system
	defaultSteamPaths = []string{
		"$HOME/.local/share/Steam",
		"$HOME/.steam/steam",
	}

	flagWindowed bool
	flagModPath  string
	flagHelp     bool
	flagVerbose  bool
)

func main() {
	flag.BoolVar(&flagWindowed, "w", false, "Windowed mode")
	flag.BoolVar(&flagHelp, "help", false, "Display this help")
	flag.BoolVar(&flagVerbose, "v", false, "Display verbose output, for debugging purposes")
	flag.StringVar(&flagModPath, "m", "", "Load mod at the given `modPath` (see below for details)")
	flag.Parse()

	if flagVerbose {
		log.SetLevel(log.DebugLevel)
	}

	var fn func(*steam)
	if flagHelp {
		help()
		os.Exit(1)
	} else if len(flag.Args()) > 0 {
		fn = cli
	} else {
		fn = ui
	}

	s, err := newSteam()
	if err != nil {
		log.Fatal(err)
	}
	fn(s)
}

// help displays the help text for the application
func help() {
	fmt.Println(applicationName + " [OPTIONS] [GAME]")
	fmt.Println("")

	flag.PrintDefaults()

	fmt.Println(`
This is a launcher application for the Steam version of Homeworld Remastered Collection. It can be used as a replacement for Gearbox's launcher which has Linux compatibility issues.

WARNING: Steam must be running already before running this tool.

GAME:

	` + fmt.Sprintf("%-10s", hw1cla) + ` Homeworld 1 Classic
	` + fmt.Sprintf("%-10s", hw2cla) + ` Homeworld 2 Classic
	` + fmt.Sprintf("%-10s", hw1rem) + ` Homeworld 1 Remastered
	` + fmt.Sprintf("%-10s", hw2rem) + ` Homeworld 2 Remastered
	` + fmt.Sprintf("%-10s", hwmp) + ` Homeworld Remastered Multiplayer

Loading a mod:

The option '-m' requires a path to the mod file. The path can be an absolute path, or a relative path to the game DATAWORKSHOPSMODS folder. The following folders should be available:
- <STEAM PATH>/steamapps/common/Homeworld/HomeworldRM/DATAWORKSHOPMODS
- <STEAM PATH>/steamapps/common/Homeworld/Homeworld2Classic/DATAWORKSHOPMODS

For example, to load the Homeworld Remastered 2.3 Players Patch using a relative path (provided the mod file and any parent folder are in the DATAWORKSHOPSMODS folder):

	` + applicationName + ` hw1rem -m 1190476337/2.3PlayersPatch.big`)
}

// cli runs the game with the expectation that configuration is coming from the CLI arguments
func cli(s *steam) {
	// Figure out the game arguments
	game := flag.Arg(0)

	// Run the game!
	if err := s.run(game); err != nil {
		log.Fatal(err.Error())
	}
}

package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"syscall"

	ps "github.com/mitchellh/go-ps"
	"golang.org/x/exp/slog"
)

type steam struct {
	clientPath   string
	protonPath   string
	gameRootPath string
}

const (
	hw1cla = "hw1cla"
	hw2cla = "hw2cla"
	hw1rem = "hw1rem"
	hw2rem = "hw2rem"
	hwmp   = "hwmp"
)

// newSteam will create a new steam client interface and verify that Steam is currently running
func newSteam() (*steam, error) {
	var err error
	ret := &steam{}

	// Discover Steam and Proton
	ret.clientPath, err = findSteam()
	if err != nil {
		return ret, err
	}
	ret.protonPath = fmt.Sprintf("%s/steamapps/common/Proton %s", ret.clientPath, defaultProtonVersion)
	ret.gameRootPath = fmt.Sprintf("%s/steamapps/common/Homeworld", ret.clientPath)

	return ret, nil
}

// run will run the game!
func (s *steam) run(game string) error {
	gamePath, args := getSteamArguments(game, flagWindowed, flagModPath, s.gameRootPath)

	// Change directory to the location of the executable
	if err := os.Chdir(gamePath); err != nil {
		return fmt.Errorf("unable to change directory to game path: %w", err)
	}

	binary := fmt.Sprintf("%s/dist/bin/wine", s.protonPath)
	args = append([]string{binary, "steam.exe"}, args...)
	env := getSteamEnvironment(s.clientPath, s.protonPath)

	slog.Debug("Running", "binary", binary, "args", args)
	return syscall.Exec(binary, args, env)
}

// findSteam will locate the steam directory on disk in order to find the Homeworld binaries
// It will also determine if Steam is running or not, which is required for the app to execute.
//
// TODO: If not found, prompt the user to fill in the missing path
func findSteam() (string, error) {
	steamPath := ""
	for _, path := range defaultSteamPaths {
		resolvedPath := os.ExpandEnv(path)
		slog.Debug("Checking " + resolvedPath + " for Steam")
		if _, err := os.Stat(resolvedPath); err == nil {
			steamPath = resolvedPath
			break
		}
	}
	if steamPath == "" {
		return steamPath, fmt.Errorf("Unable to find Steam location")
	}

	// Is Steam running?
	processes, err := ps.Processes()
	if err != nil {
		return steamPath, err
	}

	steamFound := false
	for _, process := range processes {
		if process.Executable() == "steam" {
			steamFound = true
		}
	}

	if !steamFound {
		return steamPath, fmt.Errorf("Steam does not appear to be running")
	}

	return steamPath, nil
}

// getSteamArguments will determine the exact arguments to pass to steam in order to have the game run successfully
func getSteamArguments(game string, windowed bool, modPath string, gameRootPath string) (string, []string) {
	var gameExe string
	gameOptions := []string{}

	switch game {
	case hw1cla:
		gameExe = fmt.Sprintf("%s/Homeworld1Classic/exe/Homeworld.exe", gameRootPath)
		gameOptions = append(gameOptions, "/noglddraw")

		if windowed {
			gameOptions = append(gameOptions, "/window")
		}
	case hw2cla:
		gameExe = fmt.Sprintf("%s/Homeworld2Classic/Bin/Release/Homeworld2.exe", gameRootPath)

		if windowed {
			gameOptions = append(gameOptions, "-windowed")
		}

	case hw1rem:
		gameExe = fmt.Sprintf("%s/HomeworldRM.exe", getHomeworldRMPath())
		gameOptions = append(gameOptions, "-dlccampaign HW1Campaign.big")
		gameOptions = append(gameOptions, "-campaign HomeworldClassic")
		gameOptions = append(gameOptions, "-moviepath DataHW1Campaign")
		if windowed {
			gameOptions = append(gameOptions, "-windowed")
		}

	case hw2rem:
		gameExe = fmt.Sprintf("%s/HomeworldRM.exe", getHomeworldRMPath())
		gameOptions = append(gameOptions, "-dlccampaign HW2Campaign.big")
		gameOptions = append(gameOptions, "-campaign Ascension")
		gameOptions = append(gameOptions, "-moviepath DataHW2Campaign")
		if windowed {
			gameOptions = append(gameOptions, "-windowed")
		}

	case hwmp:
		gameExe = fmt.Sprintf("%s/HomeworldRM.exe", getHomeworldRMPath())
		if windowed {
			gameOptions = append(gameOptions, "-windowed")
		}

	default:
		help()
		os.Exit(1)
	}

	if modPath != "" {
		gameOptions = append(gameOptions, "-workshopmod "+modPath)
	}

	return path.Dir(gameExe), append([]string{gameExe}, gameOptions...)
}

// getSteamEnvironment provides the environment variables required by Steam
//
// WARNING: This section should be periodically checked and updated when new versions of Steam and Proton are released.
//
// These environment variables are set by Steam when launching the game from the Steam client using Steam Play. To retrieve
// these variables inspect the run script that can be dumped from Steam in this way:
//  1. Set the game launch options in the Steam client to
//     "PROTON_DUMP_DEBUG_COMMANDS=1 %command%";
//  2. Launch the game (even if it does not work);
//  3. Find the script "/tmp/proton_<username>/run".
func getSteamEnvironment(steamPath, protonPath string) []string {
	steamCustomPaths := strings.Join([]string{
		fmt.Sprintf("%s/dist/bin", protonPath),
		fmt.Sprintf("%s/ubuntu12_32/steam-runtime/amd64/bin", steamPath),
		fmt.Sprintf("%s/ubuntu12_32/steam-runtime/amd64/usr/bin/", steamPath),
		os.Getenv("PATH"),
	}, ":")

	ldLibraryPath := strings.Join([]string{
		fmt.Sprintf("%s/dist/lib64", protonPath),
		fmt.Sprintf("%s/dist/lib", protonPath),
		fmt.Sprintf("%s/ubuntu12_32/steam-runtime/pinned_libs_32", steamPath),
		fmt.Sprintf("%s/ubuntu12_32/steam-runtime/pinned_libs_64", steamPath),
		"/usr/lib/x86_64-linux-gnu/libfakeroot",
		"/lib/i386-linux-gnu",
		"/usr/local/lib",
		"/lib/x86_64-linux-gnu",
		"/lib",
		"/lib/i386-linux-gnu/sse2",
		"/lib/i386-linux-gnu/i686",
		"/lib/i386-linux-gnu/i686/sse2",
		fmt.Sprintf("%s/ubuntu12_32/steam-runtime/i386/lib/i386-linux-gnu", steamPath),
		fmt.Sprintf("%s/ubuntu12_32/steam-runtime/i386/lib", steamPath),
		fmt.Sprintf("%s/ubuntu12_32/steam-runtime/i386/usr/lib/i386-linux-gnu", steamPath),
		fmt.Sprintf("%s/ubuntu12_32/steam-runtime/i386/usr/lib", steamPath),
		fmt.Sprintf("%s/ubuntu12_32/steam-runtime/amd64/lib/x86_64-linux-gnu", steamPath),
		fmt.Sprintf("%s/ubuntu12_32/steam-runtime/amd64/lib", steamPath),
		fmt.Sprintf("%s/ubuntu12_32/steam-runtime/amd64/usr/lib/x86_64-linux-gnu", steamPath),
		fmt.Sprintf("%s/ubuntu12_32/steam-runtime/amd64/usr/lib", steamPath),
	}, ":")

	return append(os.Environ(), []string{
		"PATH=" + steamCustomPaths,
		"TERM=xterm",
		"WINEDEBUG=-all",
		fmt.Sprintf("WINEDLLPATH=%s/dist/lib64/wine:%s/dist/lib/wine", protonPath, protonPath),
		"LD_LIBRARY_PATH=" + ldLibraryPath,
		fmt.Sprintf("WINEPREFIX=%s/steamapps/compatdata/%d/pfx/", steamPath, appID),
		"WINEESYNC=1",
		fmt.Sprintf("SteamGameId=%d", appID),
		fmt.Sprintf("SteamAppId=%d", appID),
		"WINEDLLOVERRIDES=steam.exe=b;mfplay=n;d3d11=n;d3d10=n;d3d10core=n;d3d10_1=n;dxgi=n",
		fmt.Sprintf("STEAM_COMPAT_CLIENT_INSTALL_PATH=%s", steamPath),
	}...)
}

func getHomeworldRMPath() string {
	installPath := ""
	for _, path := range defaultInstallPaths {
		resolvedPath := os.ExpandEnv(path)
		slog.Debug("Checking " + resolvedPath + " for HomeworldRM")
		if _, err := os.Stat(resolvedPath); err == nil {
			installPath = resolvedPath
			break
		}
	}
	if installPath == "" {
		help()
		os.Exit(1)
	}
	return installPath
}

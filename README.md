# Homeworld Remastered Proton Launcher

This is a launcher application for the Steam version of Homeworld Remastered Collection. It can be used as a replacement for Gearbox's launcher which has Linux compatibility issues.

Much thanks needs to be given to the version of the launcher at https://git.sr.ht/~_dev_fra/hwrc-proton-launcher which this code was originally a Go port of.

## Running

To get the application, download the binary from the Releases tab for your platform.

The application may be run in two modes:
* Command Line: Takes all inputs from the CLI. Run `./hwrc-proton-launcher -help` for more information.
* GUI: A graphical interface for launching the games. Will display if no CLI options are passed by running `./hwrc-proton-launcher`

## Known Issues

I have only tested this on my machine; there are likely incompatibility problems on other variants. If you encounter an issue file it to this repository!

I have also not yet tested the mod loading functionality.

## Attribution

This project is a port of https://git.sr.ht/~_dev_fra/hwrc-proton-launcher, converting it to a Go application and adding a UI...admittedly it an excuse for me to try out [Fyne](https://fyne.io/) by helping out the experience of running my favorite game :)

## Contributing

The tool has been written in Golang using the Fyne UI framework. The following items are required to be installed on your system in order to compile:

* [Go 1.13+](https://golang.org/dl/)
* The [Fyne Prerequisites](https://fyne.io/develop/index.html) such as GCC and graphics library headers.

Once installed, simply run the following to build and execute the application:

```sh
go build && ./hwrc-proton-launcher
```
package main

import (
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	log "github.com/sirupsen/logrus"
)

// ui runs the game with the expectation that configuration is provided by the GUI
func ui(s *steam) {
	app := app.New()

	// Expose the windowed flag to the UI
	windowed := widget.NewCheck("Windowed", func(checked bool) {
		flagWindowed = checked
	})
	windowed.SetChecked(flagWindowed)

	// Expose the modPath flag to the UI
	modPath := widget.NewEntry()
	modPath.SetText(flagModPath)
	modPath.SetPlaceHolder("Path to Mod")

	// Build it!
	w := app.NewWindow("HWRM Launcher")
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Homeworld Remastered Launcher"),
		widget.NewGroup("Classic",
			widget.NewHBox(
				widget.NewButton("Homeworld 1", func() {
					flagModPath = modPath.Text
					if err := s.run(hw1cla); err != nil {
						log.Fatal(err.Error())
					}
				}),
				widget.NewButton("Homeworld 2", func() {
					flagModPath = modPath.Text
					if err := s.run(hw2cla); err != nil {
						log.Fatal(err.Error())
					}
				}),
			),
		),
		widget.NewGroup("Remastered",
			widget.NewHBox(
				widget.NewButton("Homeworld 1", func() {
					flagModPath = modPath.Text
					if err := s.run(hw1rem); err != nil {
						log.Fatal(err.Error())
					}
				}),
				widget.NewButton("Homeworld 2", func() {
					flagModPath = modPath.Text
					if err := s.run(hw2rem); err != nil {
						log.Fatal(err.Error())
					}
				}),
			),
		),
		widget.NewGroup("Multiplayer",
			widget.NewButton("Homeworld Remastered", func() {
				flagModPath = modPath.Text
				if err := s.run(hwmp); err != nil {
					log.Fatal(err.Error())
				}
			}),
		),
		widget.NewGroup("Options",
			windowed,
			widget.NewHBox(
				widget.NewLabel("Mod Path"),
				modPath,
			),
		),
	))

	w.ShowAndRun()
}

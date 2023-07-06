package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/exp/slog"
)

// ui runs the game with the expectation that configuration is provided by the GUI
func ui(s *steam) error {
	app := app.New()

	// Expose the modPath flag to the UI
	modPath := widget.NewEntryWithData(binding.BindString(&flagModPath))
	modPath.SetPlaceHolder("Path to Mod")

	// Build it!
	w := app.NewWindow("HWRM Launcher")

	content := container.NewVBox(
		widget.NewLabel("Homeworld Remastered Launcher"),
		widget.NewCard("Classic", "",
			container.NewHBox(
				widget.NewButton("Homeworld 1", func() {
					flagModPath = modPath.Text
					if err := s.run(hw1cla); err != nil {
						slog.Error(err.Error())
					}
				}),
				widget.NewButton("Homeworld 2", func() {
					flagModPath = modPath.Text
					if err := s.run(hw2cla); err != nil {
						slog.Error(err.Error())
					}
				}),
			),
		),
		widget.NewCard("Remastered", "",
			container.NewHBox(
				widget.NewButton("Homeworld 1", func() {
					flagModPath = modPath.Text
					if err := s.run(hw1rem); err != nil {
						slog.Error(err.Error())
					}
				}),
				widget.NewButton("Homeworld 2", func() {
					flagModPath = modPath.Text
					if err := s.run(hw2rem); err != nil {
						slog.Error(err.Error())
					}
				}),
			),
		),
		widget.NewCard("Multiplayer", "",
			widget.NewButton("Homeworld Remastered", func() {
				flagModPath = modPath.Text
				if err := s.run(hwmp); err != nil {
					slog.Error(err.Error())
				}
			}),
		),
		widget.NewCard("Options", "",
			container.NewVBox(
				widget.NewCheckWithData("Windowed", binding.BindBool(&flagWindowed)),
				widget.NewForm(
					widget.NewFormItem("Mod Path", modPath),
				),
			),
		),
	)

	w.SetContent(content)
	w.ShowAndRun()
	return nil
}

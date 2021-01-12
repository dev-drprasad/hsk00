package main

import (
	"log"

	"github.com/dev-drprasad/hsk00/pkg"
	"github.com/leaanthony/mewn"
	"github.com/ncruces/zenity"
	"github.com/wailsapp/wails"
)

type Runtime struct {
	runtime *wails.Runtime
}

// WailsInit initialize wails
func (r *Runtime) WailsInit(wr *wails.Runtime) error {
	r.runtime = wr
	return nil
}

func (r *Runtime) SelectGames() []string {
	files, _ := zenity.SelectFileMutiple(zenity.Filename(""), zenity.FileFilters{{"NES ROMs", []string{"*.nes"}}})
	log.Println("files", files)
	return files
}

func (r *Runtime) SelectRootDir() string {
	file, _ := zenity.SelectFile(zenity.Filename(""), zenity.Directory())
	return file
}
func (r *Runtime) AddGames(rootDir string, categoryID int, newGames []string) error {
	return pkg.Add(rootDir, categoryID, newGames, "")
}

func main() {

	js := mewn.String("./frontend/build/static/js/main.js")
	css := mewn.String("./frontend/build/static/css/main.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:  420,
		Height: 520,
		Title:  "hsk00",
		JS:     js,
		CSS:    css,
		Colour: "#0d1117",
	})

	r := &Runtime{}

	app.Bind(r)
	app.Run()
}

package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/dev-drprasad/hsk00/pkg"
	"github.com/dev-drprasad/qt/widgets"
)

const margin = 16
const Qt__AlignHCenter = 4
const Qt__AlignCenter = 132

func newGroupCenter(groupName string, ws ...widgets.QWidget_ITF) *widgets.QGroupBox {
	rootPathGroup := widgets.NewQGroupBox2(groupName, nil)
	layout := widgets.NewQHBoxLayout2(nil)
	rootPathGroup.SetLayout(layout)

	for _, w := range ws {
		layout.AddWidget(w, 0, Qt__AlignHCenter)
	}

	return rootPathGroup
}

func newGroup(groupName string, ws ...widgets.QWidget_ITF) *widgets.QGroupBox {
	rootPathGroup := widgets.NewQGroupBox2(groupName, nil)
	layout := widgets.NewQHBoxLayout2(nil)
	rootPathGroup.SetLayout(layout)

	for _, w := range ws {
		layout.AddWidget(w, 0, 0)
	}
	return rootPathGroup
}

func getHomePath() (homePath string) {
	if runtime.GOOS == "windows" {
		homePath = os.Getenv("USERPROFILE")
	} else {
		homePath = os.Getenv("HOME")
	}
	return
}

var categories = []string{"0. Action Games", "1. Shoot Games", "2. Sport Games", "3. Fight Games", "4. Racing Games", "5. Puzzle Games"}

func main() {

	// needs to be called once before you can start using the QWidgets
	app := widgets.NewQApplication(len(os.Args), os.Args)

	// create a window
	// and sets the title to "Hello Widgets Example"
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(450, 400)
	window.SetWindowTitle("hsk00")

	// create a regular widget
	// give it a QVBoxLayout
	// and make it the central widget of the window
	layout := widgets.NewQVBoxLayout()
	layout.SetContentsMargins(margin, margin, margin, margin)
	container := widgets.NewQWidget(nil, 0)
	container.SetLayout(layout)
	window.SetCentralWidget(container)

	backupWarnWidget := widgets.NewQLabel2("", nil, 0)
	backupWarnWidget.SetText("🚨 Proceed after taking SD card backup 🚨")
	layout.AddWidget(backupWarnWidget, 0, Qt__AlignCenter)

	rootPathInput := widgets.NewQLineEdit(nil)
	rootPathInput.SetPlaceholderText("Root Path (SD Card)...")
	rootPathInput.SizePolicy().SetVerticalPolicy(widgets.QSizePolicy__Preferred)

	rootDir := ""
	selectRootButton := widgets.NewQPushButton2("Select Root", nil)
	selectRootButton.ConnectClicked(func(bool) {
		fileDialog := widgets.NewQFileDialog(nil, 0)
		fileDialog.SetFileMode(widgets.QFileDialog__DirectoryOnly)
		if fileDialog.Exec() == int(widgets.QDialog__Accepted) {
			if len(fileDialog.SelectedFiles()) == 0 {
				return
			}
			selectedDir := fileDialog.SelectedFiles()[0]
			rootPathInput.SetText(selectedDir)
			rootDir = selectedDir
		}
	})

	selectRootButton.SizePolicy().SetVerticalPolicy(widgets.QSizePolicy__Preferred)
	container.SizePolicy().SetVerticalPolicy(widgets.QSizePolicy__Preferred)

	rootPathGroup := newGroup("", rootPathInput, selectRootButton)
	layout.AddWidget(rootPathGroup, 0, 0)

	categoryID := -1
	categoryDropdown := widgets.NewQComboBox(nil)
	categoryDropdown.AddItems(append([]string{"Select Category"}, categories...))
	categoryDropdown.ConnectCurrentIndexChanged(func(i int) {
		categoryID = i - 1

	})
	categoryGroup := newGroupCenter("", categoryDropdown)
	layout.AddWidget(categoryGroup, 0, 0)

	messageWidget := widgets.NewQLabel2("", nil, 0)

	newGamePaths := []string{}
	selectGamesButton := widgets.NewQPushButton2("Select Games", nil)
	selectGamesButton.ConnectClicked(func(bool) {
		fileDialog := widgets.NewQFileDialog2(nil, "Select Games", getHomePath(), "NES ROMs (*.nes)")
		selectedFiles := fileDialog.GetOpenFileNames(window, "Select Game(s)", getHomePath(), "NES ROMs (*.nes)", "", 0)
		if len(selectedFiles) == 0 {
			log.Println("select at least one game")
			return
		}
		newGamePaths = append(newGamePaths, selectedFiles...)
		messageWidget.SetText(fmt.Sprintf("Selected %d games", len(newGamePaths)))
	})

	selectGamesGroup := newGroupCenter("", selectGamesButton)
	layout.AddWidget(selectGamesGroup, 0, 0)
	layout.AddWidget(messageWidget, 0, Qt__AlignCenter)

	addGamesButton := widgets.NewQPushButton2("Add Games", nil)
	addGamesButton.ConnectClicked(func(bool) {
		log.Println(rootDir, categoryID, newGamePaths)
		if rootDir == "" {
			messageWidget.SetText("Error: Select root directory")
			return
		}
		if categoryID < 0 {
			messageWidget.SetText("Error: Select category")
			return
		}
		if len(newGamePaths) == 0 {
			messageWidget.SetText("Error: Select at least one game")
			return
		}

		if err := pkg.Add(rootDir, categoryID, newGamePaths, "", ""); err != nil {
			messageWidget.SetText(fmt.Sprintf("Error: %s", err.Error()))
		} else {
			messageWidget.SetText("✅ Done")
		}
	})

	layout.AddWidget(addGamesButton, 1, Qt__AlignCenter)
	layout.SetAlignment2(layout, Qt__AlignCenter)

	// make the window visible
	window.Show()

	// start the main Qt event loop
	// and block until app.Exit() is called
	// or the window is closed by the user
	app.Exec()
}

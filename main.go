package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
	"os"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("File Drop App")

	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	// Создаем виджет для отображения пути к файлу
	filePathLabel := widget.NewHyperlink("Drop a file here", nil)
	copyIconButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		clipboard.Write(clipboard.FmtText, []byte(filePathLabel.URL.String()))
	})

	// Создаем дроп-зону для файла
	dropContainer := container.New(
		layout.NewVBoxLayout(),
		container.NewVBox(
			widget.NewLabelWithStyle("Drop File Here", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewIcon(theme.DocumentIcon()),
		),
		container.NewHBox(
			filePathLabel,
			copyIconButton,
		),
	)

	// Если передан аргумент командной строки, используем его как путь к файлу
	if len(os.Args) > 1 {
		filePath := os.Args[1]
		filePathLabel.SetText(fmt.Sprintf("File path from command line argument: %s", filePath))
		err := filePathLabel.SetURLFromString("https://inxo.ru")
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
	}

	myWindow.SetOnDropped(func(pos fyne.Position, url []fyne.URI) {
		filePath := url[0].String()
		filePathLabel.SetText(fmt.Sprintf("Dropped file path: %s", filePath))
		err := filePathLabel.SetURLFromString("https://inxo.ru")
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
	})

	myWindow.SetContent(dropContainer)
	myWindow.ShowAndRun()
}

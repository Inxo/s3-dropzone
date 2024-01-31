package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("File Drop App")

	// Создаем виджет для отображения пути к файлу
	filePathLabel := widget.NewLabel("Drop a file here")

	// Создаем дроп-зону для файла
	dropContainer := container.New(
		layout.NewVBoxLayout(),
		container.NewPadded(
			widget.NewLabelWithStyle("Drop File Here", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewIcon(theme.DocumentIcon()),
		),
		filePathLabel,
	)

	//myWindow.SetOnTypedRune(func(ev *fyne.TypedRuneEvent) {
	//	Пустая функция, чтобы окно реагировало на события клавиатуры
	//})

	myWindow.SetOnDropped(func(pos fyne.Position, url []fyne.URI) {
		//filePath, err := handleDropEvent(dropData)
		//if err != nil {
		//	dialog.ShowError(err, myWindow)
		//	return
		//}
		filePath := url[0].String()
		filePathLabel.SetText(fmt.Sprintf("Dropped file path: %s", filePath))
	})

	myWindow.SetContent(dropContainer)
	myWindow.ShowAndRun()
}

//func handleDropEvent(dropData fyne.DropEvent) (string, error) {
//	if dropData == nil {
//		return "", fmt.Errorf("invalid drop event")
//	}
//
//	filePath := dropData.URI().Path()
//	// Выводим путь к файлу
//	log.Printf("File dropped: %s", filePath)
//
//	return filePath, nil
//}

package main

import (
	"dropZone/short"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
	"os"
	"strings"
	"time"
)

func main() {
	myApp := app.NewWithID("ru.inxo.drop.app")
	myApp.SetIcon(resourceIconPng)
	wd := myApp.Storage().RootURI().String()

	err := os.MkdirAll(wd, os.ModePerm)
	if err != nil {
		println(err)
	}
	// Load environment variables
	//err = godotenv.Load(wd + "/.env")

	// S3 Settings Form
	bucketEntry := widget.NewEntry()
	bucketName := myApp.Preferences().String("BUCKET_NAME")
	if len(bucketName) == 0 {
		// create default
		myApp.Preferences().SetString("AWS_ACCESS_KEY_ID", "tw6kCfjGMh75do0R9I6SAG8JyvuuKI80")
		myApp.Preferences().SetString("BUCKET_NAME", "nulljet-share")
		myApp.Preferences().SetString("AWS_ENDPOINT", "https://tw-001.s3.synologyc2.net")
		myApp.Preferences().SetString("AWS_SECRET_ACCESS_KEY", "NLegC6KwfHf4zftDaWpSxnMiNhVv9KZF")
		myApp.Preferences().SetString("AWS_REGION", "tw-001")
		bucketName = myApp.Preferences().String("BUCKET_NAME")
	}
	bucketEntry.SetText(bucketName)
	regionEntry := widget.NewEntry()
	regionEntry.SetText(myApp.Preferences().String("AWS_REGION"))
	idEntry := widget.NewEntry()
	idEntry.SetText(myApp.Preferences().String("AWS_ACCESS_KEY_ID"))
	tokenEntry := widget.NewEntry()
	tokenEntry.SetText(myApp.Preferences().String("AWS_SECRET_ACCESS_KEY"))
	endpointEntry := widget.NewEntry()
	endpointEntry.SetText(myApp.Preferences().String("AWS_ENDPOINT"))

	progressEntry := widget.NewProgressBarInfinite()
	progressEntry.Hide()

	myWindow := myApp.NewWindow("File Drop App")

	saveShort := func() {
		// Handle form submission
		bucket := bucketEntry.Text
		region := regionEntry.Text
		token := tokenEntry.Text
		id := idEntry.Text
		endpoint := endpointEntry.Text

		// Perform save data
		saveData(myApp, myWindow, bucket, endpoint, region, id, token)
	}

	if desk, ok := myApp.(desktop.App); ok {
		m := fyne.NewMenu("MyApp",
			fyne.NewMenuItem("Show", func() {
				myWindow.Show()
			}))

		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(theme.StorageIcon())
	}

	myWindow.Resize(fyne.Size{
		Width:  800,
		Height: 600,
	})

	myWindow.SetCloseIntercept(func() {
		myWindow.Hide()
	})
	err = clipboard.Init()
	if err != nil {
		panic(err)
	}

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Bucket", Widget: bucketEntry},
			{Text: "Endpoint", Widget: endpointEntry},
			{Text: "Region", Widget: regionEntry},
			{Text: "Id", Widget: idEntry},
			{Text: "Token", Widget: tokenEntry},
		},
		OnSubmit:   saveShort,
		SubmitText: "Save",
	}

	// Создаем виджет для отображения пути к файлу
	filePathLabel := widget.NewHyperlink("Drop a file here", nil)
	copyIconButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		clipboard.Write(clipboard.FmtText, []byte(filePathLabel.URL.String()))
	})

	myWindow.SetIcon(theme.StorageIcon())

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
		err := filePathLabel.SetURLFromString("https://s.inxo.ru")
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		//filePathLabel.Hide()
	}

	headerSettings := widget.NewLabel("S3 Object Storage Settings")
	// Combine forms into a tab container
	tabs := container.NewAppTabs(
		container.NewTabItem("Upload", container.New(layout.NewVBoxLayout(), dropContainer)),
		container.NewTabItem("Settings", container.New(layout.NewVBoxLayout(), headerSettings, form)),
	)

	myWindow.SetContent(container.NewVBox(
		container.NewHBox(widget.NewLabel("Private share app")),
		tabs,
	))

	sync := Sync{}

	myWindow.SetOnDropped(func(pos fyne.Position, url []fyne.URI) {
		filePath := url[0].String()
		if len(url) > 1 {
			// generate page
		}
		filePathLabel.SetText(fmt.Sprintf("Dropped file path: %s", filePath))

		expireIn := os.Getenv("EXPIRE_IN")
		if len(expireIn) == 0 {
			expireIn = "+168h"
		}
		if !strings.HasPrefix(expireIn, "+") {
			expireIn = "+" + expireIn
		}

		duration, err := time.ParseDuration(expireIn)
		if err != nil {
			return
		}

		urlUploaded, err := sync.UploadToS3(filePath, duration)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		urlShort, err := short.NewLink(urlUploaded)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		err = filePathLabel.SetURLFromString(urlShort)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		filePathLabel.SetText("Download link")
	})
	myWindow.ShowAndRun()
}

func saveData(myApp fyne.App, myWindow fyne.Window, bucket string, endpoint string, region string, id string, token string) {
	// Call os.MkdirAll with the directory path
	expireIn := "24h"
	shortService := "https://s.inxo.ru/shorten"

	myApp.Preferences().SetString("BUCKET_NAME", bucket)
	myApp.Preferences().SetString("AWS_ACCESS_KEY_ID", id)
	myApp.Preferences().SetString("AWS_ENDPOINT", endpoint)
	myApp.Preferences().SetString("AWS_SECRET_ACCESS_KEY", token)
	myApp.Preferences().SetString("AWS_REGION", region)
	myApp.Preferences().SetString("EXPIRE_IN", expireIn)
	myApp.Preferences().SetString("SHORTEN_SERVICE", shortService)

	dialog.ShowInformation("Save Success", "Data save successfully!", myWindow)
}

package main

import (
	"bufio"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/joho/godotenv"
	"golang.design/x/clipboard"
	"log"
	"os"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln("No work dir")
	}
	// Load environment variables
	err = godotenv.Load(wd + "/.env")

	myApp := app.New()
	// S3 Settings Form
	bucketEntry := widget.NewEntry()
	bucketName := os.Getenv("BUCKET_NAME")
	bucketEntry.SetText(bucketName)
	regionEntry := widget.NewEntry()
	regionEntry.SetText(os.Getenv("AWS_REGION"))
	idEntry := widget.NewEntry()
	idEntry.SetText(os.Getenv("AWS_ACCESS_KEY_ID"))
	tokenEntry := widget.NewEntry()
	tokenEntry.SetText(os.Getenv("AWS_SECRET_ACCESS_KEY"))
	endpointEntry := widget.NewEntry()
	endpointEntry.SetText(os.Getenv("AWS_ENDPOINT"))

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
		saveData(myWindow, bucket, endpoint, region, id, token, wd)
	}

	if desk, ok := myApp.(desktop.App); ok {
		m := fyne.NewMenu("MyApp",
			fyne.NewMenuItem("Show", func() {
				myWindow.Show()
			}))
		desk.SetSystemTrayMenu(m)
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

	myWindow.SetIcon(theme.AccountIcon())

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
	err = sync.Init()
	if err != nil {
		dialog.ShowError(err, myWindow)
	}

	myWindow.SetOnDropped(func(pos fyne.Position, url []fyne.URI) {
		filePath := url[0].String()
		if len(url) > 1 {
			// generate page
		}
		filePathLabel.SetText(fmt.Sprintf("Dropped file path: %s", filePath))

		urlUploaded, err := sync.UploadToS3(filePath)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		err = filePathLabel.SetURLFromString(urlUploaded)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		filePathLabel.SetText("Download link")
	})
	myWindow.ShowAndRun()
}

func saveData(myWindow fyne.Window, bucket string, endpoint string, region string, id string, token string, wd string) {
	file, err := os.Create(wd + "/.env")
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	writer := bufio.NewWriter(file)
	_, _ = fmt.Fprintf(writer, "BUCKET_NAME=%s\n", bucket)
	_, _ = fmt.Fprintf(writer, "AWS_ACCESS_KEY_ID=%s\n", id)
	_, _ = fmt.Fprintf(writer, "AWS_ENDPOINT=%s\n", endpoint)
	_, _ = fmt.Fprintf(writer, "AWS_SECRET_ACCESS_KEY=%s\n", token)
	_, _ = fmt.Fprintf(writer, "AWS_REGION=%s\n", region)
	err = writer.Flush()
	if err != nil {
		dialog.ShowError(err, myWindow)
	} else {
		dialog.ShowInformation("Save Success", "Data save successfully!", myWindow)
	}
}

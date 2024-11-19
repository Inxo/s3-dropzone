package main

import (
	"capyDrop/page_maker"
	"capyDrop/short"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/kbinani/screenshot"
	"golang.design/x/clipboard"
	"image/png"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	myApp := app.NewWithID("com.n32b.drop.app")
	myApp.SetIcon(resourceIconPng)
	wd := myApp.Storage().RootURI().String()

	err := os.MkdirAll(wd, os.ModePerm)
	if err != nil {
		println(err)
	}

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
	region := myApp.Preferences().String("AWS_REGION")
	keyId := myApp.Preferences().String("AWS_ACCESS_KEY_ID")
	accessKey := myApp.Preferences().String("AWS_SECRET_ACCESS_KEY")
	endpoint := myApp.Preferences().String("AWS_ENDPOINT")

	bucketEntry.SetText(bucketName)
	regionEntry := widget.NewEntry()
	regionEntry.SetText(region)
	idEntry := widget.NewEntry()
	idEntry.SetText(keyId)
	tokenEntry := widget.NewEntry()
	tokenEntry.SetText(accessKey)
	endpointEntry := widget.NewEntry()
	endpointEntry.SetText(endpoint)

	progressEntry := widget.NewProgressBarInfinite()
	progressEntry.Hide()

	myWindow := myApp.NewWindow("File Drop App")
	myWindow.SetIcon(resourceIconPng)
	saveScreen := func() string {
		log.Println("Take Screenshot")
		myWindow.Hide()
		time.Sleep(1)
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		fmt.Println(exPath)
		n := screenshot.NumActiveDisplays()
		var fileName string
		for i := 0; i < n; i++ {
			bounds := screenshot.GetDisplayBounds(i)

			img, err := screenshot.CaptureRect(bounds)
			if err != nil {
				panic(err)
			}

			fileName = fmt.Sprintf("%d_%dx%d.png", i, bounds.Dx(), bounds.Dy())
			file, _ := os.Create(exPath + "/" + fileName)
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {

				}
			}(file)
			err = png.Encode(file, img)
			if err != nil {
				return ""
			}
		}

		myWindow.Show()
		return exPath + "/" + fileName
	}

	savePreferences := func() {
		// Handle form submission
		bucket := bucketEntry.Text
		region := regionEntry.Text
		token := tokenEntry.Text
		id := idEntry.Text
		endpoint := endpointEntry.Text

		// Perform save data
		saveData(myApp, myWindow, bucket, endpoint, region, id, token)
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
		OnSubmit:   savePreferences,
		SubmitText: "Save",
	}

	// Создаем виджет для отображения пути к файлу
	uri, err := url.Parse("https://s.inxo.ru")
	if err != nil {
		dialog.ShowError(err, myWindow)
	}
	filePathLabel := widget.NewHyperlink("File: ", uri)
	copyIconButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		if len(filePathLabel.URL.String()) > 0 {
			clipboard.Write(clipboard.FmtText, []byte(filePathLabel.URL.String()))
		}
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

	pageMaker := page_maker.PageMaker{}
	sync := Sync{PageMaker: pageMaker}
	err = sync.Init(bucketName, keyId, accessKey, endpoint, region)
	if err != nil {
		dialog.ShowError(err, myWindow)
	}

	shortUpload := func(filePath string) (string, error) {
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
			return "", err
		}

		urlUploaded, err := sync.UploadToS3(filePath, duration)
		if err != nil {
			return "", err
		}
		urlShort, err := short.NewLink(urlUploaded)
		if err != nil {
			return "", err
		}
		err = filePathLabel.SetURLFromString(urlShort)
		if err != nil {
			return "", err
		}
		filePathLabel.SetText("Download link")
		return urlUploaded, nil
	}

	// Если передан аргумент командной строки, используем его как путь к файлу
	if len(os.Args) > 1 {
		filePath := os.Args[1]
		filePathLabel.SetText(fmt.Sprintf("File path from command line argument: %s", filePath))
		err := filePathLabel.SetURLFromString("https://s.inxo.ru")
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		uploadLink, _ := shortUpload(filePath)
		pagePath := pageMaker.Do(uploadLink, "image")
		pageLink, err := shortUpload(pagePath)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		uploadLink = pageLink
	}

	makeScreen := func() {
		filePath := saveScreen()
		if len(filePath) > 1 {
			uploadLink, _ := shortUpload(filePath)
			pagePath := pageMaker.Do(uploadLink, "image")
			pageLink, err := shortUpload(pagePath)
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			uploadLink = pageLink
		}
	}

	ctrlAltS := desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl | fyne.KeyModifierAlt}
	myWindow.Canvas().AddShortcut(&ctrlAltS, func(shortcut fyne.Shortcut) {
		makeScreen()
	})
	myWindow.SetIcon(resourceIconPng)

	myWindow.SetOnDropped(func(pos fyne.Position, url []fyne.URI) {
		filePath := url[0].String()
		if len(url) > 1 {
			// generate page
		}
		_, err := shortUpload(filePath)
		dialog.ShowError(err, myWindow)
		//if makeStyledPage {
		//	mime, err := detectMime(filePath)
		//	dialog.ShowError(err, myWindow)
		//	pagePath := pageMaker.Do(link, mime)
		//	_, err = shortUpload(pagePath)
		//}
	})

	if desk, ok := myApp.(desktop.App); ok {
		m := fyne.NewMenu("MyApp",
			fyne.NewMenuItem("Show App", func() {
				myWindow.Show()
			}),
			fyne.NewMenuItem("Take Screenshot", func() {
				makeScreen()
			}))

		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(resourceIconPng)
	}

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

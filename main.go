package main

import (
	"capyDrop/capywidget"
	"capyDrop/page_maker"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
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
	"strconv"
	"time"
)

func main() {
	myApp := app.NewWithID("com.n32b.drop.app")
	myApp.SetIcon(resourceIconPng)
	wd := myApp.Storage().RootURI().String()

	tray := NewTrayIcon(myApp.(desktop.App))

	err := os.MkdirAll(wd, os.ModePerm)
	if err != nil {
		println(err)
	}

	// S3 Settings Form
	bucketEntry := widget.NewEntry()
	shortService := myApp.Preferences().String("SHORT_SERVICE")
	if len(shortService) == 0 {
		// create default
		myApp.Preferences().SetInt("EXPIRE_IN", 24)
		myApp.Preferences().SetString("SHORT_SERVICE", "https://s.n32b.com/shorten")
	}
	bucketName := myApp.Preferences().String("BUCKET_NAME")
	region := myApp.Preferences().String("AWS_REGION")
	keyId := myApp.Preferences().String("AWS_ACCESS_KEY_ID")
	accessKey := myApp.Preferences().String("AWS_SECRET_ACCESS_KEY")
	endpoint := myApp.Preferences().String("AWS_ENDPOINT")
	expireIn := myApp.Preferences().Int("EXPIRE_IN")
	shortService = myApp.Preferences().String("SHORT_SERVICE")

	bucketEntry.SetText(bucketName)
	regionEntry := widget.NewEntry()
	regionEntry.SetText(region)
	idEntry := widget.NewEntry()
	idEntry.Password = true
	idEntry.SetText(keyId)
	tokenEntry := widget.NewEntry()
	tokenEntry.Password = true
	tokenEntry.SetText(accessKey)
	endpointEntry := widget.NewEntry()
	endpointEntry.SetText(endpoint)
	expireInEntry := capywidget.NewNumericalEntry()
	expireInEntry.SetValue(expireIn)
	shortServiceEntry := widget.NewEntry()
	shortServiceEntry.SetText(shortService)

	progressEntry := widget.NewProgressBarInfinite()
	progressEntry.Hide()

	myWindow := myApp.NewWindow("File Drop App")
	myWindow.SetIcon(resourceIconPng)

	savePreferences := func() {
		// Handle form submission
		bucket := bucketEntry.Text
		region := regionEntry.Text
		token := tokenEntry.Text
		id := idEntry.Text
		endpoint := endpointEntry.Text
		expireIn, _ := strconv.Atoi(expireInEntry.Text)
		shortService := shortServiceEntry.Text

		// Perform save data
		p := Preferences{
			Bucket:       bucket,
			Endpoint:     endpoint,
			Region:       region,
			Id:           id,
			Token:        token,
			ExpireIn:     expireIn,
			ShortService: shortService,
		}

		//systray.SetTitle("Awesome App")
		//systray.SetTooltip("Pretty awesome棒棒嗒")
		p.SaveData(myApp, myWindow)
	}

	pageMaker := page_maker.PageMaker{}

	sync := Sync{PageMaker: pageMaker}
	err = sync.Init(bucketName, keyId, accessKey, endpoint, region)
	if err != nil {
		dialog.ShowError(err, myWindow)
	}

	myWindow.Resize(fyne.Size{
		Width:  800,
		Height: 600,
	})

	myWindow.SetCloseIntercept(func() {
		myWindow.Hide()
		setActivationPolicy(false)
	})
	err = clipboard.Init()
	if err != nil {
		panic(err)
	}

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Widget: widget.NewLabel("Object Storage Settings")},
			{Text: "Bucket", Widget: bucketEntry},
			{Text: "Endpoint", Widget: endpointEntry},
			{Text: "Region", Widget: regionEntry},
			{Text: "Id", Widget: idEntry},
			{Text: "Token", Widget: tokenEntry},
			{Widget: widget.NewLabel("Upload Settings")},
			{Text: "Expire In", Widget: expireInEntry},
			{Text: "Short Service", Widget: shortServiceEntry},
		},
		OnSubmit:   savePreferences,
		SubmitText: "Save",
	}

	// Создаем виджет для отображения пути к файлу
	uri, err := url.Parse("https://s.n32b.com")
	if err != nil {
		dialog.ShowError(err, myWindow)
	}
	filePathLabel := widget.NewHyperlink("n32b.com ", uri)
	copyIconButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		if len(filePathLabel.URL.String()) > 0 {
			clipboard.Write(clipboard.FmtText, []byte(filePathLabel.URL.String()))
		}
	})

	// Создаем дроп-зону для файла
	image := canvas.NewImageFromResource(resourceIconPng)
	image.FillMode = canvas.ImageFillOriginal
	openImage := func() {
		fd := dialog.NewFileOpen(func(closer fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			if closer == nil {
				return
			}
			_, err = Upload(closer.URI().String(), filePathLabel, sync, myApp, &tray)
			//dialog.ShowInformation("File", closer.URI().String(), myWindow)
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
		}, myWindow)
		if err != nil {
			dialog.ShowError(err, myWindow)
		}
		fd.Show()
	}
	btn := widget.NewButton("", openImage)
	box := container.NewPadded(btn, image)
	dropContainer := container.New(
		layout.NewVBoxLayout(),
		container.NewHBox(
			filePathLabel,
			copyIconButton,
		),
		container.NewVBox(
			widget.NewLabelWithStyle("Drop File Here", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			box,
		),
	)

	// Combine forms into a tab container
	tabs := container.NewAppTabs(
		container.NewTabItem("Upload", container.New(layout.NewVBoxLayout(), dropContainer)),
		//container.NewTabItem("Latest uploads", container.New(layout.NewVBoxLayout(), dropContainer)),
		container.NewTabItem("Settings", container.New(layout.NewVBoxLayout(), form)),
	)

	myWindow.SetContent(container.NewVBox(
		container.NewHBox(widget.NewLabel("Private share app")),
		tabs,
	))

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
			//_ = file.Close()
		}

		myWindow.Show()
		tabs.SelectIndex(0)
		return exPath + "/" + fileName
	}

	// Если передан аргумент командной строки, используем его как путь к файлу
	if len(os.Args) > 1 {
		filePath := os.Args[1]
		filePathLabel.SetText(fmt.Sprintf("File path from command line argument: %s", filePath))
		err := filePathLabel.SetURLFromString("https://s.n32b.com")
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		uploadLink, _ := Upload(filePath, filePathLabel, sync, myApp, &tray)
		pagePath := pageMaker.Do(uploadLink, "image")
		pageLink, err := Upload(pagePath, filePathLabel, sync, myApp, &tray)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		uploadLink = pageLink
	}

	makeScreen := func() {
		filePath := saveScreen()
		if len(filePath) > 1 {
			_, err = Upload(filePath, filePathLabel, sync, myApp, &tray)
			//pagePath := pageMaker.Do(uploadLink, "image")
			//pageLink, err := shortUpload(pagePath)
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			//uploadLink = pageLink
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
		_, err := Upload(filePath, filePathLabel, sync, myApp, &tray)
		if err != nil {
			dialog.ShowError(err, myWindow)
		}
		//
		//if makeStyledPage {
		//	mime, err := detectMime(filePath)
		//	dialog.ShowError(err, myWindow)
		//	pagePath := pageMaker.Do(link, mime)
		//	_, err = shortUpload(pagePath)
		//}
	})

	if desk, ok := myApp.(desktop.App); ok {
		m := fyne.NewMenu("MyApp",
			fyne.NewMenuItem("Drop File", func() {
				myWindow.Show()
				setActivationPolicy(true)
				tabs.SelectIndex(0)
			}),
			fyne.NewMenuItem("Take Screenshot", func() {
				makeScreen()
			}))

		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(resourceIconPng)
	}
	myApp.Lifecycle().SetOnStarted(func() {
		setActivationPolicy(true)
	})

	myWindow.ShowAndRun()
}

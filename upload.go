package main

import (
	"capyDrop/short"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"os"
	"strings"
	"time"
)

func Upload(filePath string, ui *widget.Hyperlink, sync Sync, myApp fyne.App) (string, error) {
	ui.SetText(fmt.Sprintf("Dropped file path: %s", filePath))

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
	urlShort, err := short.NewLink(urlUploaded, myApp)
	if err != nil {
		return "", err
	}
	err = ui.SetURLFromString(urlShort)
	if err != nil {
		return "", err
	}
	ui.SetText("Download link")
	return urlUploaded, nil
}

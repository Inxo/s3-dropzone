package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

type Preferences struct {
	Bucket       string
	Endpoint     string
	Region       string
	Id           string
	Token        string
	ExpireIn     int
	ShortService string
}

func (p Preferences) SaveData(myApp fyne.App, myWindow fyne.Window) {
	myApp.Preferences().SetString("BUCKET_NAME", p.Bucket)
	myApp.Preferences().SetString("AWS_ACCESS_KEY_ID", p.Id)
	myApp.Preferences().SetString("AWS_ENDPOINT", p.Endpoint)
	myApp.Preferences().SetString("AWS_SECRET_ACCESS_KEY", p.Token)
	myApp.Preferences().SetString("AWS_REGION", p.Region)
	myApp.Preferences().SetInt("EXPIRE_IN", p.ExpireIn)
	myApp.Preferences().SetString("SHORT_SERVICE", p.ShortService)

	dialog.ShowInformation("Save Success", "Data save successfully!", myWindow)
}

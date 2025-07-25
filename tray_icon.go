package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	_ "fyne.io/fyne/v2/theme"
	"time"
)

var iconIndex = 0

type TrayIcon struct {
	tray  desktop.App
	quit  chan bool
	icons []fyne.Resource
}

func NewTrayIcon(tray desktop.App) TrayIcon {
	icons := []fyne.Resource{
		icon1,
		icon2,
		icon3,
		icon4,
		icon5,
		icon6,
		icon7,
	}
	quit := make(chan bool)
	return TrayIcon{
		tray:  tray,
		quit:  quit,
		icons: icons,
	}
}

func (t *TrayIcon) Stop() {
	t.quit <- true
}

func (t *TrayIcon) Start() {
	go func() {
		for {
			select {
			case <-t.quit:
				t.tray.SetSystemTrayIcon(resourceIconPng)
				return
			default:
				time.Sleep(300 * time.Millisecond)
				iconIndex = (iconIndex + 1) % len(t.icons)
				t.tray.SetSystemTrayIcon(t.icons[iconIndex])
			}
		}
	}()
}

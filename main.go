package main

import (
	"crypto/ed25519"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"fyne.io/fyne/v2/widget"

	"github.com/fynelabs/selfupdate"
)

type nomad struct {
	main    fyne.Window
	store   *cityStore
	session *unsplashSession
}

func update() {
	done := make(chan struct{}, 2)

	// Used `selfupdatectl create-keys` followed by `selfupdatectl print-key`
	publicKey := ed25519.PublicKey{178, 103, 83, 57, 61, 138, 18, 249, 244, 80, 163, 162, 24, 251, 190, 241, 11, 168, 179, 41, 245, 27, 166, 70, 220, 254, 118, 169, 101, 26, 199, 129}

	// The public key above match the signature of the below file served by our CDN
	httpSource := selfupdate.NewHTTPSource(nil, "http://geoffrey-test-artefacts.fynelabs.com/nomad.exe")

	config := &selfupdate.Config{
		Source:    httpSource,
		Schedule:  selfupdate.Schedule{FetchOnStart: true, Interval: time.Minute * time.Duration(60)},
		PublicKey: publicKey,

		// This is here to force an update by announcing a time so old that nothing existed
		Current: &selfupdate.Version{Date: time.Unix(100, 0)},

		ProgressCallback: func(f float64, err error) { log.Println("Download", f, "%") },
		RestartConfirmCallback: func() bool {
			done <- struct{}{}
			return true
		},
		UpgradeConfirmCallback: func(_ string) bool { return true },
	}

	_, err := selfupdate.Manage(config)
	if err != nil {
		log.Println("Error while setting up update manager: ", err)
		return
	}

	<-done
}

func main() {

	a := app.NewWithID("com.fynelabs.nomad")
	w := a.NewWindow("Nomad")

	store := newCityStore(a.Preferences())
	session := newUnsplashSession(a.Storage(), store)
	ui := &nomad{main: w, store: store, session: session}

	var _ fyne.Theme = (*myTheme)(nil)
	a.Settings().SetTheme(&myTheme{})

	splash := ui.makeSplash()

	updateButton := widget.NewButton("update", func() { update() })
	w.SetContent(container.NewVBox(updateButton, container.NewMax(ui.makeHome(), splash)))
	w.SetPadded(false)
	w.Resize(fyne.NewSize(300, 500))

	go ui.fadeSplash(splash)
	w.ShowAndRun()

}

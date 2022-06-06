package main

import (
	"crypto/ed25519"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"github.com/fynelabs/fyneselfupdate"
	"github.com/fynelabs/selfupdate"
)

func selfManage(a fyne.App, w fyne.Window) {
	// Used `selfupdatectl create-keys` followed by `selfupdatectl print-key`
	publicKey := ed25519.PublicKey{178, 103, 83, 57, 61, 138, 18, 249, 244, 80, 163, 162, 24, 251, 190, 241, 11, 168, 179, 41, 245, 27, 166, 70, 220, 254, 118, 169, 101, 26, 199, 129}

	// The public key above match the signature of the below file served by our CDN
	httpSource := selfupdate.NewHTTPSource(nil, "http://geoffrey-test-artefacts.fynelabs.com/nomad-{{.OS}}-{{.Arch}}{{.Ext}}")

	config := fyneselfupdate.NewConfigWithTimeout(a, w, time.Duration(1)*time.Minute,
		httpSource,
		selfupdate.Schedule{FetchOnStart: true, Interval: time.Hour * time.Duration(12)},
		publicKey)

	// This is here to force an update by announcing a time so old that nothing existed
	// TODO: remove once enough test are done
	config.Current = &selfupdate.Version{Date: time.Unix(100, 0)}

	_, err := selfupdate.Manage(config)
	if err != nil {
		log.Println("Error while setting up update manager: ", err)
		return
	}
}

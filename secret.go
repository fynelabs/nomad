package main

import (
	"log"
	"os"
)

func secret() string {
	env := os.Getenv("UNSPLASH_ACCESS_KEY")
	if env == "" {
		log.Println("UNSPLASH_ACCESS_KEY is not specified")
	}
	return env
}

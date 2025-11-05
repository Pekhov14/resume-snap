package main

import (
	"log"

	"github.com/Pekhov14/resume-snap/internal/config"
	"github.com/Pekhov14/resume-snap/internal/pdf"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	err = pdf.GeneratePDF(cfg)
	if err != nil {
		log.Fatal("Error with creating pdf:", err)
	}
}

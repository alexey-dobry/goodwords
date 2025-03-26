package main

import (
	"log"
	"os"

	"github.com/alexey-dobry/goodwords/internal/analyser"
	"github.com/alexey-dobry/goodwords/internal/config"
	"github.com/alexey-dobry/goodwords/internal/logger"
)

func main() {
	l, err := logger.InitLogger()
	if err != nil {
		log.Fatal(err)
	}

	l.Info("Successfully initialized logger")

	os.Mkdir("../output", os.ModePerm)

	c, err := config.ReadConfig()
	if err != nil {
		l.Fatal(err)
	}

	l.Info("Successfully read config")

	analyser.RunAnalyser(c, l)

	l.Info("Programm execution complete")
}

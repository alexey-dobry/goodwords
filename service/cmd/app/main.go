package main

import (
	"log"

	"github.com/alexey-dobry/goodwords/internal/config"
	"github.com/alexey-dobry/goodwords/internal/logger"
	"github.com/alexey-dobry/goodwords/internal/service"
)

func main() {
	l, err := logger.InitLogger()
	if err != nil {
		log.Fatal(err)
	}

	l.Info("Successfully initialized logger")

	c, err := config.ReadConfig()
	if err != nil {
		l.Fatal(err)
	}

	l.Info("Successfully read config")

	service.SendRequests(c, l)

	l.Info("Programm execution complete")
}

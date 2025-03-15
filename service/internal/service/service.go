package service

import (
	"sync"
	//"net/http"

	"github.com/alexey-dobry/goodwords/internal/config"
	"github.com/alexey-dobry/goodwords/internal/detector"
	"github.com/alexey-dobry/goodwords/internal/models"
	"go.uber.org/zap"
)

func RequestAndAnalize(wg *sync.WaitGroup, resultChan chan detector.AnalizationResult, rd models.EndpointData, badWords []string) {

}

func SendRequests(c config.Config, l *zap.SugaredLogger) {

}

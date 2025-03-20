package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alexey-dobry/goodwords/internal/config"
	"github.com/alexey-dobry/goodwords/internal/detector"
	"github.com/alexey-dobry/goodwords/internal/models"
	"go.uber.org/zap"
)

type AnalizationResult struct {
	FoundBadWords map[string]string
	URL           string
	DataType      string
	Err           error
}

func FormatAnalizationResult(ar AnalizationResult) map[string]interface{} {
	var formatResult = map[string]interface{}{}

	formatResult["total_count"] = len(ar.FoundBadWords)

	var wordSlice []map[string]interface{}

	for key, value := range ar.FoundBadWords {
		if ar.DataType == "text" {
			wordSlice = append(wordSlice, map[string]interface{}{"word": value, "index": key})
		} else if ar.DataType == "array" {
			wordSlice = append(wordSlice, map[string]interface{}{"word": value, "expr_index": strings.Split(key, "/")[1], "index": strings.Split(key, "/")[0]})
		}
	}

	formatResult["words"] = wordSlice

	return formatResult
}

func RequestAndAnalize(wg *sync.WaitGroup, resultChan chan<- AnalizationResult, ed models.EndpointData, badWords []string) {
	var err error
	var response *http.Response

	var result = AnalizationResult{
		URL: ed.URL,
		Err: nil,
	}

	for i := 0; i < ed.MaxRetries; i++ {
		response, err = http.Get(ed.URL)
		if err == nil {
			break
		}
	}

	if err != nil {
		result.Err = err
		resultChan <- result
		wg.Done()
	}

	if ed.ReturnData == "text" {
		var responseContents string

		err = json.NewDecoder(response.Body).Decode(&responseContents)
		if err != nil {
			result.Err = err
			resultChan <- result
			wg.Done()
			return
		}

		result.FoundBadWords = detector.TextDetectBadWords(responseContents, badWords)
		result.DataType = "text"

		resultChan <- result
		wg.Done()
		return
	} else if ed.ReturnData == "array" {
		var responseContents []string

		err = json.NewDecoder(response.Body).Decode(&responseContents)
		if err != nil {
			result.Err = err
			resultChan <- result
			wg.Done()
			return
		}

		result.FoundBadWords = detector.ArrayDetectBadWords(responseContents, badWords)
		result.DataType = "array"

		resultChan <- result
		wg.Done()
		return
	}
	err = errors.New("Error occured while analizing data")
	result.Err = err

	resultChan <- result
	wg.Done()
}

func SendRequests(c *config.Config, l *zap.SugaredLogger) {
	var wg sync.WaitGroup
	var outputPath string = fmt.Sprintf("../output/%s.json", strconv.FormatInt(time.Now().Unix(), 10))

	recieverChan := make(chan AnalizationResult)

	for _, endpointData := range c.ListOfEndpoints {
		wg.Add(1)

		go func(wg *sync.WaitGroup, recieverChan chan AnalizationResult, data models.EndpointData, c *config.Config) {
			RequestAndAnalize(wg, recieverChan, data, c.BadWords)
		}(&wg, recieverChan, endpointData, c)

	}

	go func() {
		wg.Wait()
		close(recieverChan)
	}()

	outputResultChan := make(chan map[string]interface{}, 1)

	outputResult := map[string]interface{}{}

	for recievedData := range recieverChan {

		if recievedData.Err != nil {
			l.Fatalf("Error occured while requesting and analizing the data: %s", recievedData.Err)
		}

		outputResult[recievedData.URL] = FormatAnalizationResult(recievedData)

	}

	outputResultChan <- outputResult

	formattedJSON, err := json.MarshalIndent(<-outputResultChan, "", "    ")

	close(outputResultChan)

	file, err := os.Create(outputPath)
	if err != nil {
		l.Fatalf("Error occurred while trying to create output file: %s", err)
	}

	_, err = file.Write(formattedJSON)
	if err != nil {
		l.Error("Failed to write data into file")
	}
}

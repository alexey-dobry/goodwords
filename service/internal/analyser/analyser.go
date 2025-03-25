package analyser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/alexey-dobry/goodwords/internal/config"
	"github.com/alexey-dobry/goodwords/internal/models"
	"go.uber.org/zap"
)

type analizationResult struct {
	FoundBadWords [][2]string
	URL           string
	DataType      string
	Err           error
}

func requestAndAnalize(wg *sync.WaitGroup, resultChan chan<- analizationResult, ed models.EndpointData, badWords []string) {
	var err error
	var response *http.Response

	var result = analizationResult{
		URL: ed.URL,
		Err: nil,
	}

	client := http.Client{
		Timeout: time.Duration(ed.MaxTime) * time.Second,
	}

	for i := 0; i < ed.MaxRetries; i++ {
		response, err = client.Get(ed.URL)
		if err == nil {
			break
		} else if i == ed.MaxRetries-1 {
			result.Err = fmt.Errorf("too many retries")
			resultChan <- result
			wg.Done()
			return
		}
	}

	if err != nil {
		result.Err = err
		resultChan <- result
		wg.Done()
		return
	}

	if response.Status != "200 OK" {
		result.Err = fmt.Errorf("Response staus: ", response.Status)
		resultChan <- result
		wg.Done()
		return
	}

	if ed.ReturnData == "text" {
		var responseContents string

		err = json.NewDecoder(response.Body).Decode(&responseContents)
		if err != nil {
			result.Err = fmt.Errorf("Error parsing data from endpoint")
			resultChan <- result
			wg.Done()
			return
		}

		result.FoundBadWords = textDetectBadWords(responseContents, badWords)
		result.DataType = "text"

		resultChan <- result
		wg.Done()
		return

	} else if ed.ReturnData == "array" {
		var responseContents []string

		err = json.NewDecoder(response.Body).Decode(&responseContents)
		if err != nil {
			result.Err = fmt.Errorf("Error parsing data from endpoint")
			resultChan <- result
			wg.Done()
			return
		}

		result.FoundBadWords = arrayDetectBadWords(responseContents, badWords)
		result.DataType = "array"

		resultChan <- result
		wg.Done()
		return
	}

	err = fmt.Errorf("Wrong return datatype on endpoint: %s", ed.URL)
	result.Err = err

	resultChan <- result
	wg.Done()
}

func SendRequests(c *config.Config, l *zap.SugaredLogger) []byte {
	var wg sync.WaitGroup

	recieverChan := make(chan analizationResult)

	for _, endpointData := range c.ListOfEndpoints {
		wg.Add(1)

		go func(wg *sync.WaitGroup, recieverChan chan analizationResult, data models.EndpointData, c *config.Config) {
			requestAndAnalize(wg, recieverChan, data, c.BadWords)
		}(&wg, recieverChan, endpointData, c)

	}

	go func() {
		wg.Wait()
		close(recieverChan)
	}()

	outputResultChan := make(chan map[string]interface{}, 1)

	go func(ch chan<- map[string]interface{}) {
		outputResult := map[string]interface{}{}

		for recievedData := range recieverChan {
			if recievedData.Err != nil && recievedData.Err.Error() == "too many retries" {
				outputResult[recievedData.URL] = "too many retries"

				l.Error("Error occured while requesting and analyzing the data from ", recievedData.URL, " : ", recievedData.Err)

			} else if recievedData.Err != nil {
				l.Error("Error occured while requesting and analyzing the data from ", recievedData.URL, " : ", recievedData.Err)

			} else {
				outputResult[recievedData.URL] = formatAnalizationResult(recievedData)

				l.Info("Successfully analyzed data from ", recievedData.URL)

			}
		}

		outputResultChan <- outputResult

	}(outputResultChan)

	formatedJSON, err := json.MarshalIndent(<-outputResultChan, "", "    ")
	if err != nil {
		l.Error("Error marshalling data into json")
	}

	close(outputResultChan)

	return formatedJSON
}

func RunAnalyser(c *config.Config, l *zap.SugaredLogger) {
	var outputFile string = fmt.Sprintf("../output/%s.json", strconv.FormatInt(time.Now().Unix(), 10))

	formatedJSON := SendRequests(c, l)

	file, err := os.Create(outputFile)
	if err != nil {
		l.Fatalf("Error occurred while trying to create output file: %s", err)
	}

	_, err = file.Write(formatedJSON)
	if err != nil {
		l.Error("Failed to write data into file")
	}
}

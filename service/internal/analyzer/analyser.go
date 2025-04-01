package analyzer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alexey-dobry/goodwords/internal/config"
	"go.uber.org/zap"
)

type ReturnDatatype string

const (
	ReturnText  = "text"
	ReturnArray = "array"
)

type analizationResult struct {
	FoundBadWords []badWord
	URL           string
	DataType      ReturnDatatype
	Err           error
}

type badWord struct {
	Word      string
	ExprIndex int
	Index     int
}

func requestAndAnalyze(wg *sync.WaitGroup, resultChan chan<- analizationResult, ed config.ConfigEndpointData, badWords []string) {
	var err error
	var response *http.Response
	defer wg.Done()

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
			return
		}
	}

	if err != nil {
		result.Err = err
		resultChan <- result
		return
	}

	if response.StatusCode != 200 {
		result.Err = fmt.Errorf("response status code : %d", response.StatusCode)
		resultChan <- result
		return
	}

	if ed.ReturnData == ReturnText {
		var responseContents string

		err = json.NewDecoder(response.Body).Decode(&responseContents)
		if err != nil {
			result.Err = fmt.Errorf("error parsing data from endpoint")
			resultChan <- result
			return
		}

		result.FoundBadWords = textDetectBadWords(responseContents, badWords)
		result.DataType = ReturnText

		resultChan <- result
		return

	} else if ed.ReturnData == ReturnArray {
		var responseContents []string

		err = json.NewDecoder(response.Body).Decode(&responseContents)
		if err != nil {
			result.Err = fmt.Errorf("error parsing data from endpoint")
			resultChan <- result
			return
		}

		result.FoundBadWords = arrayDetectBadWords(responseContents, badWords)
		result.DataType = ReturnArray

		resultChan <- result
		return
	}

	err = fmt.Errorf("wrong return datatype on endpoint")
	result.Err = err

	resultChan <- result
	return
}

func SendRequests(c *config.Config, l *zap.SugaredLogger) []byte {
	var wg sync.WaitGroup

	receiverChan := make(chan analizationResult)

	for _, endpointData := range c.ListOfEndpoints {
		wg.Add(1)

		go func() {
			requestAndAnalyze(&wg, receiverChan, endpointData, c.BadWords)
		}()

	}

	go func() {
		wg.Wait()
		close(receiverChan)
	}()

	outputResult := map[string]interface{}{}

	for receivedData := range receiverChan {
		if receivedData.Err != nil {
			outputResult[receivedData.URL] = receivedData.Err.Error()

			l.Error("Error occurred while requesting and analyzing the data from ", receivedData.URL, " : ", receivedData.Err)
		} else {
			outputResult[receivedData.URL] = formatAnalizationResult(receivedData)

			l.Info("Successfully analyzed data from ", receivedData.URL)
		}
	}

	formattedJSON, err := json.MarshalIndent(outputResult, "", "    ")
	if err != nil {
		l.Error("Error marshalling data into json")
	}

	return formattedJSON
}

func RunAnalyzer(c *config.Config, l *zap.SugaredLogger) {
	var outputFile string = fmt.Sprintf("../output/%d.json", time.Now().Unix())

	formattedJSON := SendRequests(c, l)

	file, err := os.Create(outputFile)
	if err != nil {
		l.Fatalf("Error occurred while trying to create output file: %s", err)
	}

	_, err = file.Write(formattedJSON)
	if err != nil {
		l.Error("Failed to write data into file")
	}
}

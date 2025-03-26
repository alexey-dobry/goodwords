package analyser

import (
	"encoding/json"
	"errors"
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

func requestAndAnalize(wg *sync.WaitGroup, resultChan chan<- analizationResult, ed config.ConfigEndpointData, badWords []string) {
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
			result.Err = errors.New("too many retries")
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
		result.Err = errors.New(fmt.Sprintf("Response staus: %s", response.StatusCode))
		resultChan <- result
		return
	}

	if ed.ReturnData == ReturnText {
		var responseContents string

		err = json.NewDecoder(response.Body).Decode(&responseContents)
		if err != nil {
			result.Err = errors.New("Error parsing data from endpoint")
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
			result.Err = errors.New("Error parsing data from endpoint")
			resultChan <- result
			return
		}

		result.FoundBadWords = arrayDetectBadWords(responseContents, badWords)
		result.DataType = ReturnArray

		resultChan <- result
		return
	}

	err = errors.New(fmt.Sprintf("Wrong return datatype on endpoint: %s", ed.URL))
	result.Err = err

	resultChan <- result
	return
}

func SendRequests(c *config.Config, l *zap.SugaredLogger) []byte {
	var wg sync.WaitGroup

	recieverChan := make(chan analizationResult)

	for _, endpointData := range c.ListOfEndpoints {
		wg.Add(1)

		go func() {
			requestAndAnalize(&wg, recieverChan, endpointData, c.BadWords)
		}()

	}

	go func() {
		wg.Wait()
		close(recieverChan)
	}()

	outputResult := map[string]interface{}{}

	for recievedData := range recieverChan {
		if recievedData.Err != nil && recievedData.Err == fmt.Errorf("too many retries") {
			outputResult[recievedData.URL] = "too many retries"

			l.Error("Error occured while requesting and analysing the data from ", recievedData.URL, " : ", recievedData.Err)

		} else if recievedData.Err != nil {
			l.Error("Error occured while requesting and analysing the data from ", recievedData.URL, " : ", recievedData.Err)

		} else {
			outputResult[recievedData.URL] = formatAnalisationResult(recievedData)

			l.Info("Successfully analysed data from ", recievedData.URL)
		}
	}

	formatedJSON, err := json.MarshalIndent(outputResult, "", "    ")
	if err != nil {
		l.Error("Error marshalling data into json")
	}

	return formatedJSON
}

func RunAnalyser(c *config.Config, l *zap.SugaredLogger) {
	var outputFile string = fmt.Sprintf("../output/%d.json", time.Now().Unix())

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

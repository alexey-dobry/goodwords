package service_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexey-dobry/goodwords/internal/analyser"
	"github.com/alexey-dobry/goodwords/internal/config"
	"github.com/alexey-dobry/goodwords/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestTextResponseWithBadWord(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("hello what bad GOpher good Python")
	}))
	defer server.Close()

	mockConfig := config.Config{
		BadWords: []string{"bad gopher", "good python"},
		ListOfEndpoints: []config.ConfigEndpointData{
			config.ConfigEndpointData{
				URL:        server.URL,
				MaxTime:    5,
				MaxRetries: 5,
				ReturnData: "text",
			},
		},
	}

	l, err := logger.InitLogger()

	assert.Equal(t, err, nil)

	actualResult := analyser.SendRequests(&mockConfig, l)

	expectedResult, _ := json.MarshalIndent(
		map[string]interface{}{
			server.URL: map[string]interface{}{
				"total_count": 2,
				"words": []map[string]interface{}{
					map[string]interface{}{
						"index": "11",
						"word":  "bad gopher",
					},
					map[string]interface{}{
						"index": "22",
						"word":  "good python",
					},
				},
			},
		}, "", "    ")

	assert.Equal(t, expectedResult, actualResult)
}

func TestArrayResponseWithBadWord(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]string{"hello what bad GOpher", "hello good Python bad GOpher"})
	}))
	defer server.Close()

	mockConfig := config.Config{
		BadWords: []string{"bad gopher", "good python"},
		ListOfEndpoints: []config.ConfigEndpointData{
			config.ConfigEndpointData{
				URL:        server.URL,
				MaxTime:    5,
				MaxRetries: 5,
				ReturnData: "array",
			},
		},
	}

	l, err := logger.InitLogger()

	assert.Equal(t, err, nil)

	actualResult := analyser.SendRequests(&mockConfig, l)

	expectedResult, _ := json.MarshalIndent(
		map[string]interface{}{
			server.URL: map[string]interface{}{
				"total_count": 3,
				"words": []map[string]interface{}{
					map[string]interface{}{
						"expr_index": "0",
						"index":      "11",
						"word":       "bad gopher",
					},
					map[string]interface{}{
						"expr_index": "1",
						"index":      "18",
						"word":       "bad gopher",
					},
					map[string]interface{}{
						"expr_index": "1",
						"index":      "6",
						"word":       "good python",
					},
				},
			},
		}, "", "    ")

	assert.Equal(t, expectedResult, actualResult)
}

func TestNoResponseWithBadWord(t *testing.T) {
	mockConfig := config.Config{
		BadWords: []string{"bad gopher"},
		ListOfEndpoints: []config.ConfigEndpointData{
			config.ConfigEndpointData{
				URL:        "http://localhost:8001",
				MaxTime:    5,
				MaxRetries: 5,
				ReturnData: "text",
			},
		},
	}

	l, err := logger.InitLogger()

	assert.Equal(t, err, nil)

	actualResult := analyser.SendRequests(&mockConfig, l)

	expectedResult, _ := json.MarshalIndent(
		map[string]interface{}{
			"http://localhost:8001": "too many retries",
		}, "", "    ")

	assert.Equal(t, expectedResult, actualResult)
}

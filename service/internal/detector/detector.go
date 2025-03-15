package detector

import (
	"strings"
)

type AnalizationResult struct {
	FoundBadWords map[int]string
}

func StringDetectBadWords(data string, badwords []string) *AnalizationResult {
	var result AnalizationResult

	data = strings.ToLower(data)

	for _, badValue := range badwords {
		if index := strings.Index(data, badValue); index != -1 {
			result.FoundBadWords[index] = badValue
		}
	}

	return &result
}

func SliceDetectBadWords(data []string, badwords []string) *AnalizationResult {
	var result AnalizationResult

	for _, stringValue := range data {
		for _, badValue := range badwords {

			stringValue = strings.ToLower(stringValue)

			if index := strings.Index(stringValue, badValue); index != -1 {
				result.FoundBadWords[index] = badValue
			}
		}
	}

	return &result
}

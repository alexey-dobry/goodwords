package detector

import (
	"strconv"
	"strings"
)

func TextDetectBadWords(data string, badwords []string) map[string]string {
	var result = map[string]string{}

	data = strings.ToLower(data)

	for _, badValue := range badwords {
		if index := strings.Index(data, badValue); index != -1 {
			result[strconv.Itoa(index)] = badValue
		}
	}

	return result
}

func ArrayDetectBadWords(data []string, badwords []string) map[string]string {
	var result = map[string]string{}

	for stringIndex, stringValue := range data {

		for _, badValue := range badwords {

			stringValue = strings.ToLower(stringValue)

			if index := strings.Index(stringValue, badValue); index != -1 {

				index_of_val := strconv.Itoa(index) + "/" + strconv.Itoa(stringIndex)

				result[index_of_val] = badValue
			}
		}
	}

	return result
}

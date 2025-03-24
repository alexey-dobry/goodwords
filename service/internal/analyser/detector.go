package analyser

import (
	"strconv"
	"strings"
)

func textDetectBadWords(data string, badwords []string) [][2]string {
	result := make([][2]string, 0, 5)

	data = strings.ToLower(data)

	for _, badValue := range badwords {
		if index := strings.Index(data, strings.ToLower(badValue)); index != -1 {
			result = append(result, [2]string{strconv.Itoa(index), strings.ToLower(badValue)})
		}
	}

	return result
}

func arrayDetectBadWords(data []string, badwords []string) [][2]string {
	result := make([][2]string, 0, 5)

	for stringIndex, stringValue := range data {

		for _, badValue := range badwords {

			stringValue = strings.ToLower(stringValue)

			if index := strings.Index(stringValue, strings.ToLower(badValue)); index != -1 {

				index_of_val := strconv.Itoa(index) + "/" + strconv.Itoa(stringIndex)

				result = append(result, [2]string{index_of_val, strings.ToLower(badValue)})
			}
		}
	}

	return result
}

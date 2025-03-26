package analyser

import "strings"

func textDetectBadWords(data string, badwords []string) []badWord {
	result := make([]badWord, 0)

	data = strings.ToLower(data)

	for _, badValue := range badwords {
		if index := strings.Index(data, strings.ToLower(badValue)); index != -1 {
			result = append(result, badWord{badValue, 0, index})
		}
	}

	return result
}

func arrayDetectBadWords(data []string, badwords []string) []badWord {
	result := make([]badWord, 0)

	for stringIndex, stringValue := range data {

		for _, badValue := range badwords {

			stringValue = strings.ToLower(stringValue)

			if index := strings.Index(stringValue, strings.ToLower(badValue)); index != -1 {

				result = append(result, badWord{badValue, stringIndex, index})
			}
		}
	}

	return result
}

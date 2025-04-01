package analyzer

import "strings"

func FindAllIndexes(str, substr string) []int {
	indexes := make([]int, 0)
	start := 0
	subLen := len(substr) - 1

	if subLen == 0 {
		return indexes
	}

	for start < len(str)-subLen {
		idx := strings.Index(str[start:], substr)
		if idx == -1 {
			break
		}

		realIdx := start + idx
		indexes = append(indexes, realIdx)

		start = realIdx + subLen
	}

	return indexes
}

func textDetectBadWords(data string, badWords []string) []badWord {
	result := make([]badWord, 0)

	data = strings.ToLower(data)

	for _, badValue := range badWords {

		indexes := FindAllIndexes(data, badValue)

		if len(indexes) != 0 {

			for _, index := range indexes {
				result = append(result, badWord{badValue, 0, index})
			}

		}
	}

	return result
}

func arrayDetectBadWords(data []string, badWords []string) []badWord {
	result := make([]badWord, 0)

	for stringIndex, stringValue := range data {

		stringValue = strings.ToLower(stringValue)

		for _, badValue := range badWords {

			indexes := FindAllIndexes(stringValue, badValue)

			if len(indexes) != 0 {

				for _, index := range indexes {
					result = append(result, badWord{badValue, stringIndex, index})
				}

			}
		}
	}

	return result
}

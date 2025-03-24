package analyser

import "strings"

func formatAnalizationResult(ar analizationResult) map[string]interface{} {
	var formatResult = map[string]interface{}{}

	formatResult["total_count"] = len(ar.FoundBadWords)

	var wordSlice []map[string]interface{}

	for _, value := range ar.FoundBadWords {

		key := value[0]

		if ar.DataType == "text" {
			wordSlice = append(wordSlice, map[string]interface{}{"word": value[1], "index": key})
		} else if ar.DataType == "array" {
			wordSlice = append(wordSlice, map[string]interface{}{"word": value[1], "expr_index": strings.Split(key, "/")[1], "index": strings.Split(key, "/")[0]})
		}
	}

	formatResult["words"] = wordSlice

	return formatResult
}

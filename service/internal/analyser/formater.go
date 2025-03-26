package analyser

func formatAnalisationResult(ar analizationResult) map[string]interface{} {
	var formatResult = map[string]interface{}{}

	formatResult["total_count"] = len(ar.FoundBadWords)

	var wordSlice []map[string]interface{}

	for _, badWordData := range ar.FoundBadWords {

		if ar.DataType == ReturnText {
			wordSlice = append(wordSlice, map[string]interface{}{"word": badWordData.Word, "index": badWordData.Index})
		} else if ar.DataType == ReturnArray {
			wordSlice = append(wordSlice, map[string]interface{}{"word": badWordData.Word, "expr_index": badWordData.ExprIndex, "index": badWordData.Index})
		}
	}

	formatResult["words"] = wordSlice

	return formatResult
}

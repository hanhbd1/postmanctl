package util

func ReformatMap(nested map[string]interface{}, removeNilValue bool, keymap map[string]int) map[string]interface{} {
	newMap := make(map[string]interface{})
	for k, v := range nested {
		if _, ok := keymap[k]; ok {
			continue
		}
		if v == nil && removeNilValue {
			continue
		}
		newMap[k] = ReformatInterface(v, removeNilValue, keymap)
	}
	return newMap
}

func ReformatArrays(nested []interface{}, removeNilValue bool, keymap map[string]int) []interface{} {
	newArrays := make([]interface{}, 0)
	for _, v := range nested {
		if v == nil {
			continue
		}
		newArrays = append(newArrays, ReformatInterface(v, removeNilValue, keymap))
	}
	return newArrays
}

func ReformatInterface(nested interface{}, removeNilValue bool, keymap map[string]int) interface{} {
	var tmp interface{}
	if nested == nil {
		return nil
	}
	switch s := nested.(type) {
	case map[string]interface{}:
		tmpValue := ReformatMap(s, removeNilValue, keymap)
		tmp = tmpValue
	case []interface{}:
		tmpValue := ReformatArrays(s, removeNilValue, keymap)
		tmp = tmpValue
	default:
		tmp = nested
	}
	return tmp
}

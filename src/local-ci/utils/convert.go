package utils

func InterfaceArrayToStringArray(interfaces []interface{}) ([]string, bool) {
	strings := make([]string, len(interfaces))

	for index, value := range interfaces {
		if str, ok := value.(string); ok {
			strings[index] = str
		} else {
			return nil, false
		}
	}

	return strings, true
}

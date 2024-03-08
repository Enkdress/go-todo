package utils

func CreateReturnObject[T comparable](data []T) map[string][]T {
	returnObj := make(map[string][]T)
	if len(data) == 0 {
		returnObj["data"] = make([]T, 0, 0)
	} else {
		returnObj["data"] = data
	}

	return returnObj
}
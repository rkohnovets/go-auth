package utils

import "encoding/json"

func GetObjectJsonString(obj interface{}) (string, error) {
	result, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return "", err
	}
	return string(result), nil
}

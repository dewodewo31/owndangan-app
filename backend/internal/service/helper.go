package service

import (
	"encoding/json"

	"gorm.io/datatypes"
)

func datatypesJSON(data map[string]interface{}) datatypes.JSON {
	if data == nil {
		return datatypes.JSON(nil)
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return datatypes.JSON(nil)
	}
	return bytes
}

package files

import (
	"encoding/json"
)

type JSONFormat struct {
	CapableFormat
}

func (js *JSONFormat) Unmarshal(source []byte, configuration interface{}) error {
	return json.Unmarshal(source, configuration)
}

func NewJSONFormatUnmarshaler() *JSONFormat {
	return &JSONFormat{CapableFormat{formats: []string{".json"}}}
}

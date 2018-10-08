package files

import (
	"errors"

	"github.com/ghodss/yaml"
)

type YAMLFormat struct {
	CapableFormat
}

func (yf *YAMLFormat) Unmarshal(source []byte, configuration interface{}) error {
	if len(source) == 0 {
		return errors.New("Empty input")
	}
	return yaml.Unmarshal(source, configuration)
}

func NewYAMLFormatUnmarshaler() *YAMLFormat {
	return &YAMLFormat{CapableFormat{formats: []string{".yaml", ".yml"}}}
}

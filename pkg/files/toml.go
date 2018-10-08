package files

import (
	"errors"

	"github.com/pelletier/go-toml"
)

type TOMLFormat struct {
	CapableFormat
}

func (tf *TOMLFormat) Unmarshal(source []byte, configuration interface{}) error {
	if len(source) == 0 {
		return errors.New("Empty input")
	}
	return toml.Unmarshal(source, configuration)
}

func NewTOMLFormatUnmarshaler() *TOMLFormat {
	return &TOMLFormat{CapableFormat{formats: []string{".tom", ".toml"}}}
}

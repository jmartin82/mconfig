package files

import (
	"path/filepath"
)

type Unmarshaler interface {
	IsCapable(filename string) bool
	Unmarshal(source []byte, configuration interface{}) error
}

type CapableFormat struct {
	formats []string
}

func (cf CapableFormat) IsCapable(filename string) bool {
	ext := filepath.Ext(filename)
	for _, form := range cf.formats {
		if form == ext {
			return true
		}
	}
	return false
}

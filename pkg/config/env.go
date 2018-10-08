package config

import (
	"os"

	"github.com/subosito/gotenv"
)

type EnvFileReader interface {
	Load(filenames ...string) error
}

func NewGotEnvAdapter() *GotEnvAdapter {
	return &GotEnvAdapter{}
}

type GotEnvAdapter struct {
}

func (gea *GotEnvAdapter) Load(filenames ...string) error {
	for _, file := range filenames {
		if _, err := os.Stat(file); err == nil {
			if err := gotenv.Load(file); err != nil {
				return err
			}
		}
	}
	return nil
}

package config

import "github.com/subosito/gotenv"

type EnvFileReader interface {
	Load(filenames ...string) error
}

func NewGotEnvAdapter() *GotEnvAdapter {
	return &GotEnvAdapter{}
}

type GotEnvAdapter struct {
}

func (gea *GotEnvAdapter) Load(filenames ...string) error {
	return gotenv.Load(filenames...)
}

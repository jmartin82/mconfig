package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGotEnvAdapter_Load(t *testing.T) {
	type args struct {
		filenames []string
	}
	tests := []struct {
		name     string
		gea      *GotEnvAdapter
		args     args
		env      string
		content  string
		expected string
		wantErr  bool
	}{
		{
			"Test valid format",
			NewGotEnvAdapter(),
			args{
				[]string{".env"},
			},
			"ENV_TEST",
			"ENV_TEST=OK",
			"OK",
			false,
		},
		{
			"Test invalid format",
			NewGotEnvAdapter(),
			args{
				[]string{".env"},
			},
			"",
			"sdasdasd",
			"",
			true,
		},
		{
			"Invalid file",
			NewGotEnvAdapter(),
			args{
				[]string{".xxx"},
			},
			"",
			"",
			"",
			true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			cleanTestEnv([]string{tt.env})
			writeEnvFile(tt.args.filenames[0], tt.content)
			gea := &GotEnvAdapter{}
			if err := gea.Load(tt.args.filenames...); (err != nil) != tt.wantErr {
				t.Errorf("GotEnvAdapter.Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			if os.Getenv(tt.env) != tt.expected {
				t.Errorf("GotEnvAdapter.Load() got = %v, expected %v", os.Getenv(tt.env), tt.expected)
			}
		})
	}
}

func writeEnvFile(filename, content string) {
	ioutil.WriteFile(".env", []byte(content), 0644)
}

func cleanTestEnv(envs []string) {
	for _, env := range envs {
		os.Unsetenv(env)
	}
}

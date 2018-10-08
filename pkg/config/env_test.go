package config

import (
	"bufio"
	"fmt"
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
			false,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			cleanTestEnv([]string{tt.env})
			writeEnvFile(tt.args.filenames[0], []string{tt.content})
			defer removeTestFiles([]string{tt.args.filenames[0]})

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

func setTestEnv(envMap map[string]string) {
	for k, val := range envMap {
		os.Setenv(k, val)
	}

}

func writeEnvFile(filename string, lines []string) {
	file, _ := os.Create(filename)
	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	w.Flush()
	file.Close()
}

func cleanTestEnv(envs []string) {
	for _, env := range envs {
		os.Unsetenv(env)
	}
}

func removeTestFiles(files []string) {
	for _, file := range files {
		os.Remove(file)
	}
}

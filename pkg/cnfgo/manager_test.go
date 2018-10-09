package mconfig

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/go-test/deep"
)

type MysqlConfiguration struct {
	Host     string `env:"TEST_MYSQL_HOST"`
	Username string `env:"TEST_MYSQL_USERNAME"`
	Password string `env:"TEST_MYSQL_PASSWORD"`
	Database string `env:"TEST_MYSQL_DATABASE"`
	Port     int    `env:"TEST_MYSQL_PORT"`
}

type RedisConfiguration struct {
	Host string `env:"TEST_REDIS_HOST"`
	Port int    `env:"TEST_REDIS_PORT"`
}

type Configuration struct {
	Port  int `env:"TEST_APP_PORT"`
	Mysql MysqlConfiguration
	Redis RedisConfiguration
}

func TestReadFromFile(t *testing.T) {

	defer removeTestFiles([]string{"config.json"})
	ioutil.WriteFile("config.json", []byte("{\"Port\":3001,\"Mysql\":{\"Host\":\"192.168.0.1\",\"Username\":\"root\",\"Password\":\"test\",\"Database\":\"mconfig\",\"Port\":3306},\"Redis\":{\"Host\":\"localhost\",\"Port\":6379}}"), 0644)
	configuration := &Configuration{}
	m := getDefaultManager()
	m.ReadFromFile("config.json", configuration)

	expected := &Configuration{
		Port: 3001,
		Mysql: MysqlConfiguration{
			Host:     "192.168.0.1",
			Username: "root",
			Password: "test",
			Database: "mconfig",
			Port:     3306,
		},
		Redis: RedisConfiguration{
			Host: "localhost",
			Port: 6379,
		},
	}
	if diff := deep.Equal(configuration, expected); diff != nil {
		t.Error(diff)
	}
}

func TestReadFromEnvironment(t *testing.T) {

	cleanTestEnv([]string{
		"TEST_MYSQL_HOST",
		"TEST_MYSQL_USERNAME",
		"TEST_MYSQL_PASSWORD",
		"TEST_MYSQL_DATABASE",
		"TEST_MYSQL_PORT",
		"TEST_REDIS_HOST",
		"TEST_REDIS_PORT",
		"TEST_APP_PORT",
		"TEST_ENV_ANOTHER",
	})

	writeEnvFile(".env", []string{
		"TEST_MYSQL_HOST=192.168.0.1",
		"TEST_MYSQL_USERNAME=root",
		"TEST_MYSQL_PASSWORD=test",
		"TEST_MYSQL_DATABASE=mconfig",
		"TEST_MYSQL_PORT=3306",
		"TEST_REDIS_HOST=localhost",
		"TEST_REDIS_PORT=6379",
		"TEST_APP_PORT=3001",
	})

	defer removeTestFiles([]string{".env"})

	configuration := &Configuration{}
	m := getDefaultManager()
	m.ReadFromEnvironment(configuration)

	expected := &Configuration{
		Port: 3001,
		Mysql: MysqlConfiguration{
			Host:     "192.168.0.1",
			Username: "root",
			Password: "test",
			Database: "mconfig",
			Port:     3306,
		},
		Redis: RedisConfiguration{
			Host: "localhost",
			Port: 6379,
		},
	}

	if diff := deep.Equal(configuration, expected); diff != nil {
		t.Error(diff)
	}

}
func TestReadFromNotRegistreProvider(t *testing.T) {
	defer removeTestFiles([]string{"config.xml"})
	ioutil.WriteFile("config.xml", []byte(""), 0644)
	configuration := &Configuration{}
	m := getDefaultManager()
	if err := m.ReadFromFile("config.xml", configuration); err == nil {
		t.Errorf("Expected error reading file without provider")
	}
}

func TestInvalidConfigInput(t *testing.T) {
	configuration := Configuration{}
	m := getDefaultManager()
	if err := m.ReadFromEnvironment(configuration); err == nil {
		t.Errorf("Expected error when input is not pointer to struct")
	}

	defer removeTestFiles([]string{"config.json"})
	ioutil.WriteFile("config.json", []byte("{\"Port\":3001,\"Mysql\":{\"Host\":\"192.168.0.1\",\"Username\":\"root\",\"Password\":\"test\",\"Database\":\"mconfig\",\"Port\":3306},\"Redis\":{\"Host\":\"localhost\",\"Port\":6379}}"), 0644)
	m.ReadFromFile("config.json", configuration)
}

func TestInvalidFileName(t *testing.T) {
	configuration := &Configuration{}
	m := getDefaultManager()
	if err := m.ReadFromFile("", configuration); err == nil {
		t.Errorf("File should fail if is empty")
	}
	if err := m.ReadFromFile("dssddsf.json", configuration); err == nil {
		t.Errorf("File should fail if doesn't exist")
	}
}

func TestOKIfNotFoundEnvFilename(t *testing.T) {
	configuration := &Configuration{}
	m := getDefaultManager()
	if err := m.ReadFromEnvironment(configuration); err != nil {
		t.Errorf("File .env can not exists, Error: %s", err)
	}
}

func TestInvalidEnvFileFormat(t *testing.T) {
	configuration := &Configuration{}
	writeEnvFile(".env", []string{
		"asdasd",
	})

	defer removeTestFiles([]string{"config.json", ".env"})
	ioutil.WriteFile("config.json", []byte("{\"Port\":3001,\"Mysql\":{\"Host\":\"192.168.0.1\",\"Username\":\"root\",\"Password\":\"test\",\"Database\":\"mconfig\",\"Port\":3306},\"Redis\":{\"Host\":\"localhost\",\"Port\":6379}}"), 0644)
	if err := Parse("config.json", configuration); err == nil {
		t.Errorf("File should fail env file is not valid")
	}
}

func TestParseInvalidConfigFile(t *testing.T) {
	configuration := &Configuration{}

	if err := Parse("dssddsf.json", configuration); err == nil {
		t.Errorf("File should fail if doesn't exist")
	}
}

func TestParse(t *testing.T) {
	cleanTestEnv([]string{
		"TEST_MYSQL_PASSWORD",
	})

	writeEnvFile(".env", []string{
		"TEST_MYSQL_PASSWORD=super_secret",
	})

	defer removeTestFiles([]string{".env", "config.json"})
	ioutil.WriteFile("config.json", []byte("{\"Port\":3001,\"Mysql\":{\"Host\":\"192.168.0.1\",\"Username\":\"root\",\"Password\":\"test\",\"Database\":\"mconfig\",\"Port\":3306},\"Redis\":{\"Host\":\"localhost\",\"Port\":6379}}"), 0644)

	configuration := &Configuration{}
	Parse("config.json", configuration)

	expected := &Configuration{
		Port: 3001,
		Mysql: MysqlConfiguration{
			Host:     "192.168.0.1",
			Username: "root",
			Password: "super_secret",
			Database: "mconfig",
			Port:     3306,
		},
		Redis: RedisConfiguration{
			Host: "localhost",
			Port: 6379,
		},
	}

	if diff := deep.Equal(configuration, expected); diff != nil {
		t.Error(diff)
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

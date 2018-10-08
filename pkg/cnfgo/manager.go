package cnfgo

import (
	"errors"
	"io/ioutil"
	"reflect"

	"github.com/jmartin82/cnfgo/pkg/config"
	"github.com/jmartin82/cnfgo/pkg/files"
)

var (
	ErrNotAStructPtr = errors.New("Configuration should be a pointer to a struct type")
	ErrFileNotFound  = errors.New("Configuration file doesn't exist")
)

type ManagerFacade struct {
	unmarshalers  []files.Unmarshaler
	envFileReader config.EnvFileReader
	envTagParser  *config.EnvTagParser
	envFiles      []string
}

func (p *ManagerFacade) SetEnvFileReader(envFileReader config.EnvFileReader) {
	p.envFileReader = envFileReader
}

func (p *ManagerFacade) SetEnvFiles(files ...string) {
	p.envFiles = files
}

func (p *ManagerFacade) AddFileUnmashaler(unmarshaler files.Unmarshaler) {
	p.unmarshalers = append(p.unmarshalers, unmarshaler)
}
func (p *ManagerFacade) loadEnvFile() error {
	return p.envFileReader.Load(p.envFiles...)
}

func (p *ManagerFacade) getUnmarshaler(filename string) (files.Unmarshaler, error) {
	for _, unmarshaler := range p.unmarshalers {
		if unmarshaler.IsCapable(filename) {
			return unmarshaler, nil
		}
	}
	return nil, errors.New("Not valid config file reader")
}

func NewManager() *ManagerFacade {
	envTagParser := config.NewEnvTagParser("env")
	envFiles := []string{}
	unmarshalers := []files.Unmarshaler{}

	return &ManagerFacade{unmarshalers: unmarshalers, envFiles: envFiles, envTagParser: envTagParser}
}

func getDefaultManager() *ManagerFacade {
	managerFacade := NewManager()
	managerFacade.AddFileUnmashaler(files.NewJSONFormatUnmarshaler())
	managerFacade.AddFileUnmashaler(files.NewYAMLFormatUnmarshaler())
	managerFacade.AddFileUnmashaler(files.NewTOMLFormatUnmarshaler())
	managerFacade.SetEnvFileReader(config.NewGotEnvAdapter())
	managerFacade.SetEnvFiles(".env")
	return managerFacade
}

var Manager *ManagerFacade = getDefaultManager()

func isPtrToStruct(configuration interface{}) bool {
	refValue := reflect.ValueOf(configuration)
	typ := refValue.Type()
	return typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct
}
func ReadFromEnvironment(configuration interface{}) error {

	if !isPtrToStruct(configuration) {
		return ErrNotAStructPtr
	}

	if err := Manager.loadEnvFile(); err != nil {
		return err
	}

	Manager.envTagParser.Parse(configuration, 0)
	return nil
}

func ReadFromFile(filename string, configuration interface{}) (err error) {
	if len(filename) == 0 {
		return ErrFileNotFound
	}

	if !isPtrToStruct(configuration) {
		return ErrNotAStructPtr
	}

	unmarshaler, err := Manager.getUnmarshaler(filename)
	if err != nil {
		return err
	}

	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return unmarshaler.Unmarshal(source, configuration)
}

func Parse(filename string, configuration interface{}) (err error) {

	if err := ReadFromFile(filename, configuration); err != nil {
		return err
	}
	if err := ReadFromEnvironment(configuration); err != nil {
		return err
	}

	return nil
}

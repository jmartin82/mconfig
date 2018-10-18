package mconfig

import (
	"errors"
	"io/ioutil"
	"reflect"

	"github.com/jmartin82/mconfig/pkg/config"
	"github.com/jmartin82/mconfig/pkg/files"
)

var (
	ErrNotAStructPtr        = errors.New("configuration should be a pointer to a struct type")
	ErrFileNotFound         = errors.New("configuration file doesn't exist")
	ErrNotValidConfigReader = errors.New("not valid config file reader")
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

func (p *ManagerFacade) isPtrToStruct(configuration interface{}) bool {
	refValue := reflect.ValueOf(configuration)
	typ := refValue.Type()
	return typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct
}
func (p *ManagerFacade) ReadFromEnvironment(configuration interface{}) error {

	if !p.isPtrToStruct(configuration) {
		return ErrNotAStructPtr
	}

	if err := p.loadEnvFile(); err != nil {
		return err
	}

	p.envTagParser.Parse(configuration, 0)
	return nil
}

func (p *ManagerFacade) ReadFromFile(filename string, configuration interface{}) (err error) {
	if len(filename) == 0 {
		return ErrFileNotFound
	}

	if !p.isPtrToStruct(configuration) {
		return ErrNotAStructPtr
	}

	unmarshaler, err := p.getUnmarshaler(filename)
	if err != nil {
		return err
	}

	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return unmarshaler.Unmarshal(source, configuration)
}

func (p *ManagerFacade) getUnmarshaler(filename string) (files.Unmarshaler, error) {
	for _, unmarshaler := range p.unmarshalers {
		if unmarshaler.IsCapable(filename) {
			return unmarshaler, nil
		}
	}
	return nil, ErrNotValidConfigReader
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

var ConfigManager *ManagerFacade = getDefaultManager()

func Parse(filename string, configuration interface{}) (err error) {

	if err := ConfigManager.ReadFromFile(filename, configuration); err != nil {
		return err
	}
	if err := ConfigManager.ReadFromEnvironment(configuration); err != nil {
		return err
	}

	return nil
}

package config

import (
	"os"
	"reflect"
	"strconv"
)

type EnvTagParser struct {
	tagName string
}

func NewEnvTagParser(tagName string) *EnvTagParser {
	return &EnvTagParser{tagName: tagName}
}

func (etp *EnvTagParser) Parse(configuration interface{}, deep int) {
	value := reflect.ValueOf(configuration)

	if value.Kind() == reflect.Ptr && !value.IsNil() {
		value = value.Elem()
	}

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		kind := field.Kind()
		//recursive call to go deep in the struct
		if kind == reflect.Struct {
			etp.Parse(field.Addr().Interface(), deep+1)
			continue
		}
		if kind == reflect.Ptr && field.Type().Elem().Kind() == reflect.Struct {
			if field.IsNil() {
				field.Set(reflect.New(field.Type().Elem()))
			}
			etp.Parse(field.Interface(), deep+1)
			continue
		}

		typeField := value.Type().Field(i)
		envName := typeField.Tag.Get(etp.tagName)
		if len(envName) == 0 || !field.CanSet() {
			continue

		}

		envVar, present := os.LookupEnv(envName)
		if present {
			etp.set(kind, field, envVar, false)
		}

	}

}

func (etp *EnvTagParser) set(kind reflect.Kind, field reflect.Value, value string, ptr bool) {
	switch kind {
	case reflect.Ptr:
		etp.set(field.Type().Elem().Kind(), field, value, true)
	case reflect.Bool:
		etp.setEnvBool(field, value, ptr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		etp.setEnvInt(field, value, ptr)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		etp.setEnvUint(field, value, ptr)
	case reflect.Float32, reflect.Float64:
		etp.setEnvFloat(field, value, ptr)
	case reflect.String:
		etp.setEnvString(field, value, ptr)
	}
}

func (etp *EnvTagParser) setEnvBool(field reflect.Value, value string, ptr bool) {
	if i, err := strconv.ParseBool(value); err == nil {
		if ptr {
			field.Set(reflect.ValueOf(&i))
		} else {
			field.SetBool(i)
		}

	}
}
func (etp *EnvTagParser) setEnvInt(field reflect.Value, value string, ptr bool) {
	if i, err := strconv.ParseInt(value, 10, 64); err == nil {
		if ptr {
			switch field.Type().Elem().Kind() {
			case reflect.Int:
				v := int(i)
				field.Set(reflect.ValueOf(&v))
			case reflect.Int8:
				v := int8(i)
				field.Set(reflect.ValueOf(&v))
			case reflect.Int16:
				v := int16(i)
				field.Set(reflect.ValueOf(&v))
			case reflect.Int32:
				v := int32(i)
				field.Set(reflect.ValueOf(&v))
			default:
				field.Set(reflect.ValueOf(&i))
			}

		} else {
			field.SetInt(i)
		}

	}
}

func (etp *EnvTagParser) setEnvUint(field reflect.Value, value string, ptr bool) {
	if i, err := strconv.ParseUint(value, 10, 64); err == nil {
		if ptr {
			switch field.Type().Elem().Kind() {
			case reflect.Uint:
				v := uint(i)
				field.Set(reflect.ValueOf(&v))
			case reflect.Uint8:
				v := uint8(i)
				field.Set(reflect.ValueOf(&v))
			case reflect.Uint16:
				v := uint16(i)
				field.Set(reflect.ValueOf(&v))
			case reflect.Uint32:
				v := uint32(i)
				field.Set(reflect.ValueOf(&v))
			default:
				field.Set(reflect.ValueOf(&i))
			}
		} else {
			field.SetUint(i)
		}
	}
}
func (etp *EnvTagParser) setEnvFloat(field reflect.Value, value string, ptr bool) {
	if i, err := strconv.ParseFloat(value, 64); err == nil {
		if ptr {
			if field.Type().Elem().Kind() == reflect.Float32 {
				v := float32(i)
				field.Set(reflect.ValueOf(&v))
			} else {
				field.Set(reflect.ValueOf(&i))
			}

		} else {
			field.SetFloat(i)
		}
	}
}

func (etp *EnvTagParser) setEnvString(field reflect.Value, value string, ptr bool) {
	if ptr {
		field.Set(reflect.ValueOf(&value))
	} else {
		field.SetString(value)
	}
}

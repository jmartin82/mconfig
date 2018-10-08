package config

import (
	"testing"

	"github.com/go-test/deep"
)

type TestTypes struct {
	IsInt        int      `env:"ENV_INT"`
	IsUint       uint     `env:"ENV_UINT"`
	IsString     string   `env:"ENV_STRING"`
	IsFloat      float64  `env:"ENV_FLOAT"`
	IsBool       bool     `env:"ENV_BOOL"`
	PtrIsInt     *int     `env:"ENV_INT"`
	PtrIsInt8    *int8    `env:"ENV_INT"`
	PtrIsInt16   *int16   `env:"ENV_INT"`
	PtrIsInt32   *int32   `env:"ENV_INT"`
	PtrIsInt64   *int64   `env:"ENV_INT"`
	PtrIsUint    *uint    `env:"ENV_UINT"`
	PtrIsUint8   *uint8   `env:"ENV_UINT"`
	PtrIsUint16  *uint16  `env:"ENV_UINT"`
	PtrIsUint32  *uint32  `env:"ENV_UINT"`
	PtrIsUint64  *uint64  `env:"ENV_UINT"`
	PtrIsString  *string  `env:"ENV_STRING"`
	PtrIsFloat32 *float32 `env:"ENV_FLOAT"`
	PtrIsFloat64 *float64 `env:"ENV_FLOAT"`
	PtrIsBool    *bool    `env:"ENV_BOOL"`
	inmutable    bool     `env:"ENV_BOOL"`
}

var isInt = -1
var isUInt uint = 1
var isString = "string"
var isFloat = 1.1
var isBoolean = true

var ExpectedTestTypes = TestTypes{
	IsInt:        -1,
	IsUint:       1,
	IsString:     "string",
	IsFloat:      1.1,
	IsBool:       true,
	PtrIsInt:     &isInt,
	PtrIsInt8:    func() *int8 { c := int8(isInt); return &c }(),
	PtrIsInt16:   func() *int16 { c := int16(isInt); return &c }(),
	PtrIsInt32:   func() *int32 { c := int32(isInt); return &c }(),
	PtrIsInt64:   func() *int64 { c := int64(isInt); return &c }(),
	PtrIsUint:    &isUInt,
	PtrIsUint8:   func() *uint8 { c := uint8(isUInt); return &c }(),
	PtrIsUint16:  func() *uint16 { c := uint16(isUInt); return &c }(),
	PtrIsUint32:  func() *uint32 { c := uint32(isUInt); return &c }(),
	PtrIsUint64:  func() *uint64 { c := uint64(isUInt); return &c }(),
	PtrIsString:  &isString,
	PtrIsFloat32: func() *float32 { c := float32(isFloat); return &c }(),
	PtrIsFloat64: &isFloat,
	PtrIsBool:    &isBoolean,
}

var ExpectedTestNested = TestNested{
	IsString: "string",
	Basic: Basic{
		IsString: "string",
	},
	BasicPtr: &BasicPtr{
		IsString: &isString,
	},
}

type Basic struct {
	IsString string `env:"ENV_STRING"`
}

type BasicPtr struct {
	IsString *string `env:"ENV_STRING"`
}

type TestNested struct {
	IsString string `env:"ENV_STRING"`
	Basic    Basic
	BasicPtr *BasicPtr
}

func TestEnvTagParser_Parse(t *testing.T) {
	type fields struct {
		tagName string
	}
	type args struct {
		configuration interface{}
		deep          int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected interface{}
	}{
		{
			"Test differnt types",
			fields{tagName: "env"},
			args{&TestTypes{}, 0},
			&ExpectedTestTypes,
		},
		{
			"Test nested struct",
			fields{tagName: "env"},
			args{&TestNested{}, 0},
			&ExpectedTestNested,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setTestEnv(map[string]string{"ENV_UINT": "1", "ENV_STRING": "string", "ENV_INT": "-1", "ENV_FLOAT": "1.1", "ENV_BOOL": "true"})
			etp := NewEnvTagParser(tt.fields.tagName)
			etp.Parse(tt.args.configuration, tt.args.deep)

			if diff := deep.Equal(tt.args.configuration, tt.expected); diff != nil {
				t.Error(diff)
			}

		})
	}
}

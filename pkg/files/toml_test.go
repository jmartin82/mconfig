package files

import "testing"
import "github.com/go-test/deep"

func TestTOMUnmarshal(t *testing.T) {
	type args struct {
		source        []byte
		configuration interface{}
	}
	tests := []struct {
		name    string
		args    args
		result  *TestStruct
		wantErr bool
	}{
		{
			"Test valid format",
			args{
				[]byte("int = 1\nstring = \"x\"\n[nested]\nint = 2\nstring = \"xx\"\n"),
				&TestStruct{},
			},
			&TestStruct{Int: 1, String: "x", Nested: NestedStruct{Int: 2, String: "xx"}},
			false,
		},
		{
			"Test invalid format",
			args{
				[]byte("["),
				&TestStruct{},
			},
			&TestStruct{Int: 1, String: "x", Nested: NestedStruct{Int: 2, String: "xx"}},
			true,
		},
		{
			"Test empty input",
			args{
				[]byte(""),
				&TestStruct{},
			},
			&TestStruct{Int: 0, String: "", Nested: NestedStruct{Int: 0, String: ""}},
			true,
		},
		{
			"Test invalid stuct pointer",
			args{
				[]byte(""),
				TestStruct{},
			},
			&TestStruct{Int: 0, String: "", Nested: NestedStruct{Int: 0, String: ""}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			to := NewTOMLFormatUnmarshaler()
			if err := to.Unmarshal(tt.args.source, tt.args.configuration); (err != nil) != tt.wantErr {
				t.Errorf("TOMLFormat.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if diff := deep.Equal(tt.args.configuration, tt.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

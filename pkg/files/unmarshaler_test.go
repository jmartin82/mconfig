package files

import "testing"

type NestedStruct struct {
	Int    int
	String string
}

type TestStruct struct {
	Int    int
	String string
	Nested NestedStruct
}

func TestCapableFormat_IsCapable(t *testing.T) {
	type fields struct {
		formats []string
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"Test found",
			fields{
				formats: []string{".json", ".yml"},
			},
			args{
				filename: "test.yml",
			},
			true,
		},
		{
			"Test found",
			fields{
				formats: []string{".json", ".yml"},
			},
			args{
				filename: "test.xml",
			},
			false,
		},
		{
			"Test empty config and input",
			fields{
				formats: []string{},
			},
			args{
				filename: "",
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cf := CapableFormat{
				formats: tt.fields.formats,
			}
			if got := cf.IsCapable(tt.args.filename); got != tt.want {
				t.Errorf("CapableFormat.IsCapable() = %v, want %v", got, tt.want)
			}
		})
	}
}

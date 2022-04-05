package eco

import (
	"reflect"
	"strings"
	"testing"
)

type SampleStruct1 struct {
	Foo string `default:"Bar"`
	Baz int    `default:"100"`
}

type SampleStruct2 struct {
	Foo string        `default:"Bar"`
	Baz int           `default:"100"`
	Sub SampleStruct3 `env:"Sub"`
}

type SampleStruct3 struct {
	Foo string `default:"Baz"`
}

type SampleComplexStruct struct {
	Foo  string `default:"Bar"`
	Sub1 SampleComplexStruct_Sub
}

type SampleComplexStruct_Sub struct {
	String     string  `default:"foo"`
	Float      float64 `default:"1.1"`
	FloatEmpty float64
	Sub1       SampleComplexStruct_Sub_Sub
}

type SampleComplexStruct_Sub_Sub struct {
	I64 int64
}

type SampleComplexStruct3 struct {
	Foo  string
	Sub1 SampleComplexStruct3_Sub
}

type SampleComplexStruct3_Sub struct {
	String     string
	Float      float64
	FloatEmpty float64
	Sub1       SampleComplexStruct3_Sub_Sub
	Sub2       *SampleComplexStruct3_Sub_Sub
}

type SampleComplexStruct3_Sub_Sub struct {
	I64 int64
}

type SampleComplexStruct4 struct {
	Foo  string
	Sub1 SampleComplexStruct4_Sub
	Sub2 *SampleComplexStruct4_Sub
}

type SampleComplexStruct4_Sub struct {
	String     string
	Float      float64
	FloatEmpty float64
	Sub1       SampleComplexStruct4_Sub_Sub
}

type SampleComplexStruct4_Sub_Sub struct {
	I64 int64
}

type SampleArrayStruct struct {
	Foo []string `default:"foo,bar,baz"`
}

func TestEco_SetPrefix(t *testing.T) {
	tests := []struct {
		name    string
		prefix  string
		envs    map[string]string
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "should be same when no prefix",
			args: &SampleComplexStruct{},
			envs: map[string]string{
				"SUB1_STRING": "bar",
				"SUB1_FLOAT":  "2.2",
			},
			want: &SampleComplexStruct{
				Foo: "Bar",
				Sub1: SampleComplexStruct_Sub{
					String:     "bar",
					Float:      2.2,
					FloatEmpty: 0,
					Sub1: SampleComplexStruct_Sub_Sub{
						I64: 0,
					},
				},
			},
		},
		{
			name:   "should be same when with prefix",
			prefix: "prefix_",
			args:   &SampleComplexStruct{},
			envs: map[string]string{
				"PREFIX_SUB1_STRING":   "bar",
				"PREFIX_SUB1_FLOAT":    "2.2",
				"PREFIX_SUB1_SUB1_I64": "1",
			},
			want: &SampleComplexStruct{
				Foo: "Bar",
				Sub1: SampleComplexStruct_Sub{
					String:     "bar",
					Float:      2.2,
					FloatEmpty: 0,
					Sub1: SampleComplexStruct_Sub_Sub{
						I64: 1,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New().SetPrefix(tt.prefix)
			m := tt.args

			if e.envNamePrefix != tt.prefix {
				t.Errorf("Eco.SetPrefix() prefix = %v, want %v", e.envNamePrefix, tt.prefix)
			}

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestEco_SetPrefixAutoTrim(t *testing.T) {
	tests := []struct {
		name     string
		autoTrim bool
		prefix   string
		want     string
	}{
		{
			name:     "should trim prefix",
			autoTrim: true,
			prefix:   "prefix_",
			want:     "prefix",
		},
		{
			name:     "should not trim prefix",
			autoTrim: false,
			prefix:   "prefix_",
			want:     "prefix_",
		},
		{
			name:     "should trim from left for spaces even autoTrim is false",
			autoTrim: false,
			prefix:   "  prefix_",
			want:     "prefix_",
		},
		{
			name:     "should trim for right for spaces even autoTrim is false",
			autoTrim: false,
			prefix:   " prefix_ ",
			want:     "prefix_",
		},
		{
			name:     "should trim from left for spaces when autoTrim is true",
			autoTrim: true,
			prefix:   "  prefix_",
			want:     "prefix",
		},
		{
			name:     "should trim for right for spaces when autoTrim is true",
			autoTrim: true,
			prefix:   " prefix_  ",
			want:     "prefix",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New().SetPrefixAutoTrim(tt.autoTrim)
			if got := e.SetPrefix(tt.prefix).getPrefix(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Eco.SetPrefixAutoTrim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEco_SetArraySeparator(t *testing.T) {
	tests := []struct {
		name      string
		envs      map[string]string
		separator string
		args      interface{}
		want      interface{}
		wantErr   bool
	}{
		{
			name:      "default separator",
			separator: ",",
			args: &struct {
				Foo []string `default:"foo,bar,baz"`
			}{},
			want: &struct {
				Foo []string `default:"foo,bar,baz"`
			}{
				Foo: []string{"foo", "bar", "baz"},
			},
		},
		{
			name:      "custom separator",
			separator: "|",
			args: &struct {
				Foo []string `default:"foo|bar|baz"`
			}{},
			want: &struct {
				Foo []string `default:"foo|bar|baz"`
			}{
				Foo: []string{"foo", "bar", "baz"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New().SetArraySeparator(tt.separator)
			m := tt.args

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestEco_SetEnvNameTransformer(t *testing.T) {
	tests := []struct {
		name        string
		envs        map[string]string
		transformer envNameTransformerFunc
		args        interface{}
		want        interface{}
		wantErr     bool
	}{
		{
			name:        "default separator",
			transformer: defaultEnvNameTransformerFunc,
			args: &struct {
				Foo string
			}{},
			want: &struct {
				Foo string
			}{
				Foo: "foo",
			},
			envs: map[string]string{
				"FOO": "foo",
			},
		},
		{
			name: "custom transformer",
			transformer: func(parts []string, sep string) string {
				return "CUSTOM" + sep + strings.ToUpper(strings.Join(parts, sep))
			},
			args: &struct {
				Foo string
			}{},
			want: &struct {
				Foo string
			}{
				Foo: "foo",
			},
			envs: map[string]string{
				"CUSTOM_FOO": "foo",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New().SetEnvNameTransformer(tt.transformer)
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestEco_SetEnvNameSeparator(t *testing.T) {
	tests := []struct {
		name      string
		envs      map[string]string
		separator string
		args      interface{}
		want      interface{}
		wantErr   bool
	}{
		{
			name: "default separator",
			args: &struct {
				Foo struct {
					Bar string
				}
			}{},
			want: &struct {
				Foo struct {
					Bar string
				}
			}{
				Foo: struct {
					Bar string
				}{
					Bar: "bar",
				},
			},
			envs: map[string]string{
				"FOO_BAR": "bar",
			},
		},
		{
			name:      "custom separator",
			separator: ".",
			args: &struct {
				Foo struct {
					Bar string
				}
			}{},
			want: &struct {
				Foo struct {
					Bar string
				}
			}{
				Foo: struct {
					Bar string
				}{
					Bar: "bar",
				},
			},
			envs: map[string]string{
				"FOO.BAR": "bar",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New().SetEnvNameSeparator(tt.separator)
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestEco_SetValueGetter(t *testing.T) {
	tests := []struct {
		name    string
		envs    map[string]string
		getter  envValueGetterFunc
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "default getter",
			args: &struct {
				Foo struct {
					Bar string
				}
			}{},
			want: &struct {
				Foo struct {
					Bar string
				}
			}{
				Foo: struct {
					Bar string
				}{
					Bar: "bar",
				},
			},
			envs: map[string]string{
				"FOO_BAR": "bar",
			},
		},
		{
			name: "custom getter",
			getter: func(key string) string {
				switch key {
				case "FOO_BAR":
					return "custom_bar"
				default:
					return ""
				}
			},
			args: &struct {
				Foo struct {
					Bar string
				}
			}{},
			want: &struct {
				Foo struct {
					Bar string
				}
			}{
				Foo: struct {
					Bar string
				}{
					Bar: "custom_bar",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New().SetValueGetter(tt.getter)
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestEco_Unmarshal(t *testing.T) {
	sStr := "Bar"
	ss := struct {
		Name string
	}{}
	sas := &SampleArrayStruct{}
	ss4 := &SampleComplexStruct4_Sub{
		String: "custom_string",
	}

	var ss2 *SampleArrayStruct

	tests := []struct {
		name    string
		envs    map[string]string
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name:    "should error if argument is nil",
			args:    nil,
			wantErr: true,
		},
		{
			name: "should error if argument is not a pointer",
			args: struct {
				Name string
			}{},
			wantErr: true,
		},
		{
			name:    "should error if argument is not a pointer",
			args:    *sas,
			wantErr: true,
		},
		{
			name:    "should error if argument value is not a pointer",
			args:    ss2,
			wantErr: true,
		},
		{
			name: "should error if sub struct has any error",
			args: &SampleComplexStruct3{},
			want: &SampleComplexStruct3{
				Sub1: SampleComplexStruct3_Sub{
					Sub1: SampleComplexStruct3_Sub_Sub{},
					Sub2: &SampleComplexStruct3_Sub_Sub{},
				},
			},
			envs: map[string]string{
				"SUB1_SUB1_I64": "sub1",
				"SUB1_SUB2_I64": "sub2",
			},
			wantErr: true,
		},
		{
			name: "should continue to unmarshal if sub field is a struct",
			args: &SampleComplexStruct4{
				Sub2: ss4,
			},
			want: &SampleComplexStruct4{
				Sub1: SampleComplexStruct4_Sub{
					Sub1: SampleComplexStruct4_Sub_Sub{},
				},
				Sub2: ss4,
			},
			envs: map[string]string{
				"SUB2_STRING": "custom_string",
			},
		},
		{
			name: "should error if sub struct is not a pointer",
			args: &SampleComplexStruct3{},
			want: &SampleComplexStruct3{
				Sub1: SampleComplexStruct3_Sub{},
			},
			envs: map[string]string{
				"SUB1_SUB2_I64": "sub2",
			},
			wantErr: true,
		},
		{
			name: "should return same struct if has no default or env value",
			args: &ss,
			want: &ss,
		},
		{
			name: "should bind struct with env vars",
			args: &struct {
				Name string
			}{},
			want: &struct {
				Name string
			}{
				Name: "Foo",
			},
			envs: map[string]string{
				"NAME": "Foo",
			},
		},
		{
			name: "should bind struct with default values",
			args: &struct {
				Foo string `default:"Bar"`
				Baz int    `default:"100"`
			}{},
			want: &struct {
				Foo string `default:"Bar"`
				Baz int    `default:"100"`
			}{
				Foo: "Bar",
				Baz: 100,
			},
		},
		{
			name: "should bind TestStuct1 with default values",
			args: &SampleStruct1{},
			want: &SampleStruct1{
				Foo: "Bar",
				Baz: 100,
			},
		},
		{
			name: "should bind TestStuct1 with default value and env vars",
			args: &SampleStruct1{},
			envs: map[string]string{
				"FOO": "Foo",
			},
			want: &SampleStruct1{
				Foo: "Foo",
				Baz: 100,
			},
		},
		{
			name: "should bind TestStuct2 with default values",
			args: &SampleStruct2{},
			want: &SampleStruct2{
				Foo: "Bar",
				Baz: 100,
				Sub: SampleStruct3{
					Foo: "Baz",
				},
			},
		},
		{
			name: "should bind TestStuct2 with default values and env vars",
			args: &SampleStruct2{},
			envs: map[string]string{
				"SUB_FOO": "Foo",
			},
			want: &SampleStruct2{
				Foo: "Bar",
				Baz: 100,
				Sub: SampleStruct3{
					Foo: "Foo",
				},
			},
		},
		{
			name: "should bind struct with default values when a field is pointer",
			args: &struct {
				Foo *string `default:"Bar"`
				Baz bool    `default:"1"`
			}{},
			want: &struct {
				Foo *string `default:"Bar"`
				Baz bool    `default:"1"`
			}{
				Foo: &sStr,
				Baz: true,
			},
		},
		{
			name: "should bind struct with default values when a field is pointer and env vars",
			args: &SampleComplexStruct{},
			want: &SampleComplexStruct{
				Foo: "Bar",
				Sub1: SampleComplexStruct_Sub{
					String:     "foo",
					Float:      1.1,
					FloatEmpty: 0,
					Sub1: SampleComplexStruct_Sub_Sub{
						I64: 0,
					},
				},
			},
		},
		{
			name: "should bind struct with default values when a field is pointer",
			args: &SampleComplexStruct{},
			envs: map[string]string{
				"SUB1_STRING": "bar",
				"SUB1_FLOAT":  "2.2",
			},
			want: &SampleComplexStruct{
				Foo: "Bar",
				Sub1: SampleComplexStruct_Sub{
					String:     "bar",
					Float:      2.2,
					FloatEmpty: 0,
					Sub1: SampleComplexStruct_Sub_Sub{
						I64: 0,
					},
				},
			},
		},
		{
			name: "should bind array fields from defaults",
			args: &SampleArrayStruct{},
			envs: map[string]string{},
			want: &SampleArrayStruct{
				Foo: []string{"foo", "bar", "baz"},
			},
		},
		{
			name: "should bind array fields from env vars",
			args: &SampleArrayStruct{},
			envs: map[string]string{
				"FOO": "a,b,c",
			},
			want: &SampleArrayStruct{
				Foo: []string{"a", "b", "c"},
			},
		},
		{
			name: "should bind array fields from env vars 2",
			args: &SampleArrayStruct{},
			envs: map[string]string{
				"FOO": "string 1,string 2",
			},
			want: &SampleArrayStruct{
				Foo: []string{"string 1", "string 2"},
			},
		},
		{
			name: "should pass when a field is not exported",
			args: &struct {
				notExported string
			}{},
			want: &struct {
				notExported string
			}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr != true && tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestEcoUnmarshal_convertStrToFieldVal_String(t *testing.T) {
	str := "test"

	type Struct struct {
		FieldBlank      string
		FieldEnv        string
		FieldPtr        *string
		FieldPtrWithVal *string
		FieldDef        string `default:"foo"`
		FieldNumericDef string `default:"1"`
	}

	tests := []struct {
		name    string
		envs    map[string]string
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "should handle all fields",
			args: &Struct{
				FieldPtrWithVal: &str,
			},
			want: &Struct{
				FieldBlank:      "",
				FieldEnv:        "custom",
				FieldPtr:        nil,
				FieldPtrWithVal: &str,
				FieldDef:        "foo",
				FieldNumericDef: "1",
			},
			envs: map[string]string{
				"FIELD_ENV": "custom",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %#+v, want %#+v", m, tt.want)
			}
		})
	}
}

func TestEcoUnmarshal_convertStrToFieldVal_Int(t *testing.T) {
	i := 99

	type Struct struct {
		FieldBlank      int
		FieldEnv        int
		FieldPtr        *int
		FieldPtrWithVal *int
		FieldDef        int `default:"2"`
		FieldNumericDef int `default:"3"`
	}

	tests := []struct {
		name    string
		envs    map[string]string
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "should handle all fields",
			args: &Struct{
				FieldPtrWithVal: &i,
			},
			want: &Struct{
				FieldBlank:      0,
				FieldEnv:        1,
				FieldPtr:        nil,
				FieldPtrWithVal: &i,
				FieldDef:        2,
				FieldNumericDef: 3,
			},
			envs: map[string]string{
				"FIELD_ENV": "1",
			},
		},
		{
			name: "should handle all fields and negative values",
			args: &Struct{},
			want: &Struct{
				FieldBlank:      0,
				FieldEnv:        -1,
				FieldPtr:        nil,
				FieldPtrWithVal: nil,
				FieldDef:        2,
				FieldNumericDef: 3,
			},
			envs: map[string]string{
				"FIELD_ENV": "-1",
			},
		},
		{
			name: "should error when invalid value with env",
			args: &Struct{},
			want: &Struct{
				FieldEnv: 1,
			},
			envs: map[string]string{
				"FIELD_ENV": "aa",
			},
			wantErr: true,
		},
		{
			name: "should error when invalid value with default",
			args: &struct {
				FieldDef int `default:"aa"`
			}{},
			want: &struct {
				FieldDef int `default:"aa"`
			}{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %#+v, wantErr %#+v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %#+v, want %#+v", m, tt.want)
			}
		})
	}
}

func TestEcoUnmarshal_convertStrToFieldVal_Uint(t *testing.T) {

	i := uint(99)

	type Struct struct {
		FieldBlank      uint
		FieldEnv        uint
		FieldPtr        *uint
		FieldPtrWithVal *uint
		FieldDef        uint `default:"2"`
		FieldNumericDef uint `default:"3"`
	}

	tests := []struct {
		name    string
		envs    map[string]string
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "should handle all fields",
			args: &Struct{
				FieldPtrWithVal: &i,
			},
			want: &Struct{
				FieldBlank:      0,
				FieldEnv:        1,
				FieldPtr:        nil,
				FieldPtrWithVal: &i,
				FieldDef:        2,
				FieldNumericDef: 3,
			},
			envs: map[string]string{
				"FIELD_ENV": "1",
			},
		},
		{
			name: "should error when invalid value with env",
			args: &struct {
				FieldEnv uint
			}{},
			want: &struct {
				FieldEnv uint
			}{},
			envs: map[string]string{
				"FIELD_ENV": "aa",
			},
			wantErr: true,
		},
		{
			name: "should error when invalid value with default",
			args: &struct {
				FieldDef uint `default:"aa"`
			}{},
			want: &struct {
				FieldDef uint `default:"aa"`
			}{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestEcoUnmarshal_convertStrToFieldVal_Int64(t *testing.T) {

	i := int64(99)

	type Struct struct {
		FieldBlank      int64
		FieldEnv        int64
		FieldPtr        *int64
		FieldPtrWithVal *int64
		FieldDef        int64 `default:"2"`
		FieldNumericDef int64 `default:"3"`
	}

	tests := []struct {
		name    string
		envs    map[string]string
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "should handle all fields",
			args: &Struct{
				FieldPtrWithVal: &i,
			},
			want: &Struct{
				FieldBlank:      0,
				FieldEnv:        1,
				FieldPtr:        nil,
				FieldPtrWithVal: &i,
				FieldDef:        2,
				FieldNumericDef: 3,
			},
			envs: map[string]string{
				"FIELD_ENV": "1",
			},
		},
		{
			name: "should error when invalid value with env",
			args: &struct {
				FieldEnv int64
			}{},
			want: &struct {
				FieldEnv int64
			}{},
			envs: map[string]string{
				"FIELD_ENV": "aa",
			},
			wantErr: true,
		},
		{
			name: "should error when invalid value with default",
			args: &struct {
				FieldDef int64 `default:"aa"`
			}{},
			want: &struct {
				FieldDef int64 `default:"aa"`
			}{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestEcoUnmarshal_convertStrToFieldVal_Uint64(t *testing.T) {

	i := uint64(99)

	type Struct struct {
		FieldBlank      uint64
		FieldEnv        uint64
		FieldPtr        *uint64
		FieldPtrWithVal *uint64
		FieldDef        uint64 `default:"2"`
		FieldNumericDef uint64 `default:"3"`
	}

	tests := []struct {
		name    string
		envs    map[string]string
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "should handle all fields",
			args: &Struct{
				FieldPtrWithVal: &i,
			},
			want: &Struct{
				FieldBlank:      0,
				FieldEnv:        1,
				FieldPtr:        nil,
				FieldPtrWithVal: &i,
				FieldDef:        2,
				FieldNumericDef: 3,
			},
			envs: map[string]string{
				"FIELD_ENV": "1",
			},
		},
		{
			name: "should error when invalid value with env",
			args: &struct {
				FieldEnv uint64
			}{},
			want: &struct {
				FieldEnv uint64
			}{},
			envs: map[string]string{
				"FIELD_ENV": "aa",
			},
			wantErr: true,
		},
		{
			name: "should error when invalid value with default",
			args: &struct {
				FieldDef uint64 `default:"aa"`
			}{},
			want: &struct {
				FieldDef uint64 `default:"aa"`
			}{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestEcoUnmarshal_convertStrToFieldVal_Int32(t *testing.T) {

	i := int32(99)

	type Struct struct {
		FieldBlank      int32
		FieldEnv        int32
		FieldPtr        *int32
		FieldPtrWithVal *int32
		FieldDef        int32 `default:"2"`
		FieldNumericDef int32 `default:"3"`
	}

	tests := []struct {
		name    string
		envs    map[string]string
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "should handle all fields",
			args: &Struct{
				FieldPtrWithVal: &i,
			},
			want: &Struct{
				FieldBlank:      0,
				FieldEnv:        1,
				FieldPtr:        nil,
				FieldPtrWithVal: &i,
				FieldDef:        2,
				FieldNumericDef: 3,
			},
			envs: map[string]string{
				"FIELD_ENV": "1",
			},
		},
		{
			name: "should error when invalid value with env",
			args: &struct {
				FieldEnv int32
			}{},
			want: &struct {
				FieldEnv int32
			}{},
			envs: map[string]string{
				"FIELD_ENV": "aa",
			},
			wantErr: true,
		},
		{
			name: "should error when invalid value with default",
			args: &struct {
				FieldDef int32 `default:"aa"`
			}{},
			want: &struct {
				FieldDef int32 `default:"aa"`
			}{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestEcoUnmarshal_convertStrToFieldVal_Uint32(t *testing.T) {

	i := uint32(99)

	type Struct struct {
		FieldBlank      uint32
		FieldEnv        uint32
		FieldPtr        *uint32
		FieldPtrWithVal *uint32
		FieldDef        uint32 `default:"2"`
		FieldNumericDef uint32 `default:"3"`
	}

	tests := []struct {
		name    string
		envs    map[string]string
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "should handle all fields",
			args: &Struct{
				FieldPtrWithVal: &i,
			},
			want: &Struct{
				FieldBlank:      0,
				FieldEnv:        1,
				FieldPtr:        nil,
				FieldPtrWithVal: &i,
				FieldDef:        2,
				FieldNumericDef: 3,
			},
			envs: map[string]string{
				"FIELD_ENV": "1",
			},
		},
		{
			name: "should error when invalid value with env",
			args: &struct {
				FieldEnv uint32
			}{},
			want: &struct {
				FieldEnv uint32
			}{},
			envs: map[string]string{
				"FIELD_ENV": "aa",
			},
			wantErr: true,
		},
		{
			name: "should error when invalid value with default",
			args: &struct {
				FieldDef uint32 `default:"aa"`
			}{},
			want: &struct {
				FieldDef uint32 `default:"aa"`
			}{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestEcoUnmarshal_convertStrToFieldVal_Int16(t *testing.T) {

	i := int16(99)

	type Struct struct {
		FieldBlank      int16
		FieldEnv        int16
		FieldPtr        *int16
		FieldPtrWithVal *int16
		FieldDef        int16 `default:"2"`
		FieldNumericDef int16 `default:"3"`
	}

	tests := []struct {
		name    string
		envs    map[string]string
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "should handle all fields",
			args: &Struct{
				FieldPtrWithVal: &i,
			},
			want: &Struct{
				FieldBlank:      0,
				FieldEnv:        1,
				FieldPtr:        nil,
				FieldPtrWithVal: &i,
				FieldDef:        2,
				FieldNumericDef: 3,
			},
			envs: map[string]string{
				"FIELD_ENV": "1",
			},
		},
		{
			name: "should error when invalid value with env",
			args: &struct {
				FieldEnv int16
			}{},
			want: &struct {
				FieldEnv int16
			}{},
			envs: map[string]string{
				"FIELD_ENV": "aa",
			},
			wantErr: true,
		},
		{
			name: "should error when invalid value with default",
			args: &struct {
				FieldDef int16 `default:"aa"`
			}{},
			want: &struct {
				FieldDef int16 `default:"aa"`
			}{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestEcoUnmarshal_convertStrToFieldVal_Uint16(t *testing.T) {

	i := uint16(99)

	type Struct struct {
		FieldBlank      uint16
		FieldEnv        uint16
		FieldPtr        *uint16
		FieldPtrWithVal *uint16
		FieldDef        uint16 `default:"2"`
		FieldNumericDef uint16 `default:"3"`
	}

	tests := []struct {
		name    string
		envs    map[string]string
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "should handle all fields",
			args: &Struct{
				FieldPtrWithVal: &i,
			},
			want: &Struct{
				FieldBlank:      0,
				FieldEnv:        1,
				FieldPtr:        nil,
				FieldPtrWithVal: &i,
				FieldDef:        2,
				FieldNumericDef: 3,
			},
			envs: map[string]string{
				"FIELD_ENV": "1",
			},
		},
		{
			name: "should error when invalid value with env",
			args: &struct {
				FieldEnv uint16
			}{},
			want: &struct {
				FieldEnv uint16
			}{},
			envs: map[string]string{
				"FIELD_ENV": "aa",
			},
			wantErr: true,
		},
		{
			name: "should error when invalid value with default",
			args: &struct {
				FieldDef uint16 `default:"aa"`
			}{},
			want: &struct {
				FieldDef uint16 `default:"aa"`
			}{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestEcoUnmarshal_convertStrToFieldVal_Int8(t *testing.T) {

	i := int8(99)

	type Struct struct {
		FieldBlank      int8
		FieldEnv        int8
		FieldPtr        *int8
		FieldPtrWithVal *int8
		FieldDef        int8 `default:"2"`
		FieldNumericDef int8 `default:"3"`
	}

	tests := []struct {
		name    string
		envs    map[string]string
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "should handle all fields",
			args: &Struct{
				FieldPtrWithVal: &i,
			},
			want: &Struct{
				FieldBlank:      0,
				FieldEnv:        1,
				FieldPtr:        nil,
				FieldPtrWithVal: &i,
				FieldDef:        2,
				FieldNumericDef: 3,
			},
			envs: map[string]string{
				"FIELD_ENV": "1",
			},
		},
		{
			name: "should error when invalid value with env",
			args: &struct {
				FieldEnv int8
			}{},
			want: &struct {
				FieldEnv int8
			}{},
			envs: map[string]string{
				"FIELD_ENV": "aa",
			},
			wantErr: true,
		},
		{
			name: "should error when invalid value with default",
			args: &struct {
				FieldDef int8 `default:"aa"`
			}{},
			want: &struct {
				FieldDef int8 `default:"aa"`
			}{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestEcoUnmarshal_convertStrToFieldVal_Uint8(t *testing.T) {

	i := uint8(99)

	type Struct struct {
		FieldBlank      uint8
		FieldEnv        uint8
		FieldPtr        *uint8
		FieldPtrWithVal *uint8
		FieldDef        uint8 `default:"2"`
		FieldNumericDef uint8 `default:"3"`
	}

	tests := []struct {
		name    string
		envs    map[string]string
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "should handle all fields",
			args: &Struct{
				FieldPtrWithVal: &i,
			},
			want: &Struct{
				FieldBlank:      0,
				FieldEnv:        1,
				FieldPtr:        nil,
				FieldPtrWithVal: &i,
				FieldDef:        2,
				FieldNumericDef: 3,
			},
			envs: map[string]string{
				"FIELD_ENV": "1",
			},
		},
		{
			name: "should error when invalid value with env",
			args: &struct {
				FieldEnv uint8
			}{},
			want: &struct {
				FieldEnv uint8
			}{},
			envs: map[string]string{
				"FIELD_ENV": "aa",
			},
			wantErr: true,
		},
		{name: "should error when invalid value with default",
			args: &struct {
				FieldDef uint8 `default:"aa"`
			}{},
			want: &struct {
				FieldDef uint8 `default:"aa"`
			}{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestEcoUnmarshal_convertStrToFieldVal_Float32(t *testing.T) {
	i := float32(99.99)

	type Struct struct {
		FieldBlank      float32
		FieldEnv        float32
		FieldPtr        *float32
		FieldPtrWithVal *float32
		FieldDef        float32 `default:"2.1"`
		FieldNumericDef float32 `default:"3"`
	}

	tests := []struct {
		name    string
		envs    map[string]string
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "should handle all fields",
			args: &Struct{
				FieldPtrWithVal: &i,
			},
			want: &Struct{
				FieldBlank:      0,
				FieldEnv:        1.1,
				FieldPtr:        nil,
				FieldPtrWithVal: &i,
				FieldDef:        2.1,
				FieldNumericDef: 3,
			},
			envs: map[string]string{
				"FIELD_ENV": "1.1",
			},
		},
		{
			name: "should handle all fields and negative values",
			args: &Struct{},
			want: &Struct{
				FieldBlank:      0,
				FieldEnv:        -1.9,
				FieldPtr:        nil,
				FieldPtrWithVal: nil,
				FieldDef:        2.1,
				FieldNumericDef: 3,
			},
			envs: map[string]string{
				"FIELD_ENV": "-1.9",
			},
		},
		{
			name: "should error when invalid value with env",
			args: &Struct{},
			want: &Struct{
				FieldEnv: 1,
			},
			envs: map[string]string{
				"FIELD_ENV": "aa",
			},
			wantErr: true,
		},
		{
			name: "should error when invalid value with default",
			args: &struct {
				FieldDef float32 `default:"aa"`
			}{},
			want: &struct {
				FieldDef float32 `default:"aa"`
			}{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %#+v, wantErr %#+v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %#+v, want %#+v", m, tt.want)
			}
		})
	}
}

func TestEcoUnmarshal_convertStrToFieldVal_Float64(t *testing.T) {
	i := float64(-99.99)

	type Struct struct {
		FieldBlank      float64
		FieldEnv        float64
		FieldPtr        *float64
		FieldPtrWithVal *float64
		FieldDef        float64 `default:"2.1"`
		FieldNumericDef float64 `default:"3"`
	}

	tests := []struct {
		name    string
		envs    map[string]string
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "should handle all fields",
			args: &Struct{
				FieldPtrWithVal: &i,
			},
			want: &Struct{
				FieldBlank:      0,
				FieldEnv:        1.1,
				FieldPtr:        nil,
				FieldPtrWithVal: &i,
				FieldDef:        2.1,
				FieldNumericDef: 3,
			},
			envs: map[string]string{
				"FIELD_ENV": "1.1",
			},
		},
		{
			name: "should handle all fields and negative values",
			args: &Struct{},
			want: &Struct{
				FieldBlank:      0,
				FieldEnv:        -1.9,
				FieldPtr:        nil,
				FieldPtrWithVal: nil,
				FieldDef:        2.1,
				FieldNumericDef: 3,
			},
			envs: map[string]string{
				"FIELD_ENV": "-1.9",
			},
		},
		{
			name: "should error when invalid value with env",
			args: &Struct{},
			want: &Struct{
				FieldEnv: 1,
			},
			envs: map[string]string{
				"FIELD_ENV": "aa",
			},
			wantErr: true,
		},
		{
			name: "should error when invalid value with default",
			args: &struct {
				FieldDef float64 `default:"aa"`
			}{},
			want: &struct {
				FieldDef float64 `default:"aa"`
			}{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %#+v, wantErr %#+v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %#+v, want %#+v", m, tt.want)
			}
		})
	}
}

func TestEcoUnmarshal_convertStrToFieldVal_Bool(t *testing.T) {
	b := true

	type Struct struct {
		FieldBlank      bool
		FieldEnv        bool
		FieldPtr        *bool
		FieldPtrWithVal *bool
		FieldDef        bool `default:"true"`
		FieldNumericDef bool `default:"1"`
	}

	tests := []struct {
		name    string
		envs    map[string]string
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "should handle all fields",
			args: &Struct{
				FieldPtrWithVal: &b,
			},
			want: &Struct{
				FieldBlank:      false,
				FieldEnv:        true,
				FieldPtr:        nil,
				FieldPtrWithVal: &b,
				FieldDef:        true,
				FieldNumericDef: true,
			},
			envs: map[string]string{
				"FIELD_ENV": "1",
			},
		},
		{
			name: "should error when invalid value with env",
			args: &Struct{},
			want: &Struct{},
			envs: map[string]string{
				"FIELD_ENV": "aa",
			},
			wantErr: true,
		},
		{
			name: "should error when invalid value with default",
			args: &struct {
				FieldDef bool `default:"aa"`
			}{},
			want: &struct {
				FieldDef bool `default:"aa"`
			}{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %#+v, want %#+v", m, tt.want)
			}
		})
	}
}

func TestEcoUnmarshal_convertStrToFieldVal_SliceString(t *testing.T) {
	b := []string{"a", "b", "c"}

	type Struct struct {
		FieldBlank      []string
		FieldEnv        []string
		FieldPtr        *[]string
		FieldPtrWithVal *[]string
		FieldDef        []string `default:"d,e,f"`
		FieldNumericDef []string `default:"1,2,3"`
	}

	tests := []struct {
		name    string
		envs    map[string]string
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "should handle all fields",
			args: &Struct{
				FieldPtrWithVal: &b,
			},
			want: &Struct{
				FieldEnv:        []string{"x", "y", "z"},
				FieldPtr:        nil,
				FieldPtrWithVal: &b,
				FieldDef:        []string{"d", "e", "f"},
				FieldNumericDef: []string{"1", "2", "3"},
			},
			envs: map[string]string{
				"FIELD_ENV": "x,y,z",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %#+v, want %#+v", m, tt.want)
			}
		})
	}
}

func TestEcoUnmarshal_convertStrToFieldVal_SliceInt(t *testing.T) {
	b := []int{1, 2, 3}

	type Struct struct {
		FieldBlank      []int
		FieldEnv        []int
		FieldPtr        *[]int
		FieldPtrWithVal *[]int
		FieldDef        []int `default:"4,5,6"`
		FieldNumericDef []int `default:"1,2,3"`
	}

	tests := []struct {
		name    string
		envs    map[string]string
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "should handle all fields",
			args: &Struct{
				FieldPtrWithVal: &b,
			},
			want: &Struct{
				FieldEnv:        []int{-1, -2, -3},
				FieldPtr:        nil,
				FieldPtrWithVal: &b,
				FieldDef:        []int{4, 5, 6},
				FieldNumericDef: []int{1, 2, 3},
			},
			envs: map[string]string{
				"FIELD_ENV": "-1,-2,-3",
			},
		},
		{
			name: "should error when invalid value with env",
			args: &Struct{
				FieldPtrWithVal: &b,
			},
			want: &Struct{
				FieldEnv:        []int{-1, -2, -3},
				FieldPtr:        nil,
				FieldPtrWithVal: &b,
				FieldDef:        []int{4, 5, 6},
				FieldNumericDef: []int{1, 2, 3},
			},
			envs: map[string]string{
				"FIELD_ENV": "-1,a,-3",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %#+v, want %#+v", m, tt.want)
			}
		})
	}
}

func TestEcoUnmarshal_convertStrToFieldVal_SliceInt64(t *testing.T) {
	b := []int64{1, 2, 3}

	type Struct struct {
		FieldBlank      []int64
		FieldEnv        []int64
		FieldPtr        *[]int64
		FieldPtrWithVal *[]int64
		FieldDef        []int64 `default:"4,5,6"`
		FieldNumericDef []int64 `default:"1,2,3"`
	}

	tests := []struct {
		name    string
		envs    map[string]string
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "should handle all fields",
			args: &Struct{
				FieldPtrWithVal: &b,
			},
			want: &Struct{
				FieldEnv:        []int64{-1, -2, -3},
				FieldPtr:        nil,
				FieldPtrWithVal: &b,
				FieldDef:        []int64{4, 5, 6},
				FieldNumericDef: []int64{1, 2, 3},
			},
			envs: map[string]string{
				"FIELD_ENV": "-1,-2,-3",
			},
		},
		{
			name: "should error when invalid value with env",
			args: &Struct{
				FieldPtrWithVal: &b,
			},
			want: &Struct{
				FieldEnv:        []int64{-1, -2, -3},
				FieldPtr:        nil,
				FieldPtrWithVal: &b,
				FieldDef:        []int64{4, 5, 6},
				FieldNumericDef: []int64{1, 2, 3},
			},
			envs: map[string]string{
				"FIELD_ENV": "-1,a,-3",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %#+v, want %#+v", m, tt.want)
			}
		})
	}
}

func TestEcoUnmarshal_convertStrToFieldVal_SliceFloat32(t *testing.T) {
	b := []float32{-1.2, 2.5, 3}

	type Struct struct {
		FieldBlank      []float32
		FieldEnv        []float32
		FieldPtr        *[]float32
		FieldPtrWithVal *[]float32
		FieldDef        []float32 `default:"-5.5,-4,-3.2"`
		FieldNumericDef []float32 `default:"1,2,3"`
	}

	tests := []struct {
		name    string
		envs    map[string]string
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "should handle all fields",
			args: &Struct{
				FieldPtrWithVal: &b,
			},
			want: &Struct{
				FieldEnv:        []float32{-1.1, -2.0, 3.5},
				FieldPtr:        nil,
				FieldPtrWithVal: &b,
				FieldDef:        []float32{-5.5, -4, -3.2},
				FieldNumericDef: []float32{1, 2, 3},
			},
			envs: map[string]string{
				"FIELD_ENV": "-1.1, -2, 3.5",
			},
		},
		{
			name: "should error when invalid value with env",
			args: &Struct{
				FieldPtrWithVal: &b,
			},
			want: &Struct{
				FieldEnv:        []float32{-1.1, 2, 3.5},
				FieldPtr:        nil,
				FieldPtrWithVal: &b,
				FieldDef:        []float32{-5.5, -4, -3.2},
				FieldNumericDef: []float32{1, 2, 3},
			},
			envs: map[string]string{
				"FIELD_ENV": "-1.1,a,3.5",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %#+v, want %#+v", m, tt.want)
			}
		})
	}
}

func TestEcoUnmarshal_convertStrToFieldVal_SliceFloat64(t *testing.T) {
	b := []float64{-1.2, 2.5, 3}

	type Struct struct {
		FieldBlank      []float64
		FieldEnv        []float64
		FieldPtr        *[]float64
		FieldPtrWithVal *[]float64
		FieldDef        []float64 `default:"-5.5,-4,-3.2"`
		FieldNumericDef []float64 `default:"1,2,3"`
	}

	tests := []struct {
		name    string
		envs    map[string]string
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "should handle all fields",
			args: &Struct{
				FieldPtrWithVal: &b,
			},
			want: &Struct{
				FieldEnv:        []float64{-1.1, -2.0, 3.5},
				FieldPtr:        nil,
				FieldPtrWithVal: &b,
				FieldDef:        []float64{-5.5, -4, -3.2},
				FieldNumericDef: []float64{1, 2, 3},
			},
			envs: map[string]string{
				"FIELD_ENV": "-1.1, -2, 3.5",
			},
		},
		{
			name: "should error when invalid value with env",
			args: &Struct{
				FieldPtrWithVal: &b,
			},
			want: &Struct{
				FieldEnv:        []float64{-1.1, 2, 3.5},
				FieldPtr:        nil,
				FieldPtrWithVal: &b,
				FieldDef:        []float64{-5.5, -4, -3.2},
				FieldNumericDef: []float64{1, 2, 3},
			},
			envs: map[string]string{
				"FIELD_ENV": "-1.1,a,3.5",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %#+v, want %#+v", m, tt.want)
			}
		})
	}
}

func TestEcoUnmarshal_convertStrToFieldVal_SliceUnsupportedType(t *testing.T) {
	type Struct struct {
		FieldBlank []bool
		FieldEnv   []bool
	}

	tests := []struct {
		name    string
		envs    map[string]string
		args    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "should return unsupported type error",
			args: &Struct{},
			want: &Struct{
				FieldEnv: []bool{},
			},
			envs: map[string]string{
				"FIELD_ENV": "true, false",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := e.Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %#+v, want %#+v", m, tt.want)
			}
		})
	}
}

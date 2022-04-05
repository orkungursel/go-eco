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
	s := &SampleArrayStruct{}
	var ss2 *SampleArrayStruct
	ss4 := &SampleComplexStruct4_Sub{
		String: "custom_string",
	}

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
			args:    *s,
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
			args: &struct {
				Name string
			}{},
			want: &struct {
				Name string
			}{},
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

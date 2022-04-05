package eco

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestSetPrefix(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		want   string
	}{
		{
			name:   "set prefix",
			prefix: "prefix",
			want:   "prefix",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetPrefix(tt.prefix)
			if got := ee.getPrefix(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetPrefix() = %v, want %v", got, tt.want)
			}
			SetPrefix("")
		})
	}
}

func TestSetArraySeparator(t *testing.T) {
	tests := []struct {
		name      string
		separator string
		want      string
	}{
		{
			name:      "set separator",
			separator: ".",
			want:      ".",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetArraySeparator(tt.separator)
			if got := ee.sliceSeparator; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetArraySeparator() = %v, want %v", got, tt.want)
			}

			SetArraySeparator(".")
		})
	}
}

func TestSetEnvNameTransformer(t *testing.T) {
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
			SetEnvNameTransformer(tt.transformer)
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %v, want %v", m, tt.want)
			}

			SetEnvNameTransformer(defaultEnvNameTransformerFunc)
		})
	}
}

func TestSetEnvNameSeparator(t *testing.T) {
	tests := []struct {
		name      string
		separator string
		want      string
	}{
		{
			name:      "set separator",
			separator: ".",
			want:      ".",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetEnvNameSeparator(tt.separator)
			if got := ee.envNameSeparator; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetEnvNameSeparator() = %v, want %v", got, tt.want)
			}
			SetEnvNameSeparator("_")
		})
	}
}

func TestSetValueGetter(t *testing.T) {
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
			SetValueGetter(tt.getter)
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %v, want %v", m, tt.want)
			}

			SetValueGetter(os.Getenv)
		})
	}
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    interface{}
		envs    map[string]string
		want    interface{}
		wantErr bool
	}{
		{
			name:    "should error if argument is nil",
			args:    nil,
			wantErr: true,
		},
		{
			name: "should fill struct with default values",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.args

			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			if err := Unmarshal(m); (err != nil) != tt.wantErr {
				t.Errorf("Eco.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr != true && tt.want != nil && !reflect.DeepEqual(m, tt.want) {
				t.Errorf("Eco.Unmarshal() = %v, want %v", m, tt.want)
			}
		})
	}
}

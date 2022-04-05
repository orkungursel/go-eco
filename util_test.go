package eco

import (
	"reflect"
	"testing"
)

func Test_toSnakeCase(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should convert to snake case",
			args: args{
				str: "FooBar",
			},
			want: "foo_bar",
		},
		{
			name: "should convert to snake case",
			args: args{
				str: "_FooBar",
			},
			want: "__foo_bar",
		},
		{
			name: "should convert to snake case",
			args: args{
				str: "fooBar",
			},
			want: "foo_bar",
		},
		{
			name: "should convert to snake case",
			args: args{
				str: "fOoBar",
			},
			want: "f_oo_bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toSnakeCase(tt.args.str); got != tt.want {
				t.Errorf("toSnakeCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultEnvNameFunction(t *testing.T) {
	type Struct struct {
		FieldDefaultEnvTag                string
		FieldPtr                          *string
		FieldDef                          string `default:"foo"`
		FieldNumericDef                   string `default:"1"`
		FieldCustomEnvTag                 string `env:"custom_env_tag_name" default:"1"`
		FieldCustomEnvTagSnakeCase        string `default:"snake_case"`
		FieldCustomEnvTagSnakeCaseWithEnv string `env:"customEnvTagNameSnakeCaseWithEnv" default:"snake_case_with_env"`
		FieldCustomEnvTagUpperCaseWithEnv string `env:"CUSTOMENVTAGNAMEUPPERCASEWITHENV" default:"upper_case_with_env"`
		FieldCustomEnvTagWDef             string `env:"baz" default:"2"`
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
			args: &Struct{},
			want: &Struct{
				FieldDefaultEnvTag:                "custom_default",
				FieldPtr:                          nil,
				FieldDef:                          "foo",
				FieldNumericDef:                   "1",
				FieldCustomEnvTag:                 "custom",
				FieldCustomEnvTagSnakeCase:        "snake_case",
				FieldCustomEnvTagSnakeCaseWithEnv: "custom_value",
				FieldCustomEnvTagUpperCaseWithEnv: "custom_value2",
				FieldCustomEnvTagWDef:             "2",
			},
			envs: map[string]string{
				"FIELD_DEFAULT_ENV_TAG":                   "custom_default",
				"CUSTOM_ENV_TAG_NAME":                     "custom",
				"CUSTOM_ENV_TAG_NAME_SNAKE_CASE":          "custom_snake_case",
				"CUSTOM_ENV_TAG_NAME_SNAKE_CASE_WITH_ENV": "custom_value",
				"CUSTOMENVTAGNAMEUPPERCASEWITHENV":        "custom_value2",
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

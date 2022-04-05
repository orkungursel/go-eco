package eco

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type eco struct {
	sliceSeparator        string
	envNameSeparator      string
	envNamePrefix         string
	envNamePrefixAutoTrim bool
	envNameTransformer    envNameTransformerFunc
	envValueGetter        envValueGetterFunc
	tagNameEnv            string
	tagNameDefault        string
	tagSkipIdentifier     string
}

// New returns a new Eco instance.
func New() *eco {
	return &eco{
		sliceSeparator:        ",",
		envNameSeparator:      "_",
		envNamePrefixAutoTrim: true,
		envNameTransformer:    defaultEnvNameTransformerFunc,
		envValueGetter:        os.Getenv,
		tagNameEnv:            "env",
		tagNameDefault:        "default",
		tagSkipIdentifier:     "-",
	}
}

// SetPrefix sets the prefix for the environment variable names.
// It will be trimmed from both sides if it is not empty.
func (e *eco) SetPrefix(prefix string) *eco {
	e.envNamePrefix = strings.TrimSpace(prefix)
	return e
}

// getPrefix returns the prefix for the environment variable names.
// If auto trim is enabled, it will be trimmed from both sides.
func (e *eco) getPrefix() string {
	prefix := e.envNamePrefix
	if e.envNamePrefixAutoTrim {
		prefix = strings.TrimRight(prefix, e.envNameSeparator)
	}
	return prefix
}

// SetPrefixAutoTrim enables or disables auto trimming of the prefix.
// Default is true.
func (e *eco) SetPrefixAutoTrim(autoTrim bool) *eco {
	e.envNamePrefixAutoTrim = autoTrim
	return e
}

// SetArraySeparator sets the separator for array values.
// Default is ",".
func (e *eco) SetArraySeparator(arrSep string) *eco {
	if arrSep != "" {
		e.sliceSeparator = arrSep
	}
	return e
}

// SetEnvNameTransformer sets the function for transforming the environment variable names.
func (e *eco) SetEnvNameTransformer(transformerFunc envNameTransformerFunc) *eco {
	e.envNameTransformer = transformerFunc
	return e
}

// SetEnvNameSeparator sets the paths to the environment files.
func (e *eco) SetEnvNameSeparator(envNameSeparator string) *eco {
	if envNameSeparator != "" {
		e.envNameSeparator = envNameSeparator
	}
	return e
}

// SetValueGetter sets the function for getting the environment variable values.
func (e *eco) SetValueGetter(valueGetter envValueGetterFunc) *eco {
	if valueGetter != nil {
		e.envValueGetter = valueGetter
	}
	return e
}

// Unmarshal unmarshals the environment variables into the given struct.
func (e *eco) Unmarshal(m interface{}) error {
	if m == nil {
		return ErrRequiresNonNilPtr
	}

	t := reflect.TypeOf(m)
	if t.Kind() != reflect.Ptr {
		return ErrRequiresNonNilPtr
	}

	v := reflect.ValueOf(m)
	if v.IsNil() {
		return ErrRequiresNonNilPtr
	}

	var p []string
	if prefix := e.getPrefix(); prefix != "" {
		p = append(p, prefix)
	}

	return e.bindStructValues(m, p...)
}

// bindStructValues binds the environment variables to the given struct.
func (e *eco) bindStructValues(s interface{}, envNameParts ...string) error {
	sr := e.getStructReflection(s)

	for i := 0; i < sr.Type().NumField(); i++ {
		field := sr.Field(i)
		typeField := sr.Type().Field(i)

		// Skip unexported fields
		if !typeField.IsExported() {
			continue
		}

		tags := typeField.Tag
		isPtr := field.Type().Kind() == reflect.Ptr
		isStruct := typeField.Type.Kind() == reflect.Struct

		envTagValue, ok := tags.Lookup(e.tagNameEnv)
		if !ok || envTagValue == "" {
			// If field "env" tag is not provided, get tag from struct field name
			envTagValue = typeField.Name
		}

		envTagValue = toSnakeCase(envTagValue)

		p := envNameParts

		// if tag value is "-", skip this field
		// when looking for the environment variable name
		skip := envTagValue == e.tagSkipIdentifier
		if !skip {
			p = append(p, envTagValue)
		}

		// sanitize env variable name using the envNameFunc
		envKey := e.envNameTransformer(p, e.envNameSeparator)

		var envVal string
		var err error

		// get value from env
		envVal = e.envValueGetter(envKey)

		// if value is empty, get default value from tag
		if envVal == "" {
			defaultValue, ok := tags.Lookup(e.tagNameDefault)
			if ok {
				envVal = defaultValue
			}
		}

		// if field is a pointer, create
		if isPtr {
			if typeField.Type.Elem().Kind() != reflect.Struct && envVal == "" {
				continue
			}

			// if field is nil, create a new one with the type of the field
			if field.IsNil() {
				field.Set(reflect.New(field.Type().Elem()))
			}

			// check whether the field element kind
			// is a struct, then bind its values
			if field.Elem().Kind() == reflect.Struct {
				if err := e.bindStructValues(field.Interface(), p...); err != nil {
					return err
				}

				continue
			}
		} else if isStruct { // if field is a struct, bind it
			if err := e.bindStructValues(field.Addr().Interface(), p...); err != nil {
				return err
			}

			continue
		}

		// if value is empty, skip binding
		if envVal == "" {
			continue
		}

		// convert string value which comes from env to the type of the field
		val, err := e.convertStrToFieldVal(sr, i, envVal)
		if err != nil {
			return errors.New(err.Error() + ": " + envKey)
		}

		// set field value
		if isPtr {
			field.Elem().Set(val)
		} else {
			field.Set(val)
		}
	}

	return nil
}

// getStructReflection returns the reflection of the given struct.
func (e *eco) getStructReflection(s interface{}) reflect.Value {
	sv := reflect.ValueOf(s)

	if sv.Type().Kind() == reflect.Ptr {
		sv = sv.Elem()
	}

	return sv
}

// convertStrToFieldVal converts the given string value to the type of the field.
func (e *eco) convertStrToFieldVal(ref reflect.Value, index int, val interface{}) (out reflect.Value, err error) {
	field := ref.Field(index)
	typeField := ref.Type().Field(index)
	isPtr := typeField.Type.Kind() == reflect.Ptr

	kind := field.Kind()
	if isPtr {
		kind = field.Elem().Kind()
	}

	switch kind {
	case reflect.String:
		val = val.(string)
	case reflect.Int:
		val, err = strconv.Atoi(val.(string))
	case reflect.Uint:
		val, err = strconv.ParseUint(val.(string), 10, 0)
		val = uint(val.(uint64))
	case reflect.Int64:
		val, err = strconv.ParseInt(val.(string), 10, 64)
	case reflect.Uint64:
		val, err = strconv.ParseUint(val.(string), 10, 64)
	case reflect.Int32:
		val, err = strconv.ParseInt(val.(string), 10, 32)
		val = int32(val.(int64))
	case reflect.Uint32:
		val, err = strconv.ParseUint(val.(string), 10, 32)
		val = uint32(val.(uint64))
	case reflect.Int16:
		val, err = strconv.ParseInt(val.(string), 10, 16)
		val = int16(val.(int64))
	case reflect.Uint16:
		val, err = strconv.ParseUint(val.(string), 10, 16)
		val = uint16(val.(uint64))
	case reflect.Int8:
		val, err = strconv.ParseInt(val.(string), 10, 8)
		val = int8(val.(int64))
	case reflect.Uint8:
		val, err = strconv.ParseUint(val.(string), 10, 8)
		val = uint8(val.(uint64))
	case reflect.Float32:
		val, err = strconv.ParseFloat(val.(string), 32)
		val = float32(val.(float64))
	case reflect.Float64:
		val, err = strconv.ParseFloat(val.(string), 64)
	case reflect.Bool:
		val, err = strconv.ParseBool(val.(string))
	case reflect.Slice:
		switch typeField.Type.Elem().Kind() {
		case reflect.String:
			val = strings.Split(val.(string), e.sliceSeparator)
			vals := make([]string, len(val.([]string)))
			for i, v := range val.([]string) {
				vals[i] = strings.TrimSpace(v)
			}
			val = vals
		case reflect.Int:
			val = strings.Split(val.(string), e.sliceSeparator)
			var vals []int
			for _, v := range val.([]string) {
				i, err := strconv.Atoi(strings.TrimSpace(v))
				if err != nil {
					return reflect.ValueOf(nil), err
				}
				vals = append(vals, i)
			}
			val = vals
		case reflect.Int64:
			val = strings.Split(val.(string), e.sliceSeparator)
			var vals []int64
			for _, v := range val.([]string) {
				i, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
				if err != nil {
					return reflect.ValueOf(nil), err
				}
				vals = append(vals, i)
			}
			val = vals
		case reflect.Float32:
			val = strings.Split(val.(string), e.sliceSeparator)
			var vals []float32
			for _, v := range val.([]string) {
				f, err := strconv.ParseFloat(strings.TrimSpace(v), 32)
				if err != nil {
					return reflect.ValueOf(nil), err

				}
				vals = append(vals, float32(f))
			}
			val = vals
		case reflect.Float64:
			val = strings.Split(val.(string), e.sliceSeparator)
			var vals []float64
			for _, v := range val.([]string) {
				f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
				if err != nil {
					return reflect.ValueOf(nil), err
				}
				vals = append(vals, float64(f))
			}
			val = vals
		default:
			return reflect.ValueOf(nil), fmt.Errorf("unsupported slice type: %s", typeField.Type.Elem().Kind())
		}
	}

	return reflect.ValueOf(val), err
}

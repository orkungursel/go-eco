package eco

var ee *eco

func init() {
	ee = New()
}

// SetPrefix sets the prefix for environment variables.
func SetPrefix(prefix string) *eco {
	return ee.SetPrefix(prefix)
}

// SetArraySeparator sets the separator for array values.
// Default is ",".
func SetArraySeparator(sep string) *eco {
	return ee.SetArraySeparator(sep)
}

// SetEnvNameTransformer sets the function for transforming the environment variable names.
func SetEnvNameTransformer(transformerFunc envNameTransformerFunc) *eco {
	return ee.SetEnvNameTransformer(transformerFunc)
}

// SetEnvNameSeparator sets the separator for the environment variable names.
// Default is "_".
func SetEnvNameSeparator(envNameSeparator string) *eco {
	return ee.SetEnvNameSeparator(envNameSeparator)
}

// SetValueGetter sets the function for getting the environment variable values.
func SetValueGetter(valueGetter envValueGetterFunc) *eco {
	return ee.SetValueGetter(valueGetter)
}

// Unmarshal takes a pointer to a struct and unmarshals the environment variables to the struct.
func Unmarshal(v interface{}) error {
	return ee.Unmarshal(v)
}

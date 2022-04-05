package eco

var ee *eco

func init() {
	ee = New()
}

func SetPrefix(prefix string) *eco {
	return ee.SetPrefix(prefix)
}

func SetArraySeparator(sep string) *eco {
	return ee.SetArraySeparator(sep)
}

func SetEnvNameTransformer(transformerFunc envNameTransformerFunc) *eco {
	return ee.SetEnvNameTransformer(transformerFunc)
}

func SetEnvNameSeparator(envNameSeparator string) *eco {
	return ee.SetEnvNameSeparator(envNameSeparator)
}

func SetValueGetter(valueGetter envValueGetterFunc) *eco {
	return ee.SetValueGetter(valueGetter)
}

func Unmarshal(model interface{}) error {
	return ee.Unmarshal(model)
}

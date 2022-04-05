package eco

type envNameTransformerFunc func(parts []string, sep string) string
type envValueGetterFunc func(key string) string

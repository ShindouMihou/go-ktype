package runtime

import "reflect"

type TypeTranslator = func(t reflect.Type) string

type biasedTypeTranslatorConfig struct {
	Mapper map[string]string
}

var BiasedTypeTranslatorConfiguration = &biasedTypeTranslatorConfig{
	Mapper: map[string]string{
		"time.Time": "instant",
	},
}

// Map adds a mapping property that maps a Golang type to a more generic extension, for example,
// time.Time can be translated to "instant" which has a import map that maps "instant" to "kotlinx.datetime.Instant".
func (translator *biasedTypeTranslatorConfig) Map(name string, value string) *biasedTypeTranslatorConfig {
	translator.Mapper[name] = value
	return translator
}

func BiasedTypeTranslator(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Uint, reflect.Int, reflect.Uint8, reflect.Uint16, reflect.Int8, reflect.Int16, reflect.Uint32, reflect.Int32:
		return "int"
	case reflect.Int64, reflect.Uint64:
		return "long"
	case reflect.Bool:
		return "boolean"
	case reflect.String:
		return "string"
	case reflect.Float32:
		return "float"
	case reflect.Float64:
		return "double"
	case reflect.Pointer:
		return "nullable[" + BiasedTypeTranslator(t.Elem()) + "]"
	case reflect.Map:
		key := t.Key()
		elem := t.Elem()
		return "map[" + BiasedTypeTranslator(key) + "]" + BiasedTypeTranslator(elem)
	case reflect.Slice:
		elem := t.Elem()
		return "slice[" + BiasedTypeTranslator(elem) + "]"
	case reflect.Struct:
		fname := t.PkgPath() + "." + t.Name()
		if alias, ok := BiasedTypeTranslatorConfiguration.Mapper[fname]; ok {
			return alias
		}
		return t.Name()
	default:
		return "<unknown> (" + t.Name() + ")"
	}
}

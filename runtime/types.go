package runtime

import (
	"errors"
	"reflect"
)

type Class struct {
	Package    string              `json:"package"`
	Name       string              `json:"name"`
	Fields     map[string]string   `json:"fields"`
	Properties map[string]Property `json:"properties"`
}

type Property map[string]string

type Runtime struct {
	Classes        map[string]Class `json:"classes"`
	TypeTranslator TypeTranslator   `json:"-"`
}

func NewRuntime() *Runtime {
	return &Runtime{Classes: make(map[string]Class), TypeTranslator: BiasedTypeTranslator.Translator}
}

func (runtime *Runtime) Load(types ...any) error {
	for _, t := range types {
		class, err := runtime.classFrom(t)
		if err != nil {
			return err
		}
		runtime.Classes[class.Name] = *class
	}
	return nil
}

func (runtime *Runtime) classFrom(t any) (*Class, error) {
	if wrapper, ok := t.(Wrapper); ok {
		class, err := runtime.classFrom(wrapper.Class)
		if err != nil {
			return nil, err
		}
		class.Name = wrapper.Name
		return class, nil
	}

	rtype := reflect.TypeOf(t)

	types := make(map[string]string)
	properties := make(map[string]Property)

	for i := 0; i < rtype.NumField(); i++ {
		field := rtype.Field(i)
		types[field.Name] = runtime.TypeTranslator(field.Type)
		properties[field.Name] = make(Property)

		if json, ok := field.Tag.Lookup("json"); ok {
			properties[field.Name]["json"] = json
		}

		alias := ""

		if typeAlias, ok := field.Tag.Lookup("type_gen"); ok {
			properties[field.Name]["type_gen"] = typeAlias
			types[field.Name] = typeAlias
			alias = typeAlias
		}

		var createDependency = func(alias string, dType reflect.Type) error {
			dependencyName := dType.Name()
			if alias != "" {
				dependencyName = alias
			}
			if _, ok := runtime.Classes[dependencyName]; ok {
				return nil
			}
			if !isDestructable(dType) {
				return nil
			}
			for k := 0; k < dType.NumField(); k++ {
				kField := dType.Field(k)
				if kField.Type.AssignableTo(rtype) {
					return errors.New(rtype.Name() + " implements a field that references " + dType.Name() + " which references the former.")
				}
			}
			dependency, err := runtime.classFrom(reflect.New(dType).Elem().Interface())
			if err != nil {
				return err
			}
			if alias != "" {
				dependency.Name = alias
			}
			runtime.Classes[dependencyName] = *dependency
			return nil
		}

		if field.Type.Kind() == reflect.Struct {
			if err := createDependency(alias, field.Type); err != nil {
				return nil, err
			}
		}

		if field.Type.Kind() == reflect.Map {
			if field.Type.Elem().Kind() == reflect.Struct {
				if err := createDependency(alias, field.Type.Elem()); err != nil {
					return nil, err
				}
			}

			if field.Type.Key().Kind() == reflect.Struct {
				if err := createDependency(alias, field.Type.Elem()); err != nil {
					return nil, err
				}
			}
		}

		if field.Type.Kind() == reflect.Slice {
			if field.Type.Elem().Kind() == reflect.Struct {
				if err := createDependency(alias, field.Type.Elem()); err != nil {
					return nil, err
				}
			}
		}
	}
	return &Class{
		Package:    rtype.PkgPath(),
		Name:       rtype.Name(),
		Fields:     types,
		Properties: properties,
	}, nil
}

func isDestructable(t reflect.Type) bool {
	fname := t.PkgPath() + "." + t.Name()
	switch fname {
	case "time.Time":
		return false
	case "sync.Mutex":
		return false
	default:
		return true
	}
}

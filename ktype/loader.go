package ktype

import (
	"bytes"
	"github.com/bytedance/sonic"
	"go-ktype/runtime"
	"go-ktype/runtime/gen"
)

func LoadClassesFromJson(body []byte) (*runtime.Runtime, error) {
	var runtime runtime.Runtime
	if err := sonic.Unmarshal(body, &runtime); err != nil {
		return nil, err
	}
	return &runtime, nil
}

func LoadTypistFrom(name string, body []byte) gen.Typist {
	typist := gen.NewTypist(bytes.NewReader(body))
	gen.Typists[name] = typist
	return typist
}

func LoadImportMapFrom(body []byte) (gen.ImportMap, error) {
	var importMap gen.ImportMap
	if err := sonic.Unmarshal(body, &importMap); err != nil {
		return nil, err
	}
	return importMap, nil
}

var NoImportMap = make(map[string]string)

func NewTranspiler(typist gen.Typist, importMap gen.ImportMap, properties gen.Properties) *gen.Transpiler {
	return &gen.Transpiler{
		Typist:     typist,
		ImportMap:  importMap,
		Properties: properties,
	}
}

package main

import (
	"fmt"
	"go-ktype/ktype"
	runtime2 "go-ktype/runtime"
	"io"
	"log"
	"os"
)

func main() {
	typistBytes, err := load("generators/kotlin/data_class.kt.typist")
	if err != nil {
		log.Fatalln("failed to load typist file", err)
	}
	importMapBytes, err := load("generators/kotlin/import_maps.json")
	if err != nil {
		log.Fatalln("failed to load import maps", err)
	}
	typist := ktype.LoadTypistFrom("kotlin", typistBytes)
	importMap, err := ktype.LoadImportMapFrom(importMapBytes)
	if err != nil {
		log.Fatalln("failed to load import maps", err)
	}
	properties := ktype.NewKotlinProperties().SetPackage("pro.qucy.ktype.models").Build()
	transpiler := ktype.NewTranspiler(typist, importMap, properties)

	runtime := ktype.NewRuntimeBuilder().WithTypeTranslator(runtime2.BiasedTypeTranslator).Build()
	if err = runtime.Load(SampleClass{}, SampleClassTwo{}); err != nil {
		log.Fatalln("failed to load classes", err)
	}

	fmt.Println(transpiler.Transpile(runtime))
}

type SampleClass struct {
	Hello string `json:"hello"`
}

type SampleClassTwo struct {
	World string `json:"world"`
}

func load(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bytes, err := io.ReadAll(f)
	return bytes, err
}

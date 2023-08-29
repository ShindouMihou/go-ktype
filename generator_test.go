package main

import (
	"go-ktype/internal/testclasses"
	"go-ktype/runtime"
	"go-ktype/runtime/gen"
	"os"
	"testing"
)

func TestTranspile(t *testing.T) {
	err := gen.Load("generators/kotlin/data_class.kt.typist")
	if err != nil {
		t.Fatal(err)
	}
	typist := gen.Typists["generators/kotlin/data_class.kt.typist"]
	transpiler, err := gen.NewTranspiler(typist, "generators/kotlin/import_maps.json", map[string]string{
		"package": "pro.qucy.ktype.models",
	})
	if err != nil {
		t.Fatal(err)
	}
	rtime := runtime.NewRuntime()
	if err = rtime.Load(
		testclasses.SampleWithDependency{},
		testclasses.SampleAliasedDependency{},
		testclasses.WrappedSample,
	); err != nil {
		t.Fatal(err)
	}
	_ = os.MkdirAll(".samples/", 0750)
	f, _ := os.Create(".samples/Models.kt")
	_, err = f.WriteString(transpiler.Transpile(rtime))
	if err != nil {
		t.Fatal(err)
	}
}

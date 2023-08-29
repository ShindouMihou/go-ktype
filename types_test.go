package main

import (
	"encoding/json"
	"fmt"
	"go-ktype/internal/testclasses"
	"go-ktype/runtime"
	"testing"
)

func TestClassFrom(t *testing.T) {
	rtime := runtime.NewRuntime()
	err := rtime.Load(testclasses.SampleWithDependency{})
	if err != nil {
		t.Fatal(err)
	}
	bytes, _ := json.MarshalIndent(rtime.Classes["SampleWithDependency"], "", "  ")
	fmt.Println(string(bytes))
}

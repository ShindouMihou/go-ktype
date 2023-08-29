package go_ktype

import (
	"encoding/json"
	"fmt"
	"go-ktype/internal/testclasses"
	"go-ktype/runtime"
	"testing"
)

func TestClassFrom(t *testing.T) {
	runtime := runtime.NewRuntime()
	err := runtime.Load(testclasses.SampleWithDependency{})
	if err != nil {
		t.Fatal(err)
	}
	bytes, _ := json.MarshalIndent(runtime.Classes["SampleWithDependency"], "", "  ")
	fmt.Println(string(bytes))
}

package ktype

import "go-ktype/runtime"

type RuntimeBuilder struct {
	runtime *runtime.Runtime
}

func NewRuntimeBuilder() *RuntimeBuilder {
	return &RuntimeBuilder{runtime: runtime.NewRuntime()}
}

func (builder *RuntimeBuilder) WithTypeTranslator(translator runtime.TypeTranslator) *RuntimeBuilder {
	builder.runtime.TypeTranslator = translator
	return builder
}

func (builder *RuntimeBuilder) Build() *runtime.Runtime {
	return builder.runtime
}

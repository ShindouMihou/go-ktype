package ktype

import "go-ktype/runtime/gen"

func NewKotlinProperties() *KotlinPropertyBuilder {
	return &KotlinPropertyBuilder{properties: make(gen.Properties)}
}

type KotlinPropertyBuilder struct {
	properties gen.Properties
}

func (builder *KotlinPropertyBuilder) SetPackage(pkg string) *KotlinPropertyBuilder {
	builder.properties["package"] = pkg
	return builder
}

func (builder *KotlinPropertyBuilder) Build() gen.Properties {
	return builder.properties
}

package testclasses

import "go-ktype/runtime"

type SampleWithDependency struct {
	Dependency Dependency `json:"dependency"`
}

type Dependency struct {
	Id string `json:"id"`
}

type SampleAliasedDependency struct {
	AliasedDependency AliasedDependency `json:"aliased_dependency" type_gen:"AnAliasedDependency"`
}

type AliasedDependency struct {
	Server uint64 `json:"server"`
}

type SampleWrappedSample struct {
	Title string `json:"title"`
}

var WrappedSample = runtime.Wrapper{Name: "SomeWrappedSample", Class: SampleWrappedSample{}}

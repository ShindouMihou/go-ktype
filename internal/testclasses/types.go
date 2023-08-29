package testclasses

import (
	"go-ktype/runtime"
	"time"
)

type SampleWithDependency struct {
	Dependency Dependency `json:"dependency"`
}

type Dependency struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SampleAliasedDependency struct {
	AliasedDependency AliasedDependency `json:"aliased_dependency" type_gen:"AnAliasedDependency"`
}

type AliasedDependency struct {
	Server    uint64    `json:"server"`
	Token     string    `json:"token"`
	Timestamp time.Time `json:"timestamp"`
}

type SampleWrappedSample struct {
	Title string `json:"title"`
}

var WrappedSample = runtime.Wrapper{Name: "SomeWrappedSample", Class: SampleWrappedSample{}}

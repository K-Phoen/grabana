package sandbox

import "github.com/K-Phoen/grabana/gen/dashboard/types"

type Option func(builder *Builder) error

type Builder struct {
	internal *Sandbox
}

func Name(name string) Option {
	return func(builder *Builder) error {
		
		builder.internal.Name = name

		return nil
	}
}

func NestedStruct(nestedStruct struct {
	Foo string `json:"foo"`
}) Option {
	return func(builder *Builder) error {
		
		builder.internal.NestedStruct = nestedStruct

		return nil
	}
}

func AnythingPlz(anythingPlz any) Option {
	return func(builder *Builder) error {
		
		builder.internal.AnythingPlz = anythingPlz

		return nil
	}
}

func SomeMap(someMap map[string]int32) Option {
	return func(builder *Builder) error {
		
		builder.internal.SomeMap = &someMap

		return nil
	}
}

func OperatorStates(operatorStates map[string]OperatorState) Option {
	return func(builder *Builder) error {
		
		builder.internal.OperatorStates = &operatorStates

		return nil
	}
}


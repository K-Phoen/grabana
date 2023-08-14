package types

import "cuelang.org/go/cue"

type kindMapping struct {
	terraformType string
	golangType    string
	terraformFunc string
}

var kindMappings = map[cue.Kind]*kindMapping{
	cue.BoolKind: {
		terraformType: "Bool",
		golangType:    "bool",
		terraformFunc: "ValueBool",
	},
	cue.IntKind: {
		terraformType: "Int64",
		golangType:    "int64",
		terraformFunc: "ValueInt64",
	},
	cue.FloatKind: {
		terraformType: "Float64",
		golangType:    "float64",
		terraformFunc: "ValueFloat64",
	},
	cue.NumberKind: {
		terraformType: "Float64",
		golangType:    "float64",
		terraformFunc: "ValueFloat64",
	},
	cue.StringKind: {
		terraformType: "String",
		golangType:    "string",
		terraformFunc: "ValueString",
	},
}

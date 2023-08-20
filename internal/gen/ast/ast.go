package ast

type Kind string

const (
	KindDisjunction Kind = "disjunction"

	KindStruct Kind = "struct"
	KindEnum   Kind = "enum"
	KindMap    Kind = "map"

	KindNull   Kind = "null"
	KindAny    Kind = "any"
	KindBytes  Kind = "bytes"
	KindArray  Kind = "array"
	KindString Kind = "string"

	KindFloat32 Kind = "float32"
	KindFloat64 Kind = "float64"

	KindUint8  Kind = "uint8"
	KindUint16 Kind = "uint16"
	KindUint32 Kind = "uint32"
	KindUint64 Kind = "uint64"
	KindInt8   Kind = "int8"
	KintInt16  Kind = "int16"
	KindInt32  Kind = "int32"
	KindInt64  Kind = "int64"

	KindBool Kind = "bool"
)

type Definition struct {
	Kind         Kind
	Name         string
	Comments     []string
	IndexType    Kind              // for maps & arrays
	ValueType    *Definition       // for maps & arrays
	Branches     Definitions       // for disjunctions
	Fields       []FieldDefinition // for structs
	Values       []EnumValue       // for enums
	IsEntryPoint bool              // Dashboard is an entryPoint type. DashboardStyle isn't.

	Nullable    bool
	Constraints []TypeConstraint

	Default any // the go type of the value depends on Kind
}

func (def Definition) IsReference() bool {
	switch def.Kind {
	case KindDisjunction:
		return false
	case KindStruct:
		return false
	case KindEnum:
		return false
	case KindMap:
		return false
	case KindNull:
		return false
	case KindAny:
		return false
	case KindBytes:
		return false
	case KindArray:
		return false
	case KindString:
		return false
	case KindFloat32:
		return false
	case KindFloat64:
		return false
	case KindUint8:
		return false
	case KindUint16:
		return false
	case KindUint32:
		return false
	case KindUint64:
		return false
	case KindInt8:
		return false
	case KintInt16:
		return false
	case KindInt32:
		return false
	case KindInt64:
		return false
	case KindBool:
		return false
	}

	return true
}

type Definitions []Definition

func (defs Definitions) HasNullType() bool {
	for _, t := range defs {
		if t.Kind == KindNull {
			return true
		}
	}

	return false
}

func (defs Definitions) NonNullTypes() Definitions {
	var filteredTypes Definitions
	for _, def := range defs {
		if def.Kind == KindNull {
			continue
		}

		filteredTypes = append(filteredTypes, def)
	}

	return filteredTypes
}

type TypeConstraint struct {
	// TODO: something more descriptive here? constant?
	Op   string
	Args []any
}

type EnumValue struct {
	Type  Kind
	Name  string
	Value interface{}
}

type FieldDefinition struct {
	Name     string
	Comments []string
	Required bool
	Type     Definition
}

type File struct {
	Package string
	Types   []Definition
}

func (file *File) EntryPointType() (Definition, bool) {
	for _, typeDef := range file.Types {
		if typeDef.IsEntryPoint {
			return typeDef, true
		}
	}

	return Definition{}, false
}

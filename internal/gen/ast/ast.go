package ast

type TypeID string

const (
	TypeDisjunction TypeID = "disjunction"

	TypeStruct TypeID = "struct"
	TypeEnum   TypeID = "enum"
	TypeMap    TypeID = "map"

	TypeNull   TypeID = "null"
	TypeAny    TypeID = "any"
	TypeBytes  TypeID = "bytes"
	TypeArray  TypeID = "array"
	TypeString TypeID = "string"

	TypeFloat32 TypeID = "float32"
	TypeFloat64 TypeID = "float64"

	TypeUint8  TypeID = "uint8"
	TypeUint16 TypeID = "uint16"
	TypeUint32 TypeID = "uint32"
	TypeUint64 TypeID = "uint64"
	TypeInt8   TypeID = "int8"
	TypeInt16  TypeID = "int16"
	TypeInt32  TypeID = "int32"
	TypeInt64  TypeID = "int64"

	TypeBool TypeID = "bool"
)

type Definition struct {
	Type         TypeID
	Name         string
	Comments     []string
	IndexType    TypeID            // for maps & arrays
	ValueType    *Definition       // for maps & arrays
	Branches     Definitions       // for disjunctions
	Fields       []FieldDefinition // for structs
	Values       []EnumValue       // for enums
	IsEntryPoint bool              // Dashboard is an entryPoint type. DashboardStyle isn't.

	Nullable    bool
	Constraints []TypeConstraint
}

func (def Definition) IsReference() bool {
	switch def.Type {
	case TypeDisjunction:
		return false
	case TypeStruct:
		return false
	case TypeEnum:
		return false
	case TypeMap:
		return false
	case TypeNull:
		return false
	case TypeAny:
		return false
	case TypeBytes:
		return false
	case TypeArray:
		return false
	case TypeString:
		return false
	case TypeFloat32:
		return false
	case TypeFloat64:
		return false
	case TypeUint8:
		return false
	case TypeUint16:
		return false
	case TypeUint32:
		return false
	case TypeUint64:
		return false
	case TypeInt8:
		return false
	case TypeInt16:
		return false
	case TypeInt32:
		return false
	case TypeInt64:
		return false
	case TypeBool:
		return false
	}

	return true
}

type Definitions []Definition

func (types Definitions) HasNullType() bool {
	for _, t := range types {
		if t.Type == TypeNull {
			return true
		}
	}

	return false
}

func (types Definitions) NonNullTypes() Definitions {
	var filteredTypes Definitions
	for _, t := range types {
		if t.Type == TypeNull {
			continue
		}

		filteredTypes = append(filteredTypes, t)
	}

	return filteredTypes
}

type TypeConstraint struct {
	// TODO
	Op   string
	Args []any
}

type EnumValue struct {
	Type  TypeID
	Name  string
	Value interface{}
}

type FieldDefinition struct {
	Name     string
	Comments []string
	// Field needs to be defined
	Required bool
	Type     Definition
	// TODO
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

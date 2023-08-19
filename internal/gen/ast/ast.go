package ast

type TypeID string

const (
	TypeDisjunction TypeID = "disjunction"
	TypeNull        TypeID = "null"
	TypeBytes       TypeID = "bytes"
	TypeArray       TypeID = "array"
	TypeString      TypeID = "string"

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
	TypeAny  TypeID = "any"
)

type DefinitionType string

const (
	DefinitionStruct DefinitionType = "struct"
	DefinitionEnum   DefinitionType = "enum"
)

type TypeDefinition struct {
	Type         DefinitionType
	Name         string
	SubType      TypeID
	Comments     []string
	Fields       []FieldDefinition // for structs
	Values       []EnumValue       // for enums
	IsEntryPoint bool              // Dashboard is an entryPoint type. DashboardStyle isn't.
}

type TypeConstraint struct {
	// TODO
	Op   string
	Args []any
}

type FieldTypes []FieldType

func (types FieldTypes) HasNullType() bool {
	for _, t := range types {
		if t.IsNull() {
			return true
		}
	}

	return false
}

func (types FieldTypes) NonNullTypes() FieldTypes {
	var filteredTypes FieldTypes
	for _, t := range types {
		if t.IsNull() {
			continue
		}

		filteredTypes = append(filteredTypes, t)
	}

	return filteredTypes
}

type EnumValue struct {
	Name  string
	Value interface{}
}

type FieldType struct {
	Nullable    bool
	Type        TypeID
	SubType     FieldTypes
	Constraints []TypeConstraint
}

func (t FieldType) IsNull() bool {
	return t.Type == TypeNull
}

type FieldDefinition struct {
	Name     string
	Comments []string
	// Field needs to be defined
	Required bool
	Type     FieldType
	// TODO
}

type File struct {
	Package string
	Types   []TypeDefinition
}

func (file *File) EntryPointType() (TypeDefinition, bool) {
	for _, typeDef := range file.Types {
		if typeDef.IsEntryPoint {
			return typeDef, true
		}
	}

	return TypeDefinition{}, false
}
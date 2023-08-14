package simplecue

type TypeID string

const (
	TypeDisjunction TypeID = "disjunction"
	TypeNull        TypeID = "null"
	TypeBytes       TypeID = "bytes"
	TypeArray       TypeID = "array"
	TypeString      TypeID = "string"
	TypeFloat64     TypeID = "float64"
	TypeInt64       TypeID = "int64"
	TypeBool        TypeID = "bool"
	TypeAny         TypeID = "any"
	TypeUnknown     TypeID = "unknown"
)

type DefinitionType string

const (
	DefinitionStruct DefinitionType = "struct"
	DefinitionEnum   DefinitionType = "enum"
	DefinitionAlias  DefinitionType = "alias"
)

type TypeDefinition struct {
	Type     DefinitionType
	Name     string
	SubType  string
	Comments []string
	Fields   []FieldDefinition // for structs
	Values   []EnumValue       // for enums
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

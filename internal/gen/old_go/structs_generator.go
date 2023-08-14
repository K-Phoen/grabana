package gocode

import (
	"bytes"
	"fmt"
	"strings"

	"cuelang.org/go/cue"
)

type FieldNamer func(fieldName string) string

type Config struct {
	FieldNamer FieldNamer
}

type fieldDef struct {
	Name     string
	Type     string
	Comments []string
	Optional bool
	Tags     []string
}

type structDef struct {
	Name     string
	Comments []string
	Fields   []fieldDef
}

// GenerateStructs takes a cue.Value and generates the corresponding Go structs.
func GenerateStructs(schema cue.Value) (b []byte, err error) {
	cfg := Config{
		FieldNamer: PublicField(),
	}

	var structs []structDef
	fieldsIt, err := schema.Fields(cue.Optional(true), cue.Definitions(true))
	if err != nil {
		return nil, err
	}

	for fieldsIt.Next() {
		/*
			if fieldsIt.Value().Kind() == cue.BottomKind {
				fmt.Printf("bottom: %s\n", fieldsIt.Selector().String())
				continue
			}
		*/

		fmt.Printf("lala: %s\n", fieldsIt.Selector().String())
		def, err := generateStruct(fieldsIt.Value(), cfg)
		if err != nil {
			continue // TODO
			//return nil, err
		}

		structs = append(structs, def)
	}

	buf := new(bytes.Buffer)
	tmplVars := struct {
		Structs []structDef
	}{
		Structs: structs,
	}
	if err := tmpls.Lookup("structs.tmpl").Execute(buf, tmplVars); err != nil {
		return nil, fmt.Errorf("failed executing structs template: %w", err)
	}

	return buf.Bytes(), nil
	//return format.Source(buf.Bytes())
}

func generateStruct(value cue.Value, cfg Config) (structDef, error) {
	label, ok := value.Label()
	if !ok {
		return structDef{}, fmt.Errorf("could not extract label from struct")
	}

	if value.IncompleteKind() != cue.StructKind {
		return structDef{}, fmt.Errorf("can not generate struct from '%s' input (%s)", value.IncompleteKind(), label)
	}

	def := structDef{
		Name:     cleanupStructName(label),
		Comments: formatDoc(value.Doc()),
	}

	fmt.Printf("Exploring: '%s'\n", label)

	fieldsIt, err := value.Fields(cue.Optional(true), cue.Definitions(true))
	if err != nil {
		return def, fmt.Errorf("could not build struct fields iterator: %w", err)
	}

	for fieldsIt.Next() {
		field, err := generateStructField(fieldsIt.Value(), fieldsIt.IsOptional(), cfg)
		if err != nil {
			return def, err
		}

		def.Fields = append(def.Fields, field)
	}

	return def, nil
}

func generateStructField(value cue.Value, optional bool, cfg Config) (fieldDef, error) {
	label, ok := value.Label()
	if !ok {
		return fieldDef{}, fmt.Errorf("could not extract label from field")
	}

	fmt.Printf("\tfield: '%s', kind: '%s'\n", label, value.IncompleteKind())

	fType, err := cueTypeToGo(value)
	if err != nil {
		return fieldDef{}, fmt.Errorf("can not infer field type: '%w'", err)
	}

	return fieldDef{
		Name:     cfg.FieldNamer(label),
		Comments: formatDoc(value.Doc()),
		Optional: optional,
		Type:     fType,
		Tags: []string{
			fmt.Sprintf("json:\"%s\"", label),
		},
	}, nil
}

func cueTypeToGo(value cue.Value) (string, error) {
	switch value.IncompleteKind() {
	case cue.BoolKind:
		return "bool", nil
	case cue.IntKind:
		return "int64", nil
	case cue.FloatKind:
		return "float64", nil
	case cue.NumberKind:
		return "float64", nil
	case cue.StringKind:
		return "string", nil
	case cue.ListKind:
		return fmt.Sprintf("[]%s", "lala"), nil
	default:
		fmt.Printf("unknown type '%s'\n", value.IncompleteKind())
		return "", fmt.Errorf("unknown type %s", value.IncompleteKind())
	}
}

func PublicField() FieldNamer {
	return func(fieldName string) string {
		firstChar := string(strings.ToUpper(fieldName)[0])

		return firstChar + fieldName[1:]
	}
}

func cleanupStructName(name string) string {
	return strings.TrimPrefix(name, "#")
}

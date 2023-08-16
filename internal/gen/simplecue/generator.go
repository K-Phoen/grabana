package simplecue

import (
	"strings"

	"cuelang.org/go/cue"
)

const annotationName = "grabana"
const hintKindEnum = "enum"
const annotationKindFieldName = "kind"
const enumMembersAttr = "memberNames"

// Config governs certain variable behaviors when converting CUE to Typescript.
type Config struct {
	// Package name used to generate code into.
	Package string
}

type newGenerator struct {
}

func GenerateAST(val cue.Value, c Config) (*File, error) {
	g := &newGenerator{}

	file := File{
		Package: c.Package,
	}

	i, err := val.Fields(cue.Definitions(true))
	if err != nil {
		return nil, err
	}
	for i.Next() {
		sel := i.Selector()

		n, err := g.declareTopLevelType(selectorLabel(sel), i.Value(), sel.IsDefinition())
		if err != nil {
			return nil, err
		}

		file.Types = append(file.Types, *n)
	}

	return &file, nil
}

func (g *newGenerator) declareTopLevelType(name string, v cue.Value, isCueDefinition bool) (*TypeDefinition, error) {
	typeHint, err := getTypeHint(v)
	if err != nil {
		return nil, err
	}

	if typeHint == hintKindEnum {
		return g.declareTopLevelEnum(name, v)
	}

	switch v.IncompleteKind() {
	case cue.StructKind:
		return g.declareTopLevelStruct(name, v, isCueDefinition)
	default:
		return nil, errorWithCueRef(v, "unexpected top-level kind '%s'", v.IncompleteKind().String())
	}
}

func (g *newGenerator) declareTopLevelEnum(name string, v cue.Value) (*TypeDefinition, error) {
	// Restrict the expression of enums to ints or strings.
	allowed := cue.StringKind | cue.IntKind
	ik := v.IncompleteKind()
	if ik&allowed != ik {
		return nil, errorWithCueRef(v, "enums may only be generated from concrete strings, or ints")
	}

	values, err := g.extractEnumValues(v)
	if err != nil {
		return nil, err
	}

	subType := TypeString
	if ik == cue.IntKind {
		subType = TypeInt64
	}

	return &TypeDefinition{
		Type:     DefinitionEnum,
		SubType:  string(subType),
		Name:     name,
		Comments: commentsFromCueValue(v),
		Values:   values,
	}, nil
}

func (g *newGenerator) extractEnumValues(v cue.Value) ([]EnumValue, error) {
	_, dvals := v.Expr()
	a := v.Attribute(annotationName)

	var attrMemberNameExist bool
	var evals []string
	if a.Err() == nil {
		val, found, err := a.Lookup(0, enumMembersAttr)
		if err == nil && found {
			attrMemberNameExist = true
			evals = strings.Split(val, "|")
			if len(evals) != len(dvals) {
				return nil, errorWithCueRef(v, "enums and %s attributes size doesn't match", enumMembersAttr)
			}
		}
	}

	// We only allowed String Enum to be generated without memberName attribute
	if v.IncompleteKind() != cue.StringKind && !attrMemberNameExist {
		return nil, errorWithCueRef(v, "numeric enums may only be generated from memberNames attribute")
	}

	var fields []EnumValue
	for idx, dv := range dvals {
		var text string
		if attrMemberNameExist {
			text = evals[idx]
		} else {
			text, _ = dv.String()
		}

		if !dv.IsConcrete() {
			return nil, errorWithCueRef(v, "enums may only be generated from a disjunction of concrete strings or numbers")
		}

		val, err := cueConcreteToScalar(dv)
		if err != nil {
			return nil, err
		}
		fields = append(fields, EnumValue{
			// Simple mapping of all enum values (which we are assuming are in
			// lowerCamelCase) to corresponding CamelCase
			Name:  text,
			Value: val,
		})
	}

	return fields, nil
}

func (g *newGenerator) declareTopLevelStruct(name string, v cue.Value, isCueDefinition bool) (*TypeDefinition, error) {
	// This check might be too restrictive
	if v.IncompleteKind() != cue.StructKind {
		return nil, errorWithCueRef(v, "top-level type definitions may only be generated from structs")
	}

	typeDef := &TypeDefinition{
		Type:         DefinitionStruct,
		Name:         name,
		Comments:     commentsFromCueValue(v),
		IsEntryPoint: !isCueDefinition,
	}

	// explore struct fields
	for i, _ := v.Fields(cue.Optional(true), cue.Definitions(true)); i.Next(); {
		fieldLabel := selectorLabel(i.Selector())

		node, err := g.declareNode(i.Value())
		if err != nil {
			return nil, err
		}

		typeDef.Fields = append(typeDef.Fields, FieldDefinition{
			Name:     fieldLabel,
			Comments: commentsFromCueValue(i.Value()),
			Required: !i.IsOptional(),
			Type:     *node,
		})
	}

	return typeDef, nil
}

func (g *newGenerator) declareNode(v cue.Value) (*FieldType, error) {
	// This node is referring to another definition
	_, path := v.ReferencePath()
	if path.String() != "" {
		return &FieldType{
			Nullable: false, // TODO
			Type:     TypeID(path.String()[1:]),
		}, nil
	}

	disjunctions := appendSplit(nil, cue.OrOp, v)
	if len(disjunctions) != 1 {
		subTypes := make([]FieldType, 0, len(disjunctions))
		for _, subTypeValue := range disjunctions {
			subType, err := g.declareNode(subTypeValue)
			if err != nil {
				return nil, err
			}

			subTypes = append(subTypes, *subType)
		}

		return &FieldType{
			Type:     TypeDisjunction,
			SubType:  subTypes,
			Nullable: false, // TODO
		}, nil
	}

	// remove defaults... for now.
	defv, _ := v.Default()
	if !defv.Equals(v) {
		_, dvals := v.Expr()
		v = dvals[0]
	}

	switch v.IncompleteKind() {
	case cue.TopKind:
		return &FieldType{Type: TypeAny}, nil
	case cue.NullKind:
		return &FieldType{Type: TypeNull}, nil
	case cue.BoolKind:
		return &FieldType{Type: TypeBool}, nil
	case cue.BytesKind:
		return &FieldType{Type: TypeBytes}, nil
	case cue.StringKind:
		return &FieldType{Type: TypeString}, nil
	case cue.FloatKind, cue.NumberKind:
		return g.declareNumber(v, TypeFloat64)
	case cue.IntKind:
		return g.declareNumber(v, TypeInt64)
	case cue.ListKind:
		return g.declareList(v)
	case cue.StructKind:
		op, _ := v.Expr()

		// in cue: {...}
		if op == cue.NoOp {
			return &FieldType{Type: TypeAny}, nil
		}

		return nil, errorWithCueRef(v, "nested struct definitions are not supported")
	default:
		return nil, errorWithCueRef(v, "unexpected node with kind '%s'", v.IncompleteKind().String())
	}
}

func (g *newGenerator) declareNumber(v cue.Value, typeHint TypeID) (*FieldType, error) {
	typeDef := &FieldType{
		Type:     typeHint,
		Nullable: false,
	}

	// TODO: better representation of the type
	// we currently convert everything to int64, float64 or number

	constraints, err := g.declareNumberConstraints(v)
	if err != nil {
		return nil, err
	}

	typeDef.Constraints = constraints

	return typeDef, nil
}

func (g *newGenerator) declareNumberConstraints(v cue.Value) ([]TypeConstraint, error) {
	// typeAndConstraints can contain the following cue expressions:
	// 	- number
	// 	- int|float, number, upper bound, lower bound
	typeAndConstraints := appendSplit(nil, cue.AndOp, v)

	// nothing to do
	if len(typeAndConstraints) == 1 {
		return nil, nil
	}

	constraints := make([]TypeConstraint, 0, len(typeAndConstraints))

	constraintsStartIndex := 1

	// don't include type-related constraints
	if len(typeAndConstraints) > 1 && typeAndConstraints[0].IncompleteKind() != cue.NumberKind {
		constraintsStartIndex = 3
	}

	for _, s := range typeAndConstraints[constraintsStartIndex:] {
		constraint, err := g.extractConstraint(s)
		if err != nil {
			return nil, err
		}

		constraints = append(constraints, constraint)
	}

	return constraints, nil
}

func (g *newGenerator) extractConstraint(v cue.Value) (TypeConstraint, error) {
	toConstraint := func(operator string, arg cue.Value) (TypeConstraint, error) {
		scalar, err := cueConcreteToScalar(arg)
		if err != nil {
			return TypeConstraint{}, err
		}

		return TypeConstraint{
			Op:   operator,
			Args: []any{scalar},
		}, nil
	}

	switch op, a := v.Expr(); op {
	case cue.LessThanOp:
		return toConstraint("<", a[0])
	case cue.LessThanEqualOp:
		return toConstraint("<=", a[0])
	case cue.GreaterThanOp:
		return toConstraint(">", a[0])
	case cue.GreaterThanEqualOp:
		return toConstraint(">=", a[0])
	case cue.NotEqualOp:
		return toConstraint("!=", a[0])
	default:
		return TypeConstraint{}, errorWithCueRef(v, "unsupported op for number %v", op)
	}
}

func (g *newGenerator) declareList(v cue.Value) (*FieldType, error) {
	i, err := v.List()
	if err != nil {
		return nil, err
	}

	typeDef := &FieldType{
		Type:        TypeArray,
		Nullable:    false,
		SubType:     nil,
		Constraints: nil,
	}

	// works only for a closed/concrete list
	if v.IsConcrete() {
		for i.Next() {
			node, err := g.declareNode(i.Value())
			if err != nil {
				return nil, err
			}

			typeDef.SubType = append(typeDef.SubType, *node)
		}

		return typeDef, nil
	}

	// open list

	// If the default (all lists have a default, usually self, ugh) differs from the
	// input list, peel it off. Otherwise our AnyIndex lookup may end up getting
	// sent on the wrong path.
	defv, _ := v.Default()
	if !defv.Equals(v) {
		_, dvals := v.Expr()
		v = dvals[0]
	}

	e := v.LookupPath(cue.MakePath(cue.AnyIndex))
	if !e.Exists() {
		// unreachable?
		return nil, errorWithCueRef(v, "open list must have a type")
	}

	expr, err := g.declareNode(e)
	if err != nil {
		return nil, err
	}

	typeDef.SubType = []FieldType{*expr}

	return typeDef, nil
}

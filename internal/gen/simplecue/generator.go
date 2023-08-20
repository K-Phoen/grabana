package simplecue

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"
	"github.com/K-Phoen/grabana/internal/gen/ast"
)

const annotationName = "grabana"
const hintKindEnum = "enum"
const annotationKindFieldName = "kind"
const enumMembersAttr = "memberNames"

var errSkip = fmt.Errorf("skip")

type Config struct {
	// Package name used to generate code into.
	Package string
}

type newGenerator struct {
	file                    *ast.File
	currentTopLevelTypeName string
}

func GenerateAST(val cue.Value, c Config) (*ast.File, error) {
	g := &newGenerator{
		file: &ast.File{
			Package: c.Package,
		},
	}

	i, err := val.Fields(cue.Definitions(true))
	if err != nil {
		return nil, err
	}
	for i.Next() {
		sel := i.Selector()

		g.currentTopLevelTypeName = selectorLabel(sel)

		n, err := g.declareTopLevelType(g.currentTopLevelTypeName, i.Value(), sel.IsDefinition())
		if err != nil {
			return nil, err
		}

		g.file.Types = append(g.file.Types, *n)
	}

	return g.file, nil
}

func (g *newGenerator) declareTopLevelType(name string, v cue.Value, isCueDefinition bool) (*ast.Definition, error) {
	typeHint, err := getTypeHint(v)
	if err != nil {
		return nil, err
	}

	// Hinted as an enum
	if typeHint == hintKindEnum {
		return g.declareTopLevelEnum(name, v)
	}

	ik := v.IncompleteKind()

	// Is it a string disjunction that we can turn into an enum?
	disjunctions := appendSplit(nil, cue.OrOp, v)
	if len(disjunctions) != 1 && ik&cue.StringKind == ik {
		return g.declareTopLevelEnum(name, v)
	}

	switch v.IncompleteKind() {
	case cue.StructKind:
		return g.declareTopLevelStruct(name, v, isCueDefinition)
	default:
		return nil, errorWithCueRef(v, "unexpected top-level kind '%s'", v.IncompleteKind().String())
	}
}

func (g *newGenerator) declareTopLevelEnum(name string, v cue.Value) (*ast.Definition, error) {
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

	return &ast.Definition{
		Type:     ast.TypeEnum,
		Name:     name,
		Comments: commentsFromCueValue(v),
		Values:   values,
	}, nil
}

func (g *newGenerator) extractEnumValues(v cue.Value) ([]ast.EnumValue, error) {
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

	subType := ast.TypeString
	if v.IncompleteKind() == cue.IntKind {
		subType = ast.TypeInt64
	}

	var fields []ast.EnumValue
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
		fields = append(fields, ast.EnumValue{
			Type: subType,

			// Simple mapping of all enum values (which we are assuming are in
			// lowerCamelCase) to corresponding CamelCase
			Name:  text,
			Value: val,
		})
	}

	return fields, nil
}

func (g *newGenerator) declareTopLevelStruct(name string, v cue.Value, isCueDefinition bool) (*ast.Definition, error) {
	// This check might be too restrictive
	if v.IncompleteKind() != cue.StructKind {
		return nil, errorWithCueRef(v, "top-level type definitions may only be generated from structs")
	}

	typeDef := &ast.Definition{
		Type:         ast.TypeStruct,
		Name:         name,
		Comments:     commentsFromCueValue(v),
		IsEntryPoint: !isCueDefinition,
	}

	// explore struct fields
	fields, err := g.structFields(v)
	if err != nil {
		return nil, err
	}

	typeDef.Fields = fields

	return typeDef, nil
}

func (g *newGenerator) structFields(v cue.Value) ([]ast.FieldDefinition, error) {
	// This check might be too restrictive
	if v.IncompleteKind() != cue.StructKind {
		return nil, errorWithCueRef(v, "top-level type definitions may only be generated from structs")
	}

	var fields []ast.FieldDefinition

	// explore struct fields
	for i, _ := v.Fields(cue.Optional(true), cue.Definitions(true)); i.Next(); {
		fieldLabel := selectorLabel(i.Selector())

		node, err := g.declareNode(i.Value())
		if err == errSkip {
			continue
		}
		if err != nil {
			return nil, err
		}

		fields = append(fields, ast.FieldDefinition{
			Name:     fieldLabel,
			Comments: commentsFromCueValue(i.Value()),
			Required: !i.IsOptional(),
			Type:     *node,
		})
	}

	return fields, nil
}

func (g *newGenerator) declareNode(v cue.Value) (*ast.Definition, error) {
	// This node is referring to another definition
	_, path := v.ReferencePath()
	if path.String() != "" {
		selectors := path.Selectors()

		return &ast.Definition{
			Type: ast.TypeID(selectorLabel(selectors[len(selectors)-1])),
		}, nil
	}

	disjunctions := appendSplit(nil, cue.OrOp, v)
	if len(disjunctions) != 1 {
		allowedKindsForAnonymousEnum := cue.StringKind | cue.IntKind
		ik := v.IncompleteKind()
		if ik&allowedKindsForAnonymousEnum == ik {
			return g.declareAnonymousEnum(v)
		}

		subTypes := make([]ast.Definition, 0, len(disjunctions))
		for _, subTypeValue := range disjunctions {
			subType, err := g.declareNode(subTypeValue)
			if err != nil {
				return nil, err
			}

			subTypes = append(subTypes, *subType)
		}

		return &ast.Definition{
			Type:     ast.TypeDisjunction,
			Branches: subTypes,
			Nullable: false,
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
		return &ast.Definition{Type: ast.TypeAny}, nil
	case cue.NullKind:
		return &ast.Definition{Type: ast.TypeNull}, nil
	case cue.BoolKind:
		return &ast.Definition{Type: ast.TypeBool}, nil
	case cue.BytesKind:
		return &ast.Definition{Type: ast.TypeBytes}, nil
	case cue.StringKind:
		return &ast.Definition{Type: ast.TypeString}, nil
	case cue.FloatKind, cue.NumberKind, cue.IntKind:
		return g.declareNumber(v)
	case cue.ListKind:
		return g.declareList(v)
	case cue.StructKind:
		pathParts := v.Path().Selectors()
		selector := pathParts[len(pathParts)-1]

		// inline definition of a struct
		if selector.IsDefinition() {
			def, err := g.declareTopLevelStruct(selectorLabel(selector), v, true)
			if err != nil {
				return nil, err
			}

			g.file.Types = append(g.file.Types, *def)

			return nil, errSkip
		}

		op, opArgs := v.Expr()

		// in cue: {...}, {[string]: type}, or inline struct

		if op == cue.NoOp {
			// looking for {[string]: type}
			// (don't know how to do this properly)
			t, ok := opArgs[0].Elem()
			if ok && t.IncompleteKind() != cue.TopKind {
				typeDef, err := g.declareNode(t)
				if err != nil {
					return nil, err
				}

				return &ast.Definition{
					Type:      ast.TypeMap,
					IndexType: ast.TypeString,
					ValueType: typeDef,
				}, nil
			}
		}

		fields, err := g.structFields(v)
		if err != nil {
			return nil, err
		}

		// {...}
		if len(fields) == 0 {
			return &ast.Definition{Type: ast.TypeAny}, nil
		}

		return &ast.Definition{Type: ast.TypeStruct, Fields: fields}, nil
	default:
		return nil, errorWithCueRef(v, "unexpected node with kind '%s'", v.IncompleteKind().String())
	}
}

func (g *newGenerator) declareAnonymousEnum(v cue.Value) (*ast.Definition, error) {
	fieldName, ok := v.Label()
	if !ok {
		return nil, errorWithCueRef(v, "could not determine field name")
	}

	enumName := g.currentTopLevelTypeName + strings.Title(fieldName)
	enumType, err := g.declareTopLevelEnum(enumName, v)
	if err != nil {
		return nil, err
	}

	g.file.Types = append(g.file.Types, *enumType)

	return &ast.Definition{
		Type: ast.TypeID(enumType.Name),
	}, nil
}

func (g *newGenerator) declareNumber(v cue.Value) (*ast.Definition, error) {
	numberTypeWithConstraintsAsString, err := format.Node(v.Syntax())
	if err != nil {
		return nil, err
	}
	parts := strings.Split(string(numberTypeWithConstraintsAsString), " ")
	if len(parts) == 0 {
		return nil, errorWithCueRef(v, "something went very wrong while formatting a number expression into a string")
	}

	// dirty way of preserving the actual type from cue
	// FIXME: fails if the type has a custom bound that further restricts the values
	// IE: uint8 & < 12 will be printed as "uint & < 12
	var numberType ast.TypeID
	switch ast.TypeID(parts[0]) {
	case ast.TypeFloat32, ast.TypeFloat64:
		numberType = ast.TypeID(parts[0])
	case ast.TypeUint8, ast.TypeUint16, ast.TypeUint32, ast.TypeUint64:
		numberType = ast.TypeID(parts[0])
	case ast.TypeInt8, ast.TypeInt16, ast.TypeInt32, ast.TypeInt64:
		numberType = ast.TypeID(parts[0])
	case "uint":
		numberType = ast.TypeUint64
	case "int":
		numberType = ast.TypeInt64
	case "number":
		numberType = ast.TypeFloat64
	default:
		return nil, errorWithCueRef(v, "unknown number type '%s'", parts[0])
	}

	typeDef := &ast.Definition{
		Type:     numberType,
		Nullable: false,
	}

	constraints, err := g.declareNumberConstraints(v)
	if err != nil {
		return nil, err
	}

	typeDef.Constraints = constraints

	return typeDef, nil
}

func (g *newGenerator) declareNumberConstraints(v cue.Value) ([]ast.TypeConstraint, error) {
	// typeAndConstraints can contain the following cue expressions:
	// 	- number
	// 	- int|float, number, upper bound, lower bound
	typeAndConstraints := appendSplit(nil, cue.AndOp, v)

	// nothing to do
	if len(typeAndConstraints) == 1 {
		return nil, nil
	}

	constraints := make([]ast.TypeConstraint, 0, len(typeAndConstraints))

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

func (g *newGenerator) extractConstraint(v cue.Value) (ast.TypeConstraint, error) {
	toConstraint := func(operator string, arg cue.Value) (ast.TypeConstraint, error) {
		scalar, err := cueConcreteToScalar(arg)
		if err != nil {
			return ast.TypeConstraint{}, err
		}

		return ast.TypeConstraint{
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
		return ast.TypeConstraint{}, errorWithCueRef(v, "unsupported op for number %v", op)
	}
}

func (g *newGenerator) declareList(v cue.Value) (*ast.Definition, error) {
	i, err := v.List()
	if err != nil {
		return nil, err
	}

	typeDef := &ast.Definition{
		Type:      ast.TypeArray,
		IndexType: ast.TypeInt64,

		// FIXME: we set a default type because our logic is broken
		ValueType: &ast.Definition{
			Type: ast.TypeAny,
		},

		Nullable: false,
		//SubType:     nil,
		//Constraints: nil,
	}

	// works only for a closed/concrete list
	if v.IsConcrete() {
		// TODO
		for i.Next() {
			node, err := g.declareNode(i.Value())
			if err != nil {
				return nil, err
			}

			typeDef.ValueType = node
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

	typeDef.ValueType = expr

	return typeDef, nil
}

package simplecue

import (
	"fmt"
	"math/bits"
	"sort"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/errors"
	"github.com/K-Phoen/grabana/internal/gen/simplecue/ts"
	tsast "github.com/K-Phoen/grabana/internal/gen/simplecue/ts/ast"
	"github.com/davecgh/go-spew/spew"
)

const (
	attrname        = "cuetsy"
	attrEnumMembers = "memberNames"
	attrKind        = "kind"
)

type FilePrinter func(file *File) ([]byte, error)

// ASTTarget strings indicate the kind of declaration to which a CUE
// value should be translated. They are used in both @cuetsy attributes, and in
// calls to certain methods.
type ASTTarget string

const (
	// TargetType targets conversion of a CUE value to a type declaration.
	TargetType ASTTarget = "type"

	// TargetEnum targets conversion of a CUE value to an `enum` declaration.
	TargetEnum ASTTarget = "enum"
)

var allTargets = [...]ASTTarget{
	TargetType,
	TargetEnum,
}

// Config governs certain variable behaviors when converting CUE to Typescript.
type Config struct {
	// Export determines whether generated TypeScript symbols are exported.
	Export bool

	// Package name used to generate code into.
	Package string
}

// GenerateAny takes a cue.Value and generates the corresponding code for all
// top-level members of that value.
//
// Hidden fields are ignored.
func GenerateAny(val cue.Value, c Config, printer FilePrinter) (b []byte, err error) {
	file, err := GenerateAST(val, c)
	if err != nil {
		return nil, err
	}

	return printer(file)
}

func GenerateAST(val cue.Value, c Config) (*File, error) {
	if err := val.Validate(); err != nil {
		return nil, err
	}

	g := &generator{
		c:   c,
		val: &val,
	}

	file := File{
		Package: c.Package,
	}

	iter, err := val.Fields(
		cue.Definitions(true),
		cue.Optional(true),
	)
	if err != nil {
		return nil, err
	}

	for iter.Next() {
		n := g.declare(iter.Selector().String(), iter.Value())
		if n == nil {
			continue
		}

		file.Types = append(file.Types, *n)
	}

	return &file, g.err
}

type generator struct {
	val *cue.Value
	c   Config
	err errors.Error
}

func (g *generator) addErr(err error) {
	if err != nil {
		g.err = errors.Append(g.err, errors.Promote(err, "generate failed"))
	}
}

func (g *generator) declare(name string, v cue.Value) *TypeDefinition {
	tst, err := getKindFor(v)
	if err != nil {
		// Ignore values without attributes
		return nil
	}

	fmt.Printf("Declaring kind '%s': '%s'\n", tst, name)

	switch tst {
	case TargetEnum:
		return g.genEnumDefinition(name, v)
	/*
		case TypeEnum:
			return g.genEnum(name, v)
		case TargetAlias:
			return g.genType(name, v)
	*/
	case TargetType:
		return g.genTypeDefinition(name, v)
	default:
		return nil // TODO error out
	}
}

type KV struct {
	K, V string
}

// genEnumDefinition turns the following cue values into enum definitions:
//   - value disjunction (a | b | c): values are taken as attribute memberNames,
//     if memberNames is absent, then keys implicitly generated as CamelCase
//   - string struct: struct keys get enum keys, struct values enum values
func (g *generator) genEnumDefinition(name string, v cue.Value) *TypeDefinition {
	fmt.Printf("→ genEnumDefinition()\n")
	// FIXME compensate for attribute-applying call to Unify() on incoming Value
	op, dvals := v.Expr()
	if op == cue.AndOp {
		v = dvals[0]
		op, _ = v.Expr()
	}

	// We restrict the expression of enums to ints or strings.
	allowed := cue.StringKind | cue.IntKind
	ik := v.IncompleteKind()
	if ik&allowed != ik {
		g.addErr(valError(v, "typescript enums may only be generated from concrete strings, or ints with memberNames attribute"))
		return nil
	}

	exprs, err := orEnumDef(v)
	if err != nil {
		g.addErr(err)
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
		Values:   exprs,
	}

	/*
		ret := make([]ts.Decl, 2)
		ret[0] = tsast.TypeDecl{
			Name:        ts.Ident(name),
			Type:        tsast.EnumType{Elems: exprs},
			CommentList: commentsFor(v, true),
			Export:      g.c.Export,
		}

		defaultIdent, err := enumDefault(v)
		g.addErr(err)

		if defaultIdent == nil {
			return ret[:1]
		}

		ret[1] = tsast.VarDecl{
			Names:  ts.Names("default" + name),
			Type:   ts.Ident(name),
			Value:  tsast.SelectorExpr{Expr: ts.Ident(name), Sel: *defaultIdent},
			Export: g.c.Export,
		}
		return ret
	*/
}

// genEnum turns the following cue values into typescript enums:
//   - value disjunction (a | b | c): values are taken as attribute memberNames,
//     if memberNames is absent, then keys implicitly generated as CamelCase
//   - string struct: struct keys get enum keys, struct values enum values
func (g *generator) genEnum(name string, v cue.Value) []ts.Decl {
	// FIXME compensate for attribute-applying call to Unify() on incoming Value
	op, dvals := v.Expr()
	if op == cue.AndOp {
		v = dvals[0]
		op, _ = v.Expr()
	}

	// We restrict the expression of TS enums to ints or strings.
	allowed := cue.StringKind | cue.IntKind
	ik := v.IncompleteKind()
	if ik&allowed != ik {
		g.addErr(valError(v, "typescript enums may only be generated from concrete strings, or ints with memberNames attribute"))
		return nil
	}

	exprs, err := orEnum(v)
	if err != nil {
		g.addErr(err)
	}

	ret := make([]ts.Decl, 2)
	ret[0] = tsast.TypeDecl{
		Name:        ts.Ident(name),
		Type:        tsast.EnumType{Elems: exprs},
		CommentList: commentsFor(v, true),
		Export:      g.c.Export,
	}

	defaultIdent, err := enumDefault(v)
	g.addErr(err)

	if defaultIdent == nil {
		return ret[:1]
	}

	ret[1] = tsast.VarDecl{
		Names:  ts.Names("default" + name),
		Type:   ts.Ident(name),
		Value:  tsast.SelectorExpr{Expr: ts.Ident(name), Sel: *defaultIdent},
		Export: g.c.Export,
	}
	return ret
}

func enumDefault(v cue.Value) (*tsast.Ident, error) {
	def, ok := v.Default()
	if !ok {
		return nil, def.Err()
	}

	if v.IncompleteKind() == cue.StringKind {
		s, _ := def.String()
		return &tsast.Ident{Name: strings.Title(s)}, nil
	}

	// For Int, Float, Numeric we need to find the default value and its corresponding memberName value
	a := v.Attribute(attrname)
	val, found, err := a.Lookup(0, attrEnumMembers)
	if err != nil || !found {
		return nil, valError(v, "Looking for memberNames: found=%t err=%s", found, err)
	}
	evals := strings.Split(val, "|")

	_, dvals := v.Expr()
	for i, val := range dvals {
		valLab, _ := val.Label()
		defLab, _ := def.Label()
		if valLab == defLab {
			return &tsast.Ident{Name: evals[i]}, nil
		}
	}

	// should never reach here tho
	return nil, valError(v, "unable to find memberName corresponding to the default")
}

// List the pairs of values and member names in an enum. Err if input is not an enum
func enumPairs(v cue.Value) ([]enumPair, error) {
	// TODO should validate here. Or really, this is just evidence of how building these needs its own types
	op, dvals := v.Expr()
	if !targetsKind(v, TargetEnum) || op != cue.OrOp {
		return nil, fmt.Errorf("not an enum: %v (%s)", v, v.Path())
	}

	a := v.Attribute(attrname)
	val, found, err := a.Lookup(0, attrEnumMembers)
	if err != nil {
		return nil, valError(v, "Looking for memberNames: found=%t err=%s", found, err)
	}

	var evals []string
	if found {
		evals = strings.Split(val, "|")
	} else if v.IncompleteKind() == cue.StringKind {
		for _, part := range dvals {
			s, _ := part.String()
			evals = append(evals, strings.Title(s))
		}
	} else {
		return nil, valError(v, "must provide memberNames attribute for non-string enums")
	}

	var pairs []enumPair
	for i, eval := range evals {
		pairs = append(pairs, enumPair{
			name: eval,
			val:  dvals[i],
		})
	}

	return pairs, nil
}

type enumPair struct {
	name string
	val  cue.Value
}

func orEnum(v cue.Value) ([]ts.Expr, error) {
	_, dvals := v.Expr()
	a := v.Attribute(attrname)

	var attrMemberNameExist bool
	var evals []string
	if a.Err() == nil {
		val, found, err := a.Lookup(0, attrEnumMembers)
		if err == nil && found {
			attrMemberNameExist = true
			evals = strings.Split(val, "|")
			if len(evals) != len(dvals) {
				return nil, valError(v, "typescript enums and %s attributes size doesn't match", attrEnumMembers)
			}
		}
	}

	// We only allowed String Enum to be generated without memberName attribute
	if v.IncompleteKind() != cue.StringKind && !attrMemberNameExist {
		return nil, valError(v, "typescript numeric enums may only be generated from memberNames attribute")
	}

	var fields []ts.Expr
	for idx, dv := range dvals {
		var text string
		var id tsast.Ident
		if attrMemberNameExist {
			text = evals[idx]
			id = ts.Ident(text)
		} else {
			text, _ = dv.String()
			id = ts.Ident(strings.Title(text))
		}

		if !dv.IsConcrete() {
			return nil, valError(v, "typescript enums may only be generated from a disjunction of concrete strings or numbers")
		}

		if id.Validate() != nil {
			return nil, valError(v, "title casing of enum member %q produces an invalid typescript identifier; memberNames must be explicitly given in @cuetsy attribute", text)
		}

		val, err := tsprintConcrete(dv)
		if err != nil {
			return nil, err
		}
		fields = append(fields, tsast.AssignExpr{
			// Simple mapping of all enum values (which we are assuming are in
			// lowerCamelCase) to corresponding CamelCase
			Name:  id,
			Value: val,
		})
	}

	sort.Slice(fields, func(i, j int) bool {
		return fields[i].String() < fields[j].String()
	})

	return fields, nil
}

func orEnumDef(v cue.Value) ([]EnumValue, error) {
	_, dvals := v.Expr()
	a := v.Attribute(attrname)

	var attrMemberNameExist bool
	var evals []string
	if a.Err() == nil {
		val, found, err := a.Lookup(0, attrEnumMembers)
		if err == nil && found {
			attrMemberNameExist = true
			evals = strings.Split(val, "|")
			if len(evals) != len(dvals) {
				return nil, valError(v, "enums and %s attributes size doesn't match", attrEnumMembers)
			}
		}
	}

	// We only allowed String Enum to be generated without memberName attribute
	if v.IncompleteKind() != cue.StringKind && !attrMemberNameExist {
		return nil, valError(v, "numeric enums may only be generated from memberNames attribute")
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
			return nil, valError(v, "enums may only be generated from a disjunction of concrete strings or numbers")
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

func (g *generator) genTypeDefinition(name string, v cue.Value) *TypeDefinition {
	fmt.Printf("→ genTypeDefinition()\n")
	// We restrict the derivation of type definitions to struct kinds.
	// (More than just a struct literal match this, though.)
	if v.IncompleteKind() != cue.StructKind {
		// FIXME check for bottom here, give different error
		g.addErr(valError(v, "type definitions may only be generated from structs"))
		return nil
	}

	var elems []FieldDefinition

	iter, _ := v.Fields(cue.Optional(true), cue.Definitions(true))
	for iter != nil && iter.Next() {
		fmt.Println("iter")
		if iter.Selector().PkgPath() != "" {
			g.addErr(valError(iter.Value(), "cannot generate hidden fields; typescript has no corresponding concept"))
			return nil
		}

		fieldDef, err := g.genFieldDefinition(iter.Value())
		if err != nil || fieldDef == nil {
			fmt.Printf("skipping fieldDef(%s): err(%+v), fieldDef==nil(%+v)\n", iter.Label(), err, fieldDef == nil)

			continue // TODO
			//return nil
		}

		fieldDef.Required = !iter.IsOptional()

		elems = append(elems, *fieldDef)
	}

	return &TypeDefinition{
		Type:     DefinitionStruct,
		Name:     name,
		Comments: commentsFromCueValue(v),
		Fields:   elems,
	}
}

// Generate a typeRef for the cue.Value
func (g *generator) genFieldDefinition(v cue.Value) (*FieldDefinition, error) {
	fmt.Printf("→ genFieldDefinition()\n")
	label, ok := v.Label()
	if !ok {
		return nil, fmt.Errorf("could not extract field label")
	}

	var fieldTypeDef *FieldType
	var err error
	fieldDef := &FieldDefinition{
		Name: label,
	}

	if hasEnumReference(v) {
		fieldTypeDef, err = g.genEnumReferenceDef(v)
	} else {
		fieldTypeDef, err = g.simplifiedFieldType(v, true, false)
	}

	if err != nil {
		if !containsCuetsyReference(v) {
			g.addErr(err)
			return nil, err
		}
		g.addErr(err)
		return nil, nil
	}

	if fieldTypeDef == nil {
		return nil, fmt.Errorf("field type not generated")
	}

	fieldDef.Type = *fieldTypeDef

	return fieldDef, err
}

func hasEnumReference(v cue.Value) bool {
	// Check if we've got an enum reference at top depth or one down. If we do, it
	// changes how we generate.
	hasPred := containsPred(v, 1,
		isReference,
		func(v cue.Value) bool { return targetsKind(cue.Dereference(v), TargetEnum) },
	)

	// Check if it's setting an enum value [Enum & "value"]
	op, args := v.Expr()
	if op == cue.AndOp {
		return hasPred
	}

	// Check if it has default value [Enum & (*"defaultValue" | _)]
	for _, a := range args {
		if a.IncompleteKind() == cue.TopKind {
			return hasPred
		}
	}

	isUnion := true
	allEnums := true
	for _, a := range args {
		// Check if it is a union [(Enum & "a") | (Enum & "b")]
		if a.Kind() != a.IncompleteKind() {
			isUnion = false
		}
		// Check if all elements are enums
		_, exprs := a.Expr()
		for _, e := range exprs {
			if t, err := getKindFor(cue.Dereference(e)); err == nil && t != TargetEnum {
				allEnums = false
			}
		}
	}

	return hasPred && isUnion && allEnums
}

func hasTypeReference(v cue.Value) bool {
	hasTypeRef := containsCuetsyReference(v, TargetType)
	// Check if it's setting an enum value [Type & "value"]
	op, args := v.Expr()
	if op == cue.AndOp || op == cue.SelectorOp {
		return hasTypeRef
	}

	// Check if it has default value [Type & (*"defaultValue" | _)]
	for _, a := range args {
		if a.IncompleteKind() == cue.TopKind {
			return hasTypeRef
		}
	}

	return false
}

// Generate a typeref for a value that refers to a field
// TODO: rewrite
func (g *generator) genEnumReferenceDef(v cue.Value) (*FieldType, error) {
	fmt.Printf(" → genEnumReferenceDef()\n")
	var lit *cue.Value

	conjuncts := appendSplit(nil, cue.AndOp, v)
	var enumUnions map[cue.Value]cue.Value
	switch len(conjuncts) {
	case 0:
		ve := valError(v, "unreachable: no conjuncts while looking for enum references")
		g.addErr(ve)
		return nil, ve
	case 1:
		_, dvals := v.Expr()

		// Have to do attribute checks on the referenced field itself, so deref
		deref := cue.Dereference(v)

		var dstr string
		if len(dvals) > 1 {
			dstr, _ = dvals[1].String()
		}
		if _, ok := dvals[1].Source().(*ast.Ident); ok && targetsKind(deref) {
			return &FieldType{Type: TypeID(dstr)}, nil
			//return ts.Ident(dstr), nil
		}
		// This case is when we have a union of enums which we need to iterate them to get their values or has a default value.
		// It retrieves a list of literals with their references.

		var err error
		enumUnions, err = g.findEnumUnions(v)
		if err != nil {
			return nil, err
		}
	case 2:
		var err error
		conjuncts[1] = getDefaultEnumValue(conjuncts[1])
		lit, err = getEnumLiteral(conjuncts)
		if err != nil {
			ve := valError(v, err.Error())
			g.addErr(ve)
			return nil, ve
		}
	case 3:
		if conjuncts[1].IncompleteKind() == cue.TopKind {
			conjuncts[1] = conjuncts[0]
		}

		if !conjuncts[0].Equals(conjuncts[1]) && conjuncts[0].Subsume(conjuncts[1]) != nil {
			ve := valError(v, "complex unifications containing references to enums without overriding parent are not currently supported")
			g.addErr(ve)
			return nil, ve
		}

		var err error
		lit, err = getEnumLiteral(conjuncts[1:])
		if err != nil {
			ve := valError(v, err.Error())
			g.addErr(ve)
			return nil, ve
		}
	default:
		ve := valError(v, "complex unifications containing references to enums are not currently supported")
		g.addErr(ve)
		return nil, ve
	}

	// Search the expr tree for the actual enum. This approach is uncomfortable
	// without having the assurance that there aren't more than one possible match/a
	// guarantee from the CUE API of a stable, deterministic search order, etc.
	enumValues, referrer, has := findRefWithKind(v, TargetEnum)
	if !has {
		ve := valError(v, "does not reference a field with a cuetsy enum attribute")
		g.addErr(ve)
		return nil, fmt.Errorf("no enum attr in %s", v)
	}

	var err error
	decls := g.genEnum("foo", enumValues)
	ref := &typeRef{}

	// Construct the type component of the reference
	switch len(decls) {
	default:
		ve := valError(v, "unsupported number of expression args (%v) in reference, expected 1 or 2", len(decls))
		g.addErr(ve)
		return nil, ve
	case 1, 2:
		ref.T, err = referenceValueAs(referrer, TargetEnum)
		if err != nil {
			return nil, err
		}
	}

	// Either specify a default if one exists (one conjunct), or rewrite the type to
	// reference one of the members of the enum (two conjuncts).
	switch len(conjuncts) {
	case 1:
		if defv, hasdef := v.Default(); hasdef {
			err = g.findIdent(v, enumValues, defv, func(expr tsast.Ident) {
				ref.D = tsast.SelectorExpr{Expr: ref.T, Sel: expr}
			})
		}
		if len(enumUnions) == 0 {
			spew.Dump("enumUnions empty")
			break
		}
		var elements []tsast.Expr
		for lit, enumValues := range enumUnions {
			err = g.findIdent(v, enumValues, lit, func(ident tsast.Ident) {
				elements = append(elements, tsast.SelectorExpr{
					Expr: ref.T,
					Sel:  ident,
				})
			})
		}

		// To avoid to change the order of the elements everytime that we generate the code.
		sort.Slice(elements, func(i, j int) bool {
			return elements[i].String() < elements[j].String()
		})

		ref.T = ts.Union(elements...)
		spew.Dump(ref.T.String())
	case 2, 3:
		var rr tsast.Expr
		err = g.findIdent(v, enumValues, *lit, func(ident tsast.Ident) {
			rr = tsast.SelectorExpr{Expr: ref.T, Sel: ident}
		})

		op, args := v.Expr()
		hasInnerDefault := false
		if len(args) == 2 && op == cue.AndOp {
			_, hasInnerDefault = args[1].Default()
		}

		if _, has := v.Default(); has || hasInnerDefault {
			ref.D = rr
		} else {
			ref.T = rr
		}
	}

	return nil, nil
	//return ref, err
}

// findEnumUnions find the unions between enums like (#Enum & "a") | (#Enum & "b")
func (g generator) findEnumUnions(v cue.Value) (map[cue.Value]cue.Value, error) {
	op, values := v.Expr()
	if op != cue.OrOp {
		return nil, nil
	}

	enumsWithUnions := make(map[cue.Value]cue.Value, len(values))
	for _, val := range values {
		conjuncts := appendSplit(nil, cue.AndOp, val)
		if len(conjuncts) != 2 {
			return nil, nil
		}
		cr, lit := conjuncts[0], conjuncts[1]
		if cr.Subsume(lit) != nil {
			return nil, nil
		}

		switch val.Kind() {
		case cue.StringKind, cue.IntKind:
			enumValues, _, has := findRefWithKind(v, TargetEnum)
			if !has {
				return nil, nil
			}
			enumsWithUnions[lit] = enumValues
		default:
			_, vals := val.Expr()
			if len(vals) > 1 {
				return nil, valError(v, "%s.%s isn't a valid enum value", val.Path().String(), vals[1])
			}
			return nil, valError(v, "Invalid value in path %s", val.Path().String())
		}
	}

	return enumsWithUnions, nil
}

func (g generator) findIdent(v, ev, tv cue.Value, fn func(tsast.Ident)) error {
	if ev.Subsume(tv) != nil {
		err := valError(v, "may only apply values to an enum that are members of that enum; %#v is not a member of %#v", tv, ev)
		g.addErr(err)
		return err
	}
	pairs, err := enumPairs(ev)
	if err != nil {
		return err
	}
	for _, pair := range pairs {
		if veq(pair.val, tv) {
			fn(tsast.Ident{Name: pair.name})
			return nil
		}
	}

	// unreachable?
	return valError(v, "%#v not equal to any member of %#v, but should have been caught by subsume check", tv, ev)
}

func getEnumLiteral(conjuncts []cue.Value) (*cue.Value, error) {
	var lit *cue.Value
	// The only case we actually want to support, at least for now, is this:
	//
	//   enum: "foo" | "bar" @cuetsy(kind="enum")
	//   enumref: enum & "foo" @cuetsy(kind="type")
	//
	// Where we render enumref to TS as `Enumref: Enum.Foo`.
	// For that case, we allow at most two conjuncts, and make sure they
	// fit the pattern of the two operands above.
	aref, bref := isReference(conjuncts[0]), isReference(conjuncts[1])
	aconc, bconc := conjuncts[0].IsConcrete(), conjuncts[1].IsConcrete()
	var cr cue.Value
	if aref {
		cr, lit = conjuncts[0], &(conjuncts[1])
	} else {
		cr, lit = conjuncts[1], &(conjuncts[0])
	}

	if aref == bref || aconc == bconc || cr.Subsume(*lit) != nil {
		return nil, errors.New(fmt.Sprintf("may only unify a referenced enum with a concrete literal member of that enum. Path: %s", conjuncts[0].Path()))
	}

	return lit, nil
}

// getDefaultEnumValue is looking for default values like #Enum & (*"default" | _) struct
func getDefaultEnumValue(v cue.Value) cue.Value {
	if v.IncompleteKind() != cue.TopKind {
		return v
	}

	op, args := v.Expr()
	if op != cue.OrOp {
		return v
	}

	for _, a := range args {
		if a.IncompleteKind() == cue.TopKind {
			if def, has := a.Default(); has {
				return def
			}
		}
	}
	return v
}

// typeRef is a pair of expressions for referring to another type - the reference
// to the type, and the default value for the referrer. The default value
// may be the one provided by either the referent, or by the field doing the referring
// (in the case of a superseding mark).
type typeRef struct {
	T ts.Expr
	D ts.Expr
}

// Render a string containing a Typescript semantic equivalent to the provided
// Value for placement in a single field, if possible.
func (g generator) simplifiedFieldType(v cue.Value, isType bool, isDefault bool) (*FieldType, error) {
	fmt.Printf(" → simplifiedFieldType()\n")
	if hasEnumReference(v) {
		fmt.Printf(" → hasEnumReference(): true\n")
		return g.genEnumReferenceDef(v)
	}

	// References are orthogonal to the Kind system. Handle them first.
	if hasTypeReference(v) || containsCuetsyReference(v, TargetType) {
		fmt.Printf(" → hasTypeReference() || containsCuetsyReference(v, TargetType): true\n")
		refType, err := typeFromReference(v)
		if err != nil {
			return nil, err
		}
		if refType != nil {
			return refType, nil
		}

		return nil, valError(v, "failed to generate reference correctly for path %s", v.Path().String())
	}

	verr := v.Validate(cue.Final())
	if verr != nil {
		spew.Dump("invalid")
		return nil, verr
	}

	op, dvals := v.Expr()
	// Eliminate concretes first, to make handling the others easier.

	// Concrete values.
	// Includes "foobar", 5, [1,2,3], etc. (literal values)
	k := v.Kind()
	spew.Dump("kind", k)
	switch k {
	/*
		case cue.StructKind:
			switch op {
			case cue.SelectorOp, cue.AndOp, cue.NoOp:
				// Checks [string]something only.
				// It skips structs like {...} (cue.TopKind) to avoid undesired results.
				val := v.LookupPath(cue.MakePath(cue.AnyString))
				if val.Exists() && val.IncompleteKind() != cue.TopKind {
					expr, err := g.tsprintField(val, isType, isDefault)
					if err != nil {
						return nil, valError(v, err.Error())
					}
					kvs := []tsast.KeyValueExpr{
						{
							Key:         ts.Ident("string"),
							Value:       expr,
							CommentList: commentsFor(val.Value(), true),
						},
					}
					return tsast.ObjectLit{Elems: kvs, IsType: isType, IsMap: true}, nil
				}

				iter, err := v.Fields(cue.Optional(true))
				if err != nil {
					return nil, valError(v, "something went wrong when generate nested structs")
				}
				size, _ := v.Len().Int64()
				kvs := make([]tsast.KeyValueExpr, 0, size)
				for iter.Next() {
					value, ok := shouldIterateValue(iter.Value(), isDefault)
					if !ok {
						continue
					}
					expr, err := g.tsprintField(value, isType, isDefault)
					if err != nil {
						return nil, valError(v, err.Error())
					}
					k := iter.Label()
					if iter.IsOptional() {
						k += "?"
					}
					kvs = append(kvs, tsast.KeyValueExpr{
						Key:         ts.Ident(k),
						Value:       expr,
						CommentList: commentsFor(iter.Value(), true),
					})
				}

				return tsast.ObjectLit{Elems: kvs, IsType: isType}, nil
			default:
				return nil, valError(v, "not expecting op type %d", op)
			}

	*/
	case cue.ListKind:
		// A list is concrete (and thus its complete kind is ListKind instead of
		// BottomKind) if it specifies a finite number of elements - is
		// "closed". This is independent of the types of its elements, which may
		// be anywhere on the concreteness spectrum.
		//
		// For closed lists, we simply iterate over its component elements and
		// print their typescript representation.
		iter, _ := v.List()
		var elems []FieldType
		for iter.Next() {
			e, err := g.simplifiedFieldType(iter.Value(), isType, isDefault)
			if err != nil {
				return nil, err
			}
			elems = append(elems, *e)
		}

		// TODO
		return &FieldType{
			Type:        TypeArray,
			SubType:     elems,
			Constraints: nil,
		}, nil
	case cue.StringKind, cue.BoolKind, cue.FloatKind, cue.IntKind:
		fType, err := concreteScalarType(v)
		if err != nil {
			return nil, err
		}

		return &fType, nil

	case cue.BytesKind:
		return &FieldType{
			Type: TypeBytes,
		}, nil
	}

	// Handler for disjunctions
	disj := func(dvals []cue.Value) ([]FieldType, error) {
		parts := make([]FieldType, 0, len(dvals))
		for _, dv := range dvals {
			p, err := g.simplifiedFieldType(dv, isType, isDefault)
			if err != nil {
				return nil, err
			}
			parts = append(parts, *p)
		}
		return parts, nil
	}

	// Others: disjunctions, etc.
	ik := v.IncompleteKind()
	spew.Dump("incomplete kind", ik)
	switch ik {
	case cue.BottomKind:
		return nil, valError(v, "bottom, unsatisfiable")
	case cue.ListKind:
		// This list is open - its final element is ...<value> - and we can only
		// meaningfully convert open lists to typescript if there are zero other
		// elements.

		// First, peel off a simple default, if one exists.
		// dlist, has := v.Default()
		// if has && op == cue.OrOp {
		// 	di := analyzeList(dlist)
		// 	if len(dvals) != 2 {
		// 		panic(fmt.Sprintf("%v branches on list disjunct, can only handle 2", len(dvals)))
		// 	}
		// 	if di.eq(analyzeList(dvals[1])) {
		// 		v = dvals[0]
		// 	} else if di.eq(analyzeList(dvals[0])) {
		// 		v = dvals[1]
		// 	} else {
		// 		panic("wat - list kind had default but analysis did not match for either disjunct branch")
		// 	}
		// }

		// If the default (all lists have a default, usually self, ugh) differs from the
		// input list, peel it off. Otherwise our AnyIndex lookup may end up getting
		// sent on the wrong path.
		defv, _ := v.Default()
		if !defv.Equals(v) {
			v = dvals[0]
		}

		if op == cue.OrOp {
			subTypes, err := disj(dvals)
			if err != nil {
				return nil, err
			}

			return &FieldType{
				Type:    TypeDisjunction,
				SubType: subTypes,
			}, nil
		}

		e := v.LookupPath(cue.MakePath(cue.AnyIndex))
		if e.Exists() {
			expr, err := g.simplifiedFieldType(e, isType, isDefault)
			if err != nil {
				return nil, err
			}

			// TODO
			return &FieldType{
				Type:    TypeArray,
				SubType: []FieldType{*expr},
			}, nil
		} else {
			// unreachable?
			return nil, errors.New("open list must have a type")
		}
	case cue.NumberKind, cue.StringKind:
		// It appears there are only three cases in which we can have an
		// incomplete NumberKind or StringKind:
		//
		// 1. The corresponding literal is a bounding constraint (which subsumes
		// both int and float), e.g. >2.2, <"foo"
		// 2. There's a disjunction of concrete literals of the relevant type
		// 2. The corresponding literal is the basic type "number" or "string"
		//
		// The first case has no equivalent in typescript, and entails we error
		// out. The other two have the same handling as for other kinds, so we
		// fall through. We disambiguate by seeing if there is an expression
		// (other than Or, "|"), which is how ">" and "2.2" are represented.
		//
		// TODO get more certainty/a clearer way of ascertaining this
		switch op {
		case cue.RegexMatchOp:
			return &FieldType{
				Type:        TypeString,
				Constraints: []TypeConstraint{{Op: "regex", Args: []any{op.String()}}},
			}, nil
		case cue.NoOp, cue.OrOp, cue.AndOp:
		default:
			return nil, valError(v, "bounds constraints are not supported as they lack a direct typescript equivalent")
		}
		fallthrough
	case cue.NullKind:
		// It evaluates single null value
		if op == cue.NoOp && len(dvals) == 0 {
			return &FieldType{
				Type: TypeNull,
			}, nil
		}
		fallthrough
	case cue.FloatKind, cue.IntKind, cue.BoolKind, cue.StructKind:
		// Having eliminated the possibility of bounds/constraints, we're left
		// with disjunctions and basic types.
		switch op {
		case cue.OrOp:
			if len(dvals) == 2 && dvals[0].Kind() == cue.NullKind {
				return g.simplifiedFieldType(dvals[1], isType, isDefault)
			}

			subTypes, err := disj(dvals)
			if err != nil {
				return nil, err
			}

			return &FieldType{
				Type:    TypeDisjunction,
				SubType: subTypes,
			}, nil
		case cue.AndOp:
		// There's no op for simple unification; it's a basic type, and can
		// be trivially rendered.
		case cue.NoOp:
			// Something a list of two items like #Enum & "default" struct reaches this point.
			// The problem is that "default" is not detected as value, only by default, and we need
			// to add this value manually.
			if args := getValuesWithDefaults(v, dvals[0]); args != nil {
				subTypes, err := disj(args)
				if err != nil {
					return nil, err
				}

				return &FieldType{
					Type:    TypeDisjunction,
					SubType: subTypes,
				}, nil
			}
		default:
			return nil, valError(v, "no handler for operator: '%s' for kind '%s'", op.String(), ik)
		}
		fallthrough

	case cue.TopKind:
		fType := scalarFieldType(ik)
		return &fType, nil
	case cue.BytesKind:
		return &FieldType{
			Type: TypeBytes,
		}, nil
	}
	// Having more than one possible kind entails a disjunction, TopKind, or
	// NumberKind. We've already eliminated TopKind and NumberKind, so now check
	// if there's more than one bit set. (If there isn't, it's a bug: we've
	// missed a kind above). If so, run our disjunction-handling logic.
	if bits.OnesCount16(uint16(ik)) > 1 {
		subTypes, err := disj(dvals)
		if err != nil {
			return nil, err
		}

		return &FieldType{
			Type:    TypeDisjunction,
			SubType: subTypes,
		}, nil
	}

	return nil, valError(v, "unrecognized kind %v", ik)
}

func getValuesWithDefaults(v cue.Value, cuetsyType cue.Value) []cue.Value {
	if def, ok := v.Default(); ok {
		op, _ := cuetsyType.Expr()
		if op == cue.SelectorOp {
			return []cue.Value{cuetsyType, def}
		}
	}

	return nil
}

// ONLY call this function if it has been established that the provided Value is
// Concrete.
func tsprintConcrete(v cue.Value) (ts.Expr, error) {
	switch v.Kind() {
	case cue.NullKind:
		return ts.Null(), nil
	case cue.StringKind:
		s, _ := v.String()
		return ts.Str(s), nil
	case cue.FloatKind:
		f, _ := v.Float64()
		return ts.Float(f), nil
	case cue.NumberKind, cue.IntKind:
		i, _ := v.Int64()
		return ts.Int(i), nil
	case cue.BoolKind:
		b, _ := v.Bool()
		return ts.Bool(b), nil
	default:
		return nil, valError(v, "concrete kind not found: %s", v.Kind())
	}
}

// ONLY call this function if it has been established that the provided Value is
// Concrete.
func cueConcreteToScalar(v cue.Value) (interface{}, error) {
	switch v.Kind() {
	case cue.NullKind:
		return nil, nil
	case cue.StringKind:
		return v.String()
	case cue.FloatKind:
		return v.Float64()
	case cue.NumberKind, cue.IntKind:
		return v.Int64()
	case cue.BoolKind:
		return v.Bool()
	default:
		return nil, valError(v, "concrete kind not found: %s", v.Kind())
	}
}

// ONLY call this function if it has been established that the provided Value is
// Concrete.
func concreteScalarType(v cue.Value) (FieldType, error) {
	switch v.Kind() {
	case cue.NullKind:
		return FieldType{Type: TypeNull, Constraints: []TypeConstraint{{
			Op:   "eq",
			Args: []any{nil},
		}}}, nil
	case cue.StringKind:
		s, _ := v.String()

		return FieldType{Type: TypeString, Constraints: []TypeConstraint{{
			Op:   "eq",
			Args: []any{s},
		}}}, nil
	case cue.NumberKind, cue.FloatKind:
		f, _ := v.Float64()

		return FieldType{Type: TypeFloat64, Constraints: []TypeConstraint{{
			Op:   "eq",
			Args: []any{f},
		}}}, nil
	case cue.IntKind:
		i, _ := v.Int64()

		return FieldType{Type: TypeInt64, Constraints: []TypeConstraint{{
			Op:   "eq",
			Args: []any{i},
		}}}, nil
	case cue.BoolKind:
		b, _ := v.Bool()

		return FieldType{Type: TypeBool, Constraints: []TypeConstraint{{
			Op:   "eq",
			Args: []any{b},
		}}}, nil
	default:
		return FieldType{}, valError(v, "concrete kind not found: %s", v.Kind())
	}
}

func scalarFieldType(k cue.Kind) FieldType {
	switch k {
	case cue.BoolKind:
		return FieldType{Type: TypeBool}
	case cue.StringKind:
		return FieldType{Type: TypeString}
	case cue.IntKind:
		// TODO: better identification of actual type
		return FieldType{Type: TypeInt64}
	case cue.NumberKind, cue.FloatKind:
		return FieldType{Type: TypeFloat64}
	case cue.TopKind:
		return FieldType{Type: TypeAny}
	case cue.NullKind:
		return FieldType{Type: TypeNull}
	default:
		return FieldType{Type: TypeUnknown}
	}
}

func tsprintType(k cue.Kind) ts.Expr {
	switch k {
	case cue.BoolKind:
		return ts.Ident("boolean")
	case cue.StringKind:
		return ts.Ident("string")
	case cue.NumberKind, cue.FloatKind, cue.IntKind:
		return ts.Ident("number")
	case cue.TopKind:
		return ts.Ident("unknown")
	case cue.NullKind:
		return ts.Ident("null")
	default:
		return nil
	}
}

func valError(v cue.Value, format string, args ...interface{}) error {
	s := v.Source()
	if s == nil {
		return fmt.Errorf(format, args...)
	}

	msg := ""
	if i, ok := s.(*ast.Field); ok {
		msg = fmt.Sprintf("Found an error in the field '%s:%d:%d'. ", i.Label, s.Pos().Line(), s.Pos().Column())
	}
	f := fmt.Sprintf("%sError: %s", msg, format)
	return errors.Newf(s.Pos(), f, args...)
}

func refAsInterface(v cue.Value) (ts.Expr, error) {
	// Bail out right away if the value isn't a reference
	op, dvals := v.Expr()
	if !isReference(v) && op != cue.SelectorOp {
		return nil, fmt.Errorf("not a reference")
	}

	// Have to do attribute checks on the referenced field itself, so deref
	deref := cue.Dereference(v)
	dstr, _ := dvals[1].String()

	// FIXME It's horrifying, teasing out the type of selector kinds this way. *Horrifying*.
	switch dvals[0].Source().(type) {
	case nil:
		// A nil subject means an unqualified selector (no "."
		// literal).  This can only possibly be a reference to some
		// sibling or parent of the top-level Value being generated.
		// (We can't do cycle detection with the meager tools
		// exported in cuelang.org/go/cue, so all we have for the
		// parent case is hopium.)
		if _, ok := dvals[1].Source().(*ast.Ident); ok && targetsKind(deref, TargetType) {
			return ts.Ident(dstr), nil
		}
	case *ast.SelectorExpr:
		if targetsKind(deref, TargetType) {
			return ts.Ident(dstr), nil
		}
	case *ast.Ident:
		if targetsKind(deref, TargetType) {
			str, ok := dvals[0].Source().(fmt.Stringer)
			if !ok {
				return nil, valError(v, "expected dvals[0].Source() to implement String()")
			}

			return tsast.SelectorExpr{
				Expr: ts.Ident(str.String()),
				Sel:  ts.Ident(dstr),
			}, nil
		}
	default:
		return nil, valError(v, "unknown selector subject type %T, cannot translate", dvals[0].Source())
	}

	return nil, nil
}

// typeFromReference returns the string that should be used to create a Typescript
// reference to the given struct, if a reference is allowable.
//
// References are only permitted to other Values with an @cuetsy(kind)
// attribute. The variadic parameter determines which kinds will be treated as
// permissible. By default, all kinds are permitted.
//
// A nil expr indicates a reference is not allowable, including the case
// that the provided Value is not actually a reference. A non-nil error
// indicates a deeper problem.
func typeFromReference(v cue.Value, kinds ...ASTTarget) (*FieldType, error) {
	// Bail out right away if there's no reference anywhere in the value.
	// if !containsReference(v) {
	// 	return nil, nil
	// }
	// End goal: we want to render a reference appropriately.
	// If the top-level is a reference, then this is simple.
	//
	// If the top-level merely contains a reference, this is harder.
	// - Let's start by only supporting that case when it's because there's a default.

	// Calling Expr peels off all default paths.
	op, dvals := v.Expr()
	_ = op

	if !isReference(v) {
		_, has := v.Default()
		if hasOverrideValues(v) {
			v = dvals[1]
		} else if !has || !isReference(dvals[0]) {
			return nil, valError(v, "references within complex logic are currently unsupported")
		} else {
			v = dvals[0]
		}

		// This may break a bunch of things but let's see if it gives us a
		// defensible baseline
		op, dvals = v.Expr()
	}

	var dstr string
	if len(dvals) > 1 {
		dstr, _ = dvals[1].String()
	}

	// Have to do attribute checks on the referenced field itself, so deref
	deref := cue.Dereference(v)

	// FIXME It's horrifying, teasing out the type of selector kinds this way. *Horrifying*.
	switch dvals[0].Source().(type) {
	case nil:
		// A nil subject means an unqualified selector (no "."
		// literal).  This can only possibly be a reference to some
		// sibling or parent of the top-level Value being generated.
		// (We can't do cycle detection with the meager tools
		// exported in cuelang.org/go/cue, so all we have for the
		// parent case is hopium.)
		if _, ok := dvals[1].Source().(*ast.Ident); ok && targetsKind(deref, kinds...) {
			return &FieldType{Type: TypeID(dstr)}, nil
			//return ts.Ident(dstr), nil
		}
	case *ast.SelectorExpr:
		if targetsKind(deref, kinds...) {
			return &FieldType{Type: TypeID(dstr)}, nil
			//return ts.Ident(dstr), nil
		}
	case *ast.Ident:
		if targetsKind(deref, kinds...) {
			str, ok := dvals[0].Source().(fmt.Stringer)
			if !ok {
				return nil, valError(v, "expected dvals[0].Source() to implement String()")
			}

			return &FieldType{Type: TypeID(str.String())}, nil
			// TODO
			/*
				return tsast.SelectorExpr{
					Expr: ts.Ident(str.String()),
					Sel:  ts.Ident(dstr),
				}, nil
			*/
		}

		// It happens when we are overriding a Type parent with a default value. Because `hasOverrides` is true,
		// dstr is the default value, and we need to set the Type name here
		if str, ok := dvals[0].Source().(fmt.Stringer); ok {
			return &FieldType{Type: TypeID(str.String())}, nil
			//return ts.Ident(str.String()), nil
		}
	default:
		return nil, valError(v, "unknown selector subject type %T, cannot translate path %s", dvals[0].Source(), v.Path().String())
	}

	return nil, nil
}

// referenceValueAs returns the string that should be used to create a Typescript
// reference to the given struct, if a reference is allowable.
//
// References are only permitted to other Values with an @cuetsy(kind)
// attribute. The variadic parameter determines which kinds will be treated as
// permissible. By default, all kinds are permitted.
//
// A nil expr indicates a reference is not allowable, including the case
// that the provided Value is not actually a reference. A non-nil error
// indicates a deeper problem.
func referenceValueAs(v cue.Value, kinds ...ASTTarget) (ts.Expr, error) {
	// Bail out right away if there's no reference anywhere in the value.
	// if !containsReference(v) {
	// 	return nil, nil
	// }
	// End goal: we want to render a reference appropriately in Typescript.
	// If the top-level is a reference, then this is simple.
	//
	// If the top-level merely contains a reference, this is harder.
	// - Let's start by only supporting that case when it's because there's a default.

	// Calling Expr peels off all default paths.
	op, dvals := v.Expr()
	_ = op

	if !isReference(v) {
		_, has := v.Default()
		if hasOverrideValues(v) {
			v = dvals[1]
		} else if !has || !isReference(dvals[0]) {
			return nil, valError(v, "references within complex logic are currently unsupported")
		} else {
			v = dvals[0]
		}

		// This may break a bunch of things but let's see if it gives us a
		// defensible baseline
		op, dvals = v.Expr()
	}

	var dstr string
	if len(dvals) > 1 {
		dstr, _ = dvals[1].String()
	}

	// Have to do attribute checks on the referenced field itself, so deref
	deref := cue.Dereference(v)

	// FIXME It's horrifying, teasing out the type of selector kinds this way. *Horrifying*.
	switch dvals[0].Source().(type) {
	case nil:
		// A nil subject means an unqualified selector (no "."
		// literal).  This can only possibly be a reference to some
		// sibling or parent of the top-level Value being generated.
		// (We can't do cycle detection with the meager tools
		// exported in cuelang.org/go/cue, so all we have for the
		// parent case is hopium.)
		if _, ok := dvals[1].Source().(*ast.Ident); ok && targetsKind(deref, kinds...) {
			return ts.Ident(dstr), nil
		}
	case *ast.SelectorExpr:
		if targetsKind(deref, kinds...) {
			return ts.Ident(dstr), nil
		}
	case *ast.Ident:
		if targetsKind(deref, kinds...) {
			str, ok := dvals[0].Source().(fmt.Stringer)
			if !ok {
				return nil, valError(v, "expected dvals[0].Source() to implement String()")
			}

			return tsast.SelectorExpr{
				Expr: ts.Ident(str.String()),
				Sel:  ts.Ident(dstr),
			}, nil
		}

		// It happens when we are overriding a Type parent with a default value. Because `hasOverrides` is true,
		// dstr is the default value, and we need to set the Type name here
		if str, ok := dvals[0].Source().(fmt.Stringer); ok {
			return ts.Ident(str.String()), nil
		}
	default:
		return nil, valError(v, "unknown selector subject type %T, cannot translate path %s", dvals[0].Source(), v.Path().String())
	}

	return nil, nil
}

func commentsFor(v cue.Value, jsdoc bool) []tsast.Comment {
	docs := v.Doc()
	if s, ok := v.Source().(*ast.Field); ok {
		for _, c := range s.Comments() {
			if !c.Doc && c.Line {
				docs = append(docs, c)
			}
		}
	}

	ret := make([]tsast.Comment, 0, len(docs))
	for _, cg := range docs {
		ret = append(ret, ts.CommentFromCUEGroup(ts.Comment{
			Text:      cg.Text(),
			Multiline: cg.Doc && !cg.Line,
			JSDoc:     jsdoc,
		}))
	}
	return ret
}

func commentsFromCueValue(v cue.Value) []string {
	docs := v.Doc()
	if s, ok := v.Source().(*ast.Field); ok {
		for _, c := range s.Comments() {
			if !c.Doc && c.Line {
				docs = append(docs, c)
			}
		}
	}

	ret := make([]string, 0, len(docs))
	for _, cg := range docs {
		for _, line := range strings.Split(strings.Trim(cg.Text(), "\n "), "\n") {
			ret = append(ret, line)
		}
	}
	return ret
}

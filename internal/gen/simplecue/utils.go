package simplecue

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/format"
)

func mustDumpsyn(v cue.Value) string {
	dump, err := dumpsyn(v)
	if err != nil {
		panic(err)
	}

	return dump
}

func dumpsyn(v cue.Value) (string, error) {
	syn := v.Syntax(
		cue.Concrete(false), // allow incomplete values
		cue.Definitions(false),
		cue.Optional(true),
		cue.Attributes(true),
		cue.Docs(true),
	)

	byt, err := format.Node(syn, format.TabIndent(true))
	return string(byt), err
}

func errorWithCueRef(v cue.Value, format string, args ...interface{}) error {
	return fmt.Errorf(v.Pos().String() + ": " + fmt.Sprintf(format, args...))
}

func selectorLabel(sel cue.Selector) string {
	if sel.Type().ConstraintType() == cue.PatternConstraint {
		return "*"
	}
	switch sel.LabelType() {
	case cue.StringLabel:
		return sel.Unquoted()
	case cue.DefinitionLabel:
		return sel.String()[1:]
	}
	// We shouldn't get anything other than non-hidden
	// fields and definitions because we've not asked the
	// Fields iterator for those or created them explicitly.
	panic(fmt.Sprintf("unreachable %v", sel.Type()))
}

// from https://github.com/cue-lang/cue/blob/99e8578ac45e5e7e6ebf25794303bc916744c0d3/encoding/openapi/build.go#L490
func appendSplit(a []cue.Value, splitBy cue.Op, v cue.Value) []cue.Value {
	op, args := v.Expr()
	// dedup elements.
	k := 1
outer:
	for i := 1; i < len(args); i++ {
		for j := 0; j < k; j++ {
			if args[i].Subsume(args[j], cue.Raw()) == nil &&
				args[j].Subsume(args[i], cue.Raw()) == nil {
				continue outer
			}
		}
		args[k] = args[i]
		k++
	}
	args = args[:k]

	if op == cue.NoOp && len(args) == 1 {
		// TODO: this is to deal with default value removal. This may change
		// when we completely separate default values from values.
		a = append(a, args...)
	} else if op != splitBy {
		a = append(a, v)
	} else {
		for _, v := range args {
			a = appendSplit(a, splitBy, v)
		}
	}
	return a
}

func getTypeHint(v cue.Value) (string, error) {
	// Direct lookup of attributes with Attribute() seems broken-ish, so do our
	// own search as best we can, allowing ValueAttrs, which include both field
	// and decl attributes.
	var found bool
	var attr cue.Attribute
	for _, a := range v.Attributes(cue.ValueAttr) {
		if a.Name() == annotationName {
			found = true
			attr = a
		}
	}

	if !found {
		return "", nil
	}

	tt, found, err := attr.Lookup(0, annotationKindFieldName)
	if err != nil {
		return "", err
	}

	if !found {
		return "", errorWithCueRef(v, "no value for the %q key in @%s attribute", annotationKindFieldName, annotationName)
	}
	return tt, nil
}

// ONLY call this function if it has been established that the provided Value is
// Concrete.
func cueConcreteToScalar(v cue.Value) (interface{}, error) {
	switch v.Kind() {
	case cue.NullKind:
		return nil, nil
	case cue.StringKind:
		return v.String()
	case cue.NumberKind, cue.FloatKind:
		return v.Float64()
	case cue.IntKind:
		return v.Int64()
	case cue.BoolKind:
		return v.Bool()
	default:
		return nil, errorWithCueRef(v, "can not convert kind to scalar: %s", v.Kind())
	}
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

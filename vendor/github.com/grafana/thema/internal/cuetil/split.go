package cuetil

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/token"
)

// AppendSplit recursively splits an expression in a single cue.Value by a
// single operation, flattening it into the slice of cue.Value that
// are joined by the provided operation in the input value.
//
// Most calls to this should pass nil for the third parameter.
func AppendSplit(v cue.Value, splitBy cue.Op, a []cue.Value) []cue.Value {
	if !v.Exists() {
		return a
	}
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
			a = AppendSplit(v, splitBy, a)
		}
	}
	return a
}

// PrintPosList dumps the cue.Value.Pos for each unified element of the provided
// cue.Value.
//
// Useful for debugging values with complex multiple unified antecedents.
func PrintPosList(v cue.Value) {
	for i, dval := range AppendSplit(v, cue.AndOp, nil) {
		fmt.Println(i, dval.Pos())
	}
}

// PoslistWithoutThema returns all token.Pos associated with a given cue.Value, omitting
// any token.Pos that point to thema.
func PoslistWithoutThema(v cue.Value) []token.Pos {
	vals := AppendSplit(v, cue.AndOp, nil)
	poslist := make([]token.Pos, 0, len(vals))
	for _, dval := range vals {
		// TODO not sure if we should expect os-sensitive path separators here or not
		if pos := dval.Pos(); pos != token.NoPos && !strings.Contains(pos.Filename(), "github.com/grafana/thema") {
			poslist = append(poslist, pos)
		}
	}
	return poslist
}

// FirstNonThemaPos returns the first [token.Pos] in the slice returned by
// PoslistWithoutThema, or [token.NoPos] if no such pos exists.
func FirstNonThemaPos(v cue.Value) token.Pos {
	pl := PoslistWithoutThema(v)
	if len(pl) == 0 {
		return token.NoPos
	}
	return pl[0]
}

//
// func RemoveThemaValues(v cue.Value) []cue.Value {
// 	if vlist := AppendSplit(v, cue.AndOp, nil); len(vlist) > 1 {
// 		others := make([]cue.Value, 0, len(vlist))
// 		for _, av := range vlist {
// 			_, path := av.ReferencePath()
// 			if path.String() != "#Lineage" {
// 				others = append(others, av)
// 			}
// 		}
// 	}
// }

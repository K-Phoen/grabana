package thema

import (
	"bytes"
	"errors"

	"cuelang.org/go/cue"
	cuejson "cuelang.org/go/pkg/encoding/json"
)

// TODO clean up signature to only return cue.Value
func doHydrate(sch, data cue.Value) (cue.Value, error) {
	switch sch.IncompleteKind() {
	case cue.ListKind:
		// if list element exist
		ele := sch.LookupPath(cue.MakePath(cue.AnyIndex))

		// if data is not a concrete list, we must have list elements exist to be used to trim defaults
		if ele.Exists() {
			if ele.IncompleteKind() == cue.BottomKind {
				panic("unreachable")
			}
			iter, err := data.List()
			if err != nil {
				panic("unreachable - kind was list")
			}
			var iterlist []cue.Value
			for iter.Next() {
				ref, err := getBranch(ele, iter.Value())
				if err != nil {
					return data, err
				}
				re, err := doHydrate(iter.Value(), ref)
				if err == nil {
					iterlist = append(iterlist, re)
				}
			}
			liInstance := sch.Context().NewList(iterlist...)
			if liInstance.Err() != nil {
				return data, liInstance.Err()
			}
			return liInstance, nil
		}
		return data.Unify(sch), nil
	case cue.StructKind:
		iter, err := sch.Fields(cue.Optional(true))
		if err != nil {
			return data, err
		}
		for iter.Next() {
			label, _ := iter.Value().Label()
			lv := data.LookupPath(cue.MakePath(cue.Str(label)))
			if err != nil {
				continue
			}
			if lv.Exists() {
				res, err := doHydrate(lv, iter.Value())
				if err != nil {
					continue
				}
				data = data.FillPath(cue.MakePath(cue.Str(label)), res)
			} else if !iter.IsOptional() {
				data = data.FillPath(cue.MakePath(cue.Str(label)), iter.Value().Eval())
			}
		}
		return data, nil
	default:
		data = data.Unify(sch)
	}
	return data, nil
}

func convertCUEValueToString(inputCUE cue.Value) (string, error) {
	re, err := cuejson.Marshal(inputCUE)
	if err != nil {
		return re, err
	}

	result := []byte(re)
	result = bytes.Replace(result, []byte("\\u003c"), []byte("<"), -1)
	result = bytes.Replace(result, []byte("\\u003e"), []byte(">"), -1)
	result = bytes.Replace(result, []byte("\\u0026"), []byte("&"), -1)
	return string(result), nil
}

func getDefault(icue cue.Value) (cue.Value, bool) {
	d, exist := icue.Default()
	if exist && d.Kind() == cue.ListKind {
		len, err := d.Len().Int64()
		if err != nil {
			return d, false
		}
		var defaultExist bool
		if len <= 0 {
			op, vals := icue.Expr()
			if op == cue.OrOp {
				for _, val := range vals {
					vallen, _ := val.Len().Int64()
					if val.Kind() == cue.ListKind && vallen <= 0 {
						defaultExist = true
						break
					}
				}
				if !defaultExist {
					exist = false
				}
			} else {
				exist = false
			}
		}
	}
	return d, exist
}

func isCueValueEqual(inputdef cue.Value, input cue.Value) bool {
	d, exist := getDefault(inputdef)
	if exist {
		return input.Subsume(d) == nil && d.Subsume(input) == nil
	}
	return false
}

// func ddoDehydrate(sch cue.Value, data cue.Value) cue.Value {
// 	rv := sch.Runtime().CompileString("_", cue.Filename("dehydrated"))
//
// 	switch sch.IncompleteKind() {
// 	case cue.StructKind:
// 		for iter, _ := data.Fields(cue.Optional(true), cue.Definitions(true)); iter.Next(); {
// 			datav, p := iter.Value(), cue.MakePath(iter.Selector())
// 			schv := sch.LookupPath(p)
// 			_, has := schv.Default()
// 			// No schema [default] value means we need the data
// 			if !has || !schv.Exists() {
// 				rv = rv.FillPath(p, datav)
// 			} else if dehyd := ddoDehydrate(schv, datav); dehyd.Exists() {
// 				rv = rv.FillPath(p, datav)
// 			}
// 		}
// 		return rv
// 	case cue.ListKind:
// 		schdef, has := sch.Default()
// 		if !has {
// 			return rv.FillPath(cue.Path{}, data)
// 		}
// 		// Just don't try if the schema default is a complex structure.
// 		// Maybe we can improve this later.
// 		if op, _ := schdef.Expr(); op != cue.NoOp {
// 			return rv.FillPath(cue.Path{}, data)
// 		}
// 		// TODO seems like we need more here?
// 		if !eq(schdef, data) {
// 			return rv.FillPath(cue.Path{}, data)
// 		}
// 	default:
// 		schdef, has := sch.Default()
// 		if !has || !eq(schdef, data) {
// 			return rv.FillPath(cue.Path{}, data)
// 		}
// 	}
//
// 	return rv
// }

// func eq(a, b cue.Value) bool {
// 	return (a.Subsume(b) == nil) && (b.Subsume(a) == nil)
// }

// TODO clean up signature to only return cue.Value
func doDehydrate(sch, data cue.Value) (cue.Value, bool, error) {
	// To include all optional fields, we need to use sch for iteration,
	// since the lookuppath with optional field doesn't work very well
	rv := sch.Context().CompileString("", cue.Filename("helper"))
	if rv.Err() != nil {
		return data, false, rv.Err()
	}

	switch sch.IncompleteKind() {
	case cue.StructKind:
		// Get all fields including optional fields
		iter, _ := sch.Fields(cue.Optional(true))
		keySet := make(map[string]bool)
		for iter.Next() {
			lable, _ := iter.Value().Label()
			keySet[lable] = true
			lv := data.LookupPath(cue.MakePath(cue.Str(lable)))
			if lv.Exists() {
				re, isEqual, err := doDehydrate(iter.Value(), lv)
				if err == nil && !isEqual {
					rv = rv.FillPath(cue.MakePath(cue.Str(lable)), re)
				}
			}
		}
		// Get all the fields that are not defined in schema
		iter, _ = data.Fields()
		for iter.Next() {
			lable, _ := iter.Value().Label()
			if exists := keySet[lable]; !exists {
				rv = rv.FillPath(cue.MakePath(cue.Str(lable)), iter.Value())
			}
		}
		return rv, false, nil
	case cue.ListKind:
		if isCueValueEqual(sch, data) {
			return rv, true, nil
		}

		// take every element of the list
		ele := sch.LookupPath(cue.MakePath(cue.AnyIndex))

		// if data is not a concrete list, we must have list elements exist to be used to trim defaults
		if ele.Exists() {
			if ele.IncompleteKind() == cue.BottomKind {
				return rv, true, nil
			}

			iter, err := data.List()
			if err != nil {
				return rv, true, nil
			}
			var iterlist []cue.Value
			for iter.Next() {
				ref, err := getBranch(ele, iter.Value())
				if err != nil {
					iterlist = append(iterlist, iter.Value())
					continue
				}
				re, isEqual, err := doDehydrate(ref, iter.Value())
				if err == nil && !isEqual {
					iterlist = append(iterlist, re)
				} else {
					iterlist = append(iterlist, iter.Value())
				}
			}
			liInstance := sch.Context().NewList(iterlist...)
			return liInstance, false, liInstance.Err()
		}
		// now when ele is empty, we don't trim anything
		return data, false, nil

	default:
		if isCueValueEqual(sch, data) {
			return data, true, nil
		}
		return data, false, nil
	}
}

func getBranch(sch cue.Value, data cue.Value) (cue.Value, error) {
	op, defs := sch.Expr()
	if op == cue.OrOp {
		for _, def := range defs {
			err := def.Unify(data).Validate(cue.Concrete(true))
			if err == nil {
				return def, nil
			}
		}
		// no matching branches? wtf
		return sch, errors.New("no branch is found for list")
	}
	return sch, nil
}

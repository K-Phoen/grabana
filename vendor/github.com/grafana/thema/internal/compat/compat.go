package compat

import (
	"cuelang.org/go/cue"
)

// ThemaCompatible is the canonical Thema algorithm for checking that the
// [cue.Value] s is (backwards) compatible with p. A nil return indicates
// compatibility.
//
// The behavior of this function is undefined if s and p are not closed
// structs. TODO check this and error if conditions aren't met
func ThemaCompatible(p, s cue.Value) error {
	return s.Subsume(p, cue.Raw(), cue.All())
}

// type CompatInvariantError struct {
// 	rawlin    cue.Value
// 	violation [2]SyntacticVersion
// 	detail    error
// }
//
// func (e *CompatInvariantError) Error() string {
// 	if e.violation[0][0] == e.violation[1][0] {
// 		// TODO better
// 		return e.detail.Error()
// 	}
// 	return fmt.Sprintf("schema %s must be backwards incompatible with schema %s", e.violation[1], e.violation[0])
// }

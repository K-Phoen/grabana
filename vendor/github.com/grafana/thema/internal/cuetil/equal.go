package cuetil

import (
	"cuelang.org/go/cue"
)

// Equal reports nil when the two cue values subsume each other or an error otherwise
func Equal(val1 cue.Value, val2 cue.Value) error {
	if err := val1.Subsume(val2, cue.Raw()); err != nil {
		return err
	}

	if err := val2.Subsume(val1, cue.Raw()); err != nil {
		return err
	}

	return nil
}

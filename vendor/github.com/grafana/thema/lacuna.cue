package thema

// A Lacuna represents a semantic gap in a Lens's mapping between schemas.
//
// For any given mapping between schema, there may exist some valid values and
// intended semantics on either side that are impossible to precisely translate.
// When such gaps occur, and an actual schema instance falls into such a gap,
// the Lens is expected to emit Lacuna that describe the general nature of the
// translation gap.
//
// A lacuna may be unconditional (the gap exists for all possible instances
// being translated between the schema pair) or conditional (the gap only exists
// when certain values appear in the instance being translated between schema).
// When translating, the predicate in each lacuna's "condition" field is checked
// to determine whether the lacuna is applicable to the particular instance
// being translated.
#Lacuna: {
	// A reference to a field and its value in a schema/instance.
	#FieldRef: {
		// TODO would be great to be able to constrain that value should always be a reference,
		// and path is (a modified version of) the string representation of the reference
		value: _
		path:  string
	}

	// Indicates whether the lacuna applies to the particular object being translated.
	// If this field evaluates to false, this lacuna is not emitted for a particular
	// translation. If unspecified, defaults to true.
	condition: bool | *true

	// The field path(s) and their value(s) in the pre-translation instance
	// that are relevant to the lacuna.
	sourceFields: [...#FieldRef]

	// The field path(s) and their value(s) in the post-translation instance
	// that are relevant to the lacuna.
	targetFields: [...#FieldRef]

	// At least one of sourceFields or targetFields must be non-empty.
	// TODO(must) https://github.com/cue-lang/cue/issues/943
	// must(len(sourceFields) > 0 || len(targetFields) > 0, "at least one of sourceFields or targetFields must be non-empty")
	// _mustlen: >0 & (len(sourceFields) + len(targetFields))

	// A human-readable message describing the gap in translation.
	message: string

	type: or([ for t in #LacunaTypes {t}])
}

#LacunaTypes: [N=string]: #LacunaType & {name: N}
#LacunaTypes: {
	// Placeholder lacunas indicate that a field in the target instance has
	// been filled with a placeholder value.
	//
	// Use Placeholder when introducing a new required field that lacks a default,
	// and it is necessary to fill the field with some value to meet lens
	// validity requirements.
	//
	// A placeholder is NOT a schema-defined default. It is expressly the
	// opposite: a lens-defined value that exists solely to be replaced by the
	// calling program.
	Placeholder: {
		id: 1
	}

	// DroppedField lacunas indicate that field(s) in the source instance were
	// dropped in a manner that potentially lost some of their contained semantics.
	//
	// When a lens drops multiple fields, prefer to create one DroppedField
	// lacuna per distinct cause. For example, if multiple instance fields are
	// dropped from a single open struct because they were absent from the
	// schema, all of those fields should be included in a single DroppedField.
	DroppedField: {
		id: 2
	}

	// LossyFieldMapping lacunas indicate that no clear mapping existed from the
	// source field value to the intended semantics of any valid target field
	// value. 
	// 
	// Only use this lacuna type when there exists at least one valid source
	// value with a clear, lossless mapping to the target value.
	LossyFieldMapping: {
		id: 3
	}

	// ChangedDefault lacunas indicate that the source field value was the
	// schema-specified default, and the default changed in the target field,
	// and the value in the instance was changed as well.
	//
	// NOTE the semantics of field presence/absence in the instance are subtle here, and this may need refinement
	ChangedDefault: {
		id: 4
	}
}

#LacunaType: {
	name: string
	id:   int // FIXME this is a dumb way of trying to express identity
}

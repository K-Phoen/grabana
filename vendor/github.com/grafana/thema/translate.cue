package thema

import "list"

// TODO docs
#TranslatedInstance: {
	from:   #SyntacticVersion
	to:     #SyntacticVersion
	result: _
	lacunas: [...#Lacuna]
}

// Translate takes an instance, a lineage, and to and from versions.
//
// The provided instance must be a valid instance of the 'from' schema.
//
// The instance is iteratively transformed through the lineage's lenses,
// starting from the version the instance is valid against, and
// continuing until the target schema version is reached.
//
// The result is the instance in final translated form, the schema versions
// at which the translation started and ended, and any lacunas emitted during
// translation.
//
// TODO functionize
#Translate: {
	L=lin: _
	I=inst: {...}
	FV=from: #SyntacticVersion
	TV=to:   #SyntacticVersion

	let cmp = (_cmpSV & {l: FV, r: TV}).out

	// TODO validate I is instance of FV schema

	out: {
		steps: [...#TranslatedInstance]
		result: {
			from:   FV
			to:     TV
			result: *steps[len(steps)-1].result | I
		}
	}
	out: {
		steps: [
			// to version older/smaller than from version
			if cmp == 1 {
				// schrange is the subset of schemas being traversed in this
				// translation, inclusive of the starting schema.
				let hi = (L._flatidx & {v: FV}).out
				let lo = (L._flatidx & {v: TV}).out

				_accum: [{result: I, to: FV}, for i, pos in list.Range(hi-1, lo-1, -1) {
					// alias pointing to the previous item in the list we're building
					let prior = _accum[i]

					// the actual schema def
					let schdef = L.schemas[pos]

					// TODO does having this field in the result, even hidden, cause a problem? does using an alias cause the computation to run more than once?
					_lens: L._backwardLenses[(L._flatidx & {v: schdef.version}).out] & {
						// pass the prior result along as the input to this lens
						input: prior.result

						// guard against bugs in this impl - these are already concretely specified
						// on the lens itself by the user, so these version numbers must be correct
						from: prior.to
						to:   schdef.version
					}

					from: _lens.from
					to:   _lens.to
					// TODO initial input isn't necessarily unified with schema - does that make translated output meaningfully different?
					result: {_lens.result, schdef._#schema}
					lacunas: [ for lac in _lens.lacunas if lac.condition {lac.lacuna}]
				}]

				// Final value excludes first element (initial input) from accum
				list.Drop(_accum, 1)
			},
			// to version newer/larger than from version
			if cmp == -1 {
				// schrange contains the subset of schemas being traversed in this
				// translation, inclusive of the starting schema.
				let lo = (L._flatidx & {v: FV}).out
				let hi = (L._flatidx & {v: TV}).out
				let schrange = list.Slice(L.schemas, lo+1, hi+1)

				_accum: [{result: I, to: FV}, for i, schdef in schrange {
					// alias pointing to the previous item in the list we're building
					let prior = _accum[i]

					if prior.to[0] == schdef.version[0] {
						// Only a minor version. Backwards compatibility rules dictate that the mapping
						// algorithm for the forward lens is generic: simple unification.
						from: prior.to
						to:   schdef.version
						//						result: prior.result & schdef._#schema
						result: {prior.result, schdef._#schema}
					}

					if prior.to[0] < schdef.version[0] {
						// TODO does having this field in the result, even hidden, cause a problem? does using an alias cause the computation to run more than once?
						_lens: L._forwardLenses[schdef.version[0]-1] & {
							// pass the prior result along as the input to this next lens
							input: prior.result

							// guard against bugs in this impl - these are already concretely specified
							// on the lens itself by the user, so these version numbers must be correct
							from: prior.to
							to:   schdef.version
						}

						from: prior.to
						to:   schdef.version
						//						result: {_lens.result, schdef._#schema, {lidx: lensidx}}
						result: {_lens.result, schdef._#schema}
						lacunas: [ for lac in _lens.lacunas if lac.condition {lac.lacuna}]
						// Crossing a major version. The forward lens explicitly defined in the schema
						// provides the mapping algorithm.
					}
				}]

				// Final value excludes first element (initial input) from accum
				list.Drop(_accum, 1)
			},
			// to version same as from version is a no-op
			if cmp == 0 {[]},
		][0]
	}
}

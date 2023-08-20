package thema

import (
	"struct"
	"list"
)

// Lineage is the top-level container in thema, holding the complete
// evolutionary history of a particular kind of object: every schema that has
// ever existed for that object, and the lenses that allow translating between
// those schema versions.
#Lineage: {
	// The name of the thing specified by the schemas in this lineage.
	//
	// A lineage's name must not change as it evolves.
	name: string
	// TODO(must) https://github.com/cue-lang/cue/issues/943
	// name: must(isconcrete(name), "all lineages must have a name")

	// joinSchema governs the shape of schema that may be expressed in a
	// lineage. It is the least upper bound, or join, of the acceptable schema
	// value space; the schemas defined in this lineage must be instances of the
	// joinSchema.
	//
	// All Thema schemas must be struct-kinded. Consequently, if a lineage defines
	// a joinSchema, it must be a struct containing at least one field.
	//
	// A lineage's joinSchema must never change as the lineage evolves.

	//	joinSchema?: struct.MinFields(1)
	joinSchema: _

	// schemas is the ordered list of all schemas in the lineage.
	//
	// Each element is a #SchemaDef, injected with the joinSchema for this lineage.
	//
	// It is recommended but not required that schema entries in this list be defined
	// in ascending order by version. Thema tooling that modifies and emits lineages
	// definitions may produce schemas sorted in ascending order, rather than original
	// source order.
	// TODO switch to descending order - newest on top is nicer to read
	schemas: [#SchemaDef, ...#SchemaDef]
	//			schemas: [...#SchemaDef]

	if joinSchema != _|_ {
		schemas: [{_join: joinSchema}, ...{{_join: joinSchema}}]
	}

	// lenses contains all the mappings between all the schemas in the lineage.
	//
	// For a lineage to be valid, it must contain lenses such that an instance of
	// each schema can be translated into both its predecessor and successor schemas.
	//
	// Because minor version changes are backwards compatible by definition, a lens
	// implicitly exists, and no lens definition may be defined by the lineage author.
	// However, all other version transitions require an explicit lens definition:
	//
	//  - A lens mapping forward across every breaking change/new major version
	//  - A lens mapping backward across every change
	//
	// Thus, for a lineage with schema versions [0,0], [0,1], [1,0], [2,0], [2,1],
	// the following lenses must exist (implicit lenses are wrapped in parentheses):
	//
	//  [0,1] -> [0,0]
	// ([0,0] -> [0,1])
	//  [1,0] -> [0,1]
	//  [0,1] -> [1,0]
	//  [2,0] -> [1,0]
	//  [1,0] -> [2,0]
	//  [2,1] -> [2,0]
	// ([2,0] -> [2,1])
	//
	// To be valid, a lineage must define the exact set of explicit lenses entailed by its
	// set of schema versions. It is not permitted to explicitly define a lens across
	// non-breaking changes.
	//
	// The above ordering of lenses, sorted ascending first by 'to' then by 'from' version,
	// is the recommended but not required order in which lenses should be defined in this
	// list. Thema tooling that modifies and emits lineages definitions may produce lenses
	// sorted in ascending order, rather than original source order.
	// TODO switch to descending order - newest on top is nicer to read
	lenses: [...#Lens]

	_atLeastOneSchema: len(schemas) > 0

	SS=_schemas: [...]
	if _atLeastOneSchema == true {
		_schemas: schemas
	}

	if _atLeastOneSchema == false {
		_schemas: [#SchemaDef & {version: [0, 0]}]
	}

	_forwardLenses: [ for lens in lenses if {lens.to[0] > lens.from[0]} {lens}]
	_backwardLenses: [ for lens in lenses if {lens.to[0] <= lens.from[0]} {lens}]

	// preserved for debugging
	//	lensVersions: {
	//		backward: {
	//			both: [ for lens in _backwardLenses {"\(lens.from[0]).\(lens.from[1])->v\(lens.to[0]).\(lens.to[1])"}]
	//			to: [ for lens in _backwardLenses {"v\(lens.to[0]).\(lens.to[1])"}]
	//			from: [ for lens in _backwardLenses {"v\(lens.from[0]).\(lens.from[1])"}]
	//		}
	//		forward: {
	//			both: [ for lens in _forwardLenses {"\(lens.from[0]).\(lens.from[1])->v\(lens.to[0]).\(lens.to[1])"}]
	//			to: [ for lens in _forwardLenses {"v\(lens.to[0]).\(lens.to[1])"}]
	//			from: [ for lens in _forwardLenses {"v\(lens.from[0]).\(lens.from[1])"}]
	//		}
	//	}

	// _counts tracks the number of versions in each major version in the lineage.
	// The index corresponds to the major version number, and the value is the
	// number of minor versions within that major.
	_counts: [...uint64] & list.MinItems(1)

	//	counts: _counts
	// TODO check subsumption (backwards compat) of each schema with its successor natively in CUE
	if len(SS) == 1 {
		_counts: [0]
	}

	if len(SS) > 1 {
		_pos: [0, for i, sch in list.Drop(SS, 1) if SS[i].version[0] < sch.version[0] {i + 1}]
		_counts: [ for i, idx in list.Slice(_pos, 0, len(_pos)-1) {
			_pos[i+1] - list.Sum(list.Slice(_pos, 0, i+1))
		}, len(SS) - _pos[len(_pos)-1]]

		// The following approach to the above:
		//
		//		let pos = [0, for i, sch in SS[1:] if SS[i].version[0] < sch.version[0] { i+1 }]
		//		_counts: [for i, idx in pos[:len(pos)-1] {
		//			pos[i+1]-list.Sum(pos[:i+1])
		//		}, len(SS)-pos[len(pos)-1]]
		//
		// causes the following cue internals panic:
		// panic: getNodeContext: nodeContext out of sync [recovered]
		//	panic: getNodeContext: nodeContext out of sync
	}

	// _basis keeps the index of the first schema in each major version
	// within the overall canonical schema sort ordering. This allows trivial
	// schema retrieval from a syntactic version.
	_basis: [0, for maj, _ in list.Drop(_counts, 1) {
		list.Sum(list.Take(_counts, maj+1))
	}]

	// Pick is a pseudofunction that returns the schema from this Lineage
	// (lin) that corresponds to the provided SyntacticVersion (v). Bounds
	// constraints enforce that the provided version number exists within the
	// lineage.
	//
	// Pick is the only correct mechanism to retrieve a lineage's declared schema.
	// Retrieving a lineage's schemas by direct indexing will not check invariants,
	// apply compositions or joinSchemas.
	//#Pick: {
	//	// The schema version to retrieve.
	//	v: #SyntacticVersion & [<len(L._counts), <=L._counts[v[0]]]
	//	// TODO(must) https://github.com/cue-lang/cue/issues/943
	//	// must(isconcrete(v[0]), "must specify a concrete major version")
//
	//	out: SS[_basis[v[0]]+v[1]]._#schema
	//}

	// PickDef takes the same arguments as Pick, but returns the entire
	// #SchemaDef rather than only the schema body itself.
	//#PickDef: {
	//	// The schema version to retrieve.
	//	v: #SyntacticVersion & [<len(L._counts), <=L._counts[v[0]]]
	//	// TODO(must) https://github.com/cue-lang/cue/issues/943
	//	// must(isconcrete(v[0]), "must specify a concrete sequence number")
//
	//	out: SS[_basis[v[0]]+v[1]]
	//}

	// latestVersion always contains the SyntacticVersion of the lineage's
	// "latest" schema - the schema with the largest version number.
	//
	// Take care in using this. If any code that depends on schema contents relies
	// on it, that code will break as soon as a breaking schema change is made. This
	// may be desirable within a tight development loop - e.g., for a finite team,
	// working within a single VCS repository - in order to force updating code that
	// must be kept in sync.
	//
	// But relying on it in, for example, an API client based on Thema lineages
	// undermines the entire goal of Thema, as it would forces breaking changes
	// immediately on the client's users, rather than allowing them to update at
	// their own pace.
	//
	// TODO functionize
	#LatestVersion: SS[len(SS)-1].version

	_flatidx: {
		v: #SyntacticVersion
		// TODO check what happens when out of bounds
		out: _basis[v[0]] + v[1]
	}
}

// #SchemaDef represents a single schema declaration in Thema. In addition to
// the schema itself, it contains the schema's version, optional examples,
// composition instructions, and lenses that map to or from the schema, as
// required by Thema's invariants.
//
// Note that the version number must be explicitly declared, even though the
// correct value is algorithmically determined.
#SchemaDef: {
	// version is the Syntactic Version number of the schema. While this property
	// is settable by lineage authors, it has exactly one correct value for any
	// particular #SchemaDef in any lineage, algorithmically determined by its
	// position in the list of schemas and the number of its predecessors that
	// make breaking changes to their schemas.
	//
	// Despite there being only one correct choice, lineage authors must still
	// explicitly declare the schema version. Future improvements in Thema may make
	// this unnecessary, but explicitly declaring the version is always useful for
	// readability.
	//
	// The entire lineage is considered invalid if the version number in this field
	// is inconsistent with the algorithmically determined set of [non-]breaking changes.
	version: #SyntacticVersion | *[0, 0]

	schema: _

	_join: _

	// Thema's internal handle for the user-provided schema definition. This
	// handle is used by all helpers/operations in the thema package. As a
	// CUE definition, use of this handle entails that all thema schemas are
	// always recursively closed by default.
	//
	// This handle is also unified with the joinSchema of the containing lineage.

	_#schema: _join & schema

	//	_schemaIsNonEmpty: struct.MinFields(1) & _#schema

	// examples is an optional set of named examples of the schema, intended
	// for use in documentation or other non-functional contexts.
	examples?: [string]: _#schema
}

// Lens defines a transformation that maps the fields of one schema in a lineage to the
// fields of another schema, as well as the lacunas that may exist for specific objects
// when translating instances between these schemas.
#Lens: {
	// The schema version that is the input or source of the mapping.
	from: #SyntacticVersion

	// The schema version that is the result or target of the mapping.
	to: #SyntacticVersion

	// input is filled with an object instance of the 'from' schema in order to use the
	// lens to translate between schema versions.
	input: _

	// The relation between the schemas identified by the 'from' and 'to' versions,
	// expressed as a mapping from 'input' to this field.
	//
	// The value must be an instance of the 'to' schema, constructed through
	// references to the 'from' schema.
	//
	// For example, if the schemas corresponding to the 'from' and 'to' versions are:
	//
	//   from: { a: string }
	//   to:   { b: string }
	//
	// and the goal is to remap the field 'a' to be called 'b', result should be written as:
	//
	//   result: { b: L.input.a }
	result: struct.MinFields(0)

	// lacunas describe semantic gaps in the transform's mapping. See lacuna docs
	// for more information (TODO).
	lacunas: [...#Lacuna]
}

// SyntacticVersion is an ordered pair of non-negative integers. It represents
// the version of a schema within a lineage, or the version of an instance that
// is valid with respect to the schema of that version.
#SyntacticVersion: [uint64, uint64]

// TODO functionize
_cmpSV: FN={
	l:   #SyntacticVersion
	r:   #SyntacticVersion
	out: -1 | 0 | 1
	out: {
		if FN.l[0] < FN.r[0] {-1}
		if FN.l[0] > FN.r[0] {1}
		if FN.l[0] == FN.r[0] && FN.l[1] < FN.r[1] {-1}
		if FN.l[0] == FN.r[0] && FN.l[1] > FN.r[1] {1}
		if FN.l == FN.r {0}
	}
}

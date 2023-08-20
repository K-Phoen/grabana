package thema

// SearchAndValidate is a pseudofunction that takes a lineage (lin) and some candidate
// data (inst) as an argument, and searches the lineage for a schema against which
// that data is valid.
//
// A #LinkedInstance is "returned" (out) corresponding to the first schema in the
// lineage for which the data is a valid instance of the schema. Data is not checked
// for concreteness. If there is no match, bottom is returned.
//
// TODO optionally check for concreteness of otherwise valid data
// TODO functionize
#SearchAndValidate: fn={
	lin: _
	inst: {...} // TODO consistently rename to 'object' or something

	out: #LinkedInstance
	out: [ for _, sch in fn.lin.schemas if ((inst & sch._#schema) != _|_) {
		v:    sch.version
		lin:  fn.lin
		inst: inst & sch._#schema
	}, _|_][0]
}

#ValidFor: {
	lin: _
	inst: {...} // TODO consistently rename to 'object' or something

	out: #SyntacticVersion
	out: [ for _, sch in lin.schemas if ((sch._#schema & inst) != _|_) {sch.version}][0]
}

// #LinkedInstance represents data that is an instance of some schema, the
// version of that schema, and the lineage of the schema.
#LinkedInstance: {
	inst: {...} // TODO consistently rename to 'object' or something
	lin:        #Lineage & {_atLeastOneSchema: true}
	v:          #SyntacticVersion              // TODO rename to 'version'

	// TODO need proper validation/subsumption check here, not simple unification
	//	_valid: inst & (#Pick & {lin: L.lin, v: v}).out
}

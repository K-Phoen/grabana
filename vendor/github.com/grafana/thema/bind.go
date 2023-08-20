package thema

import (
	"bytes"
	"fmt"

	"cuelang.org/go/cue"
	cerrors "cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/token"
	"github.com/cockroachdb/errors"

	terrors "github.com/grafana/thema/errors"
	"github.com/grafana/thema/internal/compat"
)

// maybeLineage is an intermediate processing structure used to validate
// inputs as actual lineages
//
// it's important that these flags are populated in order to avoid false negatives.
// no system ensures this, it's all human reasoning
type maybeLineage struct {
	// user lineage definition, NOT unified with thema.#Lineage
	raw cue.Value

	// user lineage definition, unified with thema.#Lineage
	uni cue.Value

	// original input cue.Value representing the lineage. May or may not be unified
	// with thema.#Lineage
	orig cue.Value

	rt *Runtime

	// pos of the original input for the lineage
	pos token.Pos

	// bind options passed by the caller
	cfg *bindConfig

	schlist []*schemaDef

	allv []SyntacticVersion

	implens []ImperativeLens

	lensmap map[lensID]ImperativeLens

	// The raw input value is the root of a package instance
	// rawIsPackage bool
}

// to, from
type lensID struct {
	From, To SyntacticVersion
}

func lid(from, to SyntacticVersion) lensID {
	return lensID{from, to}
}

func (id lensID) String() string {
	return fmt.Sprintf("%s -> %s", id.From, id.To)
}

func (ml *maybeLineage) checkGoValidity(cfg *bindConfig) error {
	schiter, err := ml.uni.LookupPath(cue.MakePath(cue.Str("schemas"))).List()
	if err != nil {
		panic(fmt.Sprintf("unreachable - should have already verified schemas field exists and is list: %+v", cerrors.Details(err, nil)))
	}
	vpath := cue.MakePath(cue.Str("version"))

	var previous *schemaDef
	for schiter.Next() {
		// Only thing not natively enforced in CUE is that the #SchemaDef.version field is concrete
		svval := schiter.Value().LookupPath(vpath)
		iter, err := svval.List()
		if err != nil {
			panic(fmt.Sprintf("unreachable - should have already verified #SchemaDef.version field exists and is list: %+v", err))
		}
		for iter.Next() {
			if !iter.Value().IsConcrete() {
				return errors.Mark(mkerror(iter.Value(), "#SchemaDef.version must have concrete major and minor versions"), terrors.ErrInvalidLineage)
			}
		}
		sch := &schemaDef{}
		err = svval.Decode(&sch.v)
		if err != nil {
			panic(fmt.Sprintf("unreachable - could not decode syntactic version: %+v", err))
		}

		if err := ml.checkSchemasOrder(previous, sch); err != nil {
			return err
		}

		sch.ref = schiter.Value()
		sch.def = sch.ref.LookupPath(pathSchDef)
		if previous != nil && !cfg.skipbuggychecks {
			compaterr := compat.ThemaCompatible(previous.def, sch.def)
			if sch.v[1] == 0 && compaterr == nil {
				// Major version change, should be backwards incompatible
				return errors.Mark(mkerror(sch.ref.LookupPath(pathSch), "schema %s must be backwards incompatible with schema %s: introduce a breaking change, or redeclare as version %s", sch.v, previous.v, synv(previous.v[0], previous.v[1]+1)), terrors.ErrInvalidLineage)
			}
			if sch.v[1] != 0 && compaterr != nil {
				// Minor version change, should be backwards compatible
				return errors.Mark(mkerror(sch.ref.LookupPath(pathSch), "schema %s is not backwards compatible with schema %s:\n%s", sch.v, previous.v, cerrors.Details(compaterr, nil)), terrors.ErrInvalidLineage)
			}
		}

		ml.schlist = append(ml.schlist, sch)
		ml.allv = append(ml.allv, sch.v)
		previous = sch
	}

	return nil
}

func (ml *maybeLineage) checkSchemasOrder(prev, curr *schemaDef) error {
	if prev == nil {
		return nil
	}

	if curr.v.Less(prev.v) {
		return errors.Mark(mkerror(curr.ref.LookupPath(pathSch), "schema version %s is not greater than previous schema version %s", curr.v, prev.v), terrors.ErrInvalidSchemasOrder)
	}

	return nil
}

func (ml *maybeLineage) checkExists(cfg *bindConfig) error {
	p := ml.raw.Path().String()
	// The candidate lineage must exist.
	// TODO can we do any better with contextualizing these errors?
	if !ml.raw.Exists() {
		if p != "" {
			return errors.Mark(errors.Newf("not a lineage: no cue value at path %q", p), terrors.ErrValueNotExist)
		}

		return errors.WithStack(terrors.ErrValueNotExist)
	}
	return nil
}

func (ml *maybeLineage) checkLineageShape(cfg *bindConfig) error {
	// Check certain paths specifically, because these are common getting started errors of just arranging
	// CUE statements in the right way that deserve more targeted guidance
	for _, path := range []string{"name", "schemas"} {
		val := ml.raw.LookupPath(cue.MakePath(cue.Str(path)))
		if !val.Exists() {
			return errors.Mark(mkerror(ml.raw, "not a lineage, missing #Lineage.%s", path), terrors.ErrValueNotALineage)
		}
		if !val.IsConcrete() {
			return errors.Mark(mkerror(val, "invalid lineage, #Lineage.%s must be concrete", path), terrors.ErrInvalidLineage)
		}
	}

	// The candidate lineage must be an instance of #Lineage. However, we can't validate the whole
	// structure, because lenses will fail validation. This is because we currently expect them to be written:
	//
	// {
	// 		input: _
	// 		result: {
	// 			foo: input.foo
	// 		}
	// }
	//
	// means that those structures won't pass Validate until we've injected an actual object there.
	if err := ml.uni.Validate(cue.Final()); err != nil {
		return errors.Mark(cerrors.Promote(err, "not an instance of thema.#Lineage"), terrors.ErrInvalidLineage)
	}

	return nil
}

// Checks the validity properties of lineages that are expressible natively in CUE.
func (ml *maybeLineage) checkNativeValidity(cfg *bindConfig) error {
	// The candidate lineage must be error-free.
	// TODO replace this with Err, this check isn't actually what we want up here. Only schemas themselves must be cycle-free
	if err := ml.raw.Validate(cue.Concrete(false)); err != nil {
		return errors.Mark(cerrors.Promote(err, "lineage is invalid"), terrors.ErrInvalidLineage)
	}
	if err := ml.uni.Validate(cue.Concrete(false)); err != nil {
		return errors.Mark(cerrors.Promote(err, "lineage is invalid"), terrors.ErrInvalidLineage)
	}

	return nil
}

func (ml *maybeLineage) checkLensesOrder() error {
	// Two distinct validation paths, depending on whether the lenses were defined in
	// Go or CUE.
	if len(ml.implens) > 0 {
		return ml.checkGoLensCompleteness()
	}

	lensIter, err := ml.uni.LookupPath(cue.MakePath(cue.Str("lenses"))).List()
	if err != nil {
		return nil // no lenses found
	}

	var previous *lensVersionDef
	for lensIter.Next() {
		curr, err := newLensVersionDef(lensIter.Value())
		if err != nil {
			return err
		}

		if err := doCheck(previous, &curr); err != nil {
			return err
		}

		previous = &curr
	}

	return nil
}

func (ml *maybeLineage) checkGoLensCompleteness() error {
	// TODO(sdboyer) it'd be nice to consolidate all the errors so that the user always sees a complete set of problems
	all := make(map[lensID]bool)
	for _, lens := range ml.implens {
		id := lid(lens.From, lens.To)
		if all[id] {
			return fmt.Errorf("duplicate Go migration %s", id)
		}
		if lens.Mapper == nil {
			return fmt.Errorf("nil Go migration func for %s", id)
		}
		all[id] = true
	}

	var missing []lensID

	var prior SyntacticVersion
	for _, sch := range ml.schlist[1:] {
		// there must always at least be a reverse lens
		v := sch.Version()
		revid := lid(v, prior)

		if !all[revid] {
			missing = append(missing, revid)
		} else {
			delete(all, revid)
		}

		if v[0] != prior[0] {
			// if we crossed a major version, there must also be a forward lens
			fwdid := lid(prior, v)
			if !all[fwdid] {
				missing = append(missing, fwdid)
			} else {
				delete(all, fwdid)
			}
		}
		prior = v
	}

	// TODO is it worth making each sub-item into its own error type?
	if len(missing) > 0 {
		b := new(bytes.Buffer)

		fmt.Fprintf(b, "Go migrations not provided for the following version pairs:\n")
		for _, mlid := range missing {
			fmt.Fprint(b, "\t", mlid, "\n")
		}
		return errors.Mark(errors.New(b.String()), terrors.ErrMissingLenses)
	}

	if len(all) > 0 {
		b := new(bytes.Buffer)

		fmt.Fprintf(b, "Go migrations erroneously provided for the following version pairs:\n")
		// walk the slice so output is reliably ordered
		for _, lens := range ml.implens {
			// if it's not in the list it's because it was expected & already processed
			elid := lid(lens.From, lens.To)
			if _, has := all[elid]; !has {
				continue
			}
			if !synvExists(ml.allv, lens.To) {
				fmt.Fprintf(b, "\t%s (schema version %s does not exist)", elid, lens.To)
			} else if !synvExists(ml.allv, lens.From) {
				fmt.Fprintf(b, "\t%s (schema version %s does not exist)", elid, lens.From)
			} else if elid.To == elid.From {
				fmt.Fprintf(b, "\t%s (self-migrations not allowed)", elid)
			} else if elid.To.Less(elid.From) {
				// reverse lenses
				// only possibility is non-sequential versions connected
				fmt.Fprintf(b, "\t%s (%s is predecessor of %s, not %s)", elid, ml.allv[searchSynv(ml.allv, elid.From)-1], elid.From, elid.To)
			} else {
				// forward lenses
				// either a minor lens was provided, or non-sequential versions connected
				if lens.To[0] != lens.From[0] {
					fmt.Fprintf(b, "\t%s (minor version upgrades are handled automatically)", elid)
				} else {
					fmt.Fprintf(b, "\t%s (%s is successor of %s, not %s)", elid, ml.allv[searchSynv(ml.allv, elid.From)+1], elid.From, elid.To)
				}
			}
		}
		return errors.Mark(errors.New(b.String()), terrors.ErrErroneousLenses)
	}

	ml.lensmap = make(map[lensID]ImperativeLens, len(ml.implens))
	for _, lens := range ml.implens {
		ml.lensmap[lid(lens.From, lens.To)] = lens
	}

	return nil
}

type lensVersionDef struct {
	to   SyntacticVersion
	from SyntacticVersion
}

func newLensVersionDef(val cue.Value) (lensVersionDef, error) {
	v := lensVersionDef{}
	to, err := v.version(val, "to")
	if err != nil {
		return lensVersionDef{}, err
	}

	from, err := v.version(val, "from")
	if err != nil {
		return lensVersionDef{}, err
	}

	return lensVersionDef{to: to, from: from}, err
}

func doCheck(prev, curr *lensVersionDef) error {
	if prev == nil {
		return nil
	}

	if curr == nil {
		return nil
	}

	if curr.to.Less(prev.to) {
		return errors.Mark(
			errors.Errorf("lens version [to: %s, from: %s] is not greater than previous lens version [to: %s, from: %s]", curr.to, curr.from, prev.to, prev.from),
			terrors.ErrInvalidLensesOrder)
	}

	if prev.to == curr.to && curr.from.Less(prev.from) {
		return errors.Mark(
			errors.Errorf("lens version [to: %s, from: %s] is not greater than previous lens version [to: %s, from: %s]", curr.to, curr.from, prev.to, prev.from),
			terrors.ErrInvalidLensesOrder)
	}

	return nil
}

func (lensVersionDef) version(val cue.Value, p string) (SyntacticVersion, error) {
	vPath := cue.MakePath(cue.Str(p))
	vval := val.Value().LookupPath(vPath)

	v := SyntacticVersion{}
	if err := vval.Value().Decode(&v); err != nil {
		return v, errors.Mark(mkerror(val, fmt.Sprintf("failed to decode lens version %s from: %s", vval, val)), terrors.ErrInvalidLineage)
	}

	return v, nil
}

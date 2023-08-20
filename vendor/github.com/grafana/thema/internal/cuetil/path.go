package cuetil

import "cuelang.org/go/cue"

// TrimPathPrefix strips the provided prefix from the provided path, if the
// prefix exists.
//
// If path and prefix are equivalent, and there is at least one additional
// selector in the provided path.
func TrimPathPrefix(path, prefix cue.Path) cue.Path {
	sels, psels := path.Selectors(), prefix.Selectors()
	if len(sels) == 1 {
		return path
	}
	var i int
	for ; i < len(psels) && i < len(sels); i++ {
		if !SelEq(psels[i], sels[i]) {
			break
		}
	}
	return cue.MakePath(sels[i:]...)
}

// ReplacePathPrefix replaces oldprefix with newprefix if p begins with oldprefix.
func ReplacePathPrefix(p, oldprefix, newprefix cue.Path) cue.Path {
	ps, ops, nps := p.Selectors(), oldprefix.Selectors(), newprefix.Selectors()
	if len(ops) >= len(ps) {
		return p
	}

	if !PathsAreEq(cue.MakePath(ps[:len(ops)]...), cue.MakePath(ops...)) {
		return p
	}

	pn := make([]cue.Selector, len(nps)+len(ps)-len(ops))
	copy(pn, nps)
	copy(pn[len(nps):], ps[len(ops):])
	return cue.MakePath(pn...)
}

// PathsAreEq tests whether two [cue.Path] are equivalent. Paths that vary only by
// optionality are considered equivalent.
func PathsAreEq(p1, p2 cue.Path) bool {
	return pathsAreEq(p1.Selectors(), p2.Selectors())
}

func pathsAreEq(p1s, p2s []cue.Selector) bool {
	if len(p1s) != len(p2s) {
		return false
	}
	for i := 0; i < len(p2s); i++ {
		if !SelEq(p2s[i], p1s[i]) {
			return false
		}
	}
	return true
}

// PathHasPrefix tests whether the [cue.Path] p begins with prefix.
func PathHasPrefix(p, prefix cue.Path) bool {
	ps, pres := p.Selectors(), prefix.Selectors()
	if len(pres) > len(ps) {
		return false
	}
	return pathsAreEq(ps[:len(pres)], pres)
}

// LastSelectorEq tests whether the final selector in the provided path is
// equivalent to the provided selector. Selectors that vary only by optionality
// are considered equivalent.
func LastSelectorEq(p cue.Path, sel cue.Selector) bool {
	sels := p.Selectors()
	last := sels[len(sels)-1]
	return SelEq(last, sel)
}

// SelEq indicates whether two selectors are equivalent. Selectors are equivalent if
// they are either exactly equal, or if they are equal ignoring path optionality.
func SelEq(s1, s2 cue.Selector) bool {
	return s1 == s2 || s1.Optional() == s2.Optional()
}

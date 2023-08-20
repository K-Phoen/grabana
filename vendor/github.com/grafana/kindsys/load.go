package kindsys

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"
	"testing/fstest"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	"github.com/grafana/thema/load"
	"github.com/yalue/merged_fs"
)

var defaultFramework cue.Value
var fwOnce sync.Once
var ctx = cuecontext.New()

// cueContext returns a singleton instance of [cue.Context].
func cueContext() *cue.Context {
	return ctx
}

func init() {
	loadpFrameworkOnce()
}

func loadpFrameworkOnce() {
	fwOnce.Do(func() {
		var err error
		defaultFramework, err = doLoadFrameworkCUE(cueContext())
		if err != nil {
			panic(err)
		}
		ctx = cuecontext.New()
	})
}

func doLoadFrameworkCUE(ctx *cue.Context) (cue.Value, error) {
	v, err := BuildInstance(ctx, ".", "kindsys", nil)
	if err != nil {
		return v, err
	}

	if err = v.Validate(cue.Concrete(false), cue.All()); err != nil {
		return cue.Value{}, fmt.Errorf("kindsys framework loaded cue.Value has err: %w", err)
	}

	return v, nil
}

func BuildInstance(ctx *cue.Context, relpath string, pkg string, overlay fs.FS) (cue.Value, error) {
	bi, err := LoadInstance(relpath, pkg, overlay)
	if err != nil {
		return cue.Value{}, err
	}

	if ctx == nil {
		ctx = cueContext()
	}

	v := ctx.BuildInstance(bi)
	if v.Err() != nil {
		return v, fmt.Errorf("%s not a valid CUE instance: %w", relpath, v.Err())
	}
	return v, nil
}

// LoadInstance returns a build.Instance populated with the CueSchemaFS at the root and
// an optional overlay filesystem.
func LoadInstance(relpath string, pkg string, overlay fs.FS) (*build.Instance, error) {
	relpath = filepath.ToSlash(relpath)

	var f fs.FS = CueSchemaFS
	var err error
	if overlay != nil {
		f, err = prefixWithCUE(relpath, overlay)
		if err != nil {
			return nil, err
		}
	}

	if pkg != "" {
		return load.InstanceWithThema(f, relpath, load.Package(pkg))
	}
	return load.InstanceWithThema(f, relpath)
}

// prefixWithCUE constructs an fs.FS that merges the provided fs.FS with the
// embedded FS containing kindsys core CUE files, CueSchemaFS. The provided
// prefix should be the relative path from the repository root to the directory
// root of the provided inputfs.
//
// The returned fs.FS is suitable for passing to a CUE loader, such as
// [load.InstanceWithThema].
func prefixWithCUE(prefix string, inputfs fs.FS) (fs.FS, error) {
	m, err := prefixFS(prefix, inputfs)
	if err != nil {
		return nil, err
	}
	return merged_fs.NewMergedFS(m, CueSchemaFS), nil
}

// TODO such a waste, replace with stateless impl that just transforms paths on the fly
func prefixFS(prefix string, fsys fs.FS) (fs.FS, error) {
	m := make(fstest.MapFS)

	prefix = filepath.FromSlash(prefix)
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		b, err := fs.ReadFile(fsys, filepath.ToSlash(path))
		if err != nil {
			return err
		}
		// fstest can recognize only forward slashes.
		m[filepath.ToSlash(filepath.Join(prefix, path))] = &fstest.MapFile{Data: b}
		return nil
	})
	return m, err
}

// CUEFramework returns a cue.Value representing all the kindsys framework raw
// CUE files.
//
// For low-level use in constructing other types and APIs, while still letting
// us define all the frameworky CUE bits in a single package. Other Go types
// make the constructs in the returned cue.Value easy to use.
//
// Calling this with a nil [cue.Context] (the singleton returned from
// [CUEContext]) will memoize certain CUE operations. Prefer passing nil unless
// a different cue.Context is specifically required.
func CUEFramework(ctx *cue.Context) cue.Value {
	if ctx == nil || ctx == cueContext() {
		// Ensure framework is loaded, even if this func is called
		// from an init() somewhere.
		loadpFrameworkOnce()
		return defaultFramework
	}
	// Error guaranteed to be nil here because erroring would have caused init() to panic
	v, _ := doLoadFrameworkCUE(ctx) // nolint:errcheck
	return v
}

// ToKindProps takes a cue.Value expected to represent a kind of the category
// specified by the type parameter and populates the Go type from the cue.Value.
func ToKindProps[T KindProperties](v cue.Value) (T, error) {
	props := new(T)
	if !v.Exists() {
		return *props, ErrValueNotExist
	}

	fw := CUEFramework(v.Context())
	var kdef cue.Value

	anyprops := any(*props).(SomeKindProperties)
	switch anyprops.(type) {
	case CoreProperties:
		kdef = fw.LookupPath(cue.MakePath(cue.Str("Core")))
	case CustomProperties:
		kdef = fw.LookupPath(cue.MakePath(cue.Str("Custom")))
	case ComposableProperties:
		kdef = fw.LookupPath(cue.MakePath(cue.Str("Composable")))
	default:
		// unreachable so long as all the possibilities in KindProperties have switch branches
		panic("unreachable")
	}

	item := v.Unify(kdef)
	if item.Err() != nil {
		return *props, errors.Wrap(errors.Promote(ErrValueNotAKind, ""), item.Err())
	}

	if err := item.Decode(props); err != nil {
		// Should only be reachable if CUE and Go framework types have diverged
		panic(errors.Details(err, nil))
	}

	return *props, nil
}

// ToDef takes a cue.Value expected to represent a kind of the category
// specified by the type parameter and populates a Def from the CUE value.
// The cue.Value in Def.V will be the unified value of the parameter cue.Value
// and the kindsys CUE kind (Core, Custom, Composable).
func ToDef[T KindProperties](v cue.Value) (Def[T], error) {
	def := Def[T]{}
	props := new(T)
	if !v.Exists() {
		return def, ErrValueNotExist
	}

	fw := CUEFramework(v.Context())
	var kdef cue.Value

	anyprops := any(*props).(SomeKindProperties)
	switch anyprops.(type) {
	case CoreProperties:
		kdef = fw.LookupPath(cue.MakePath(cue.Str("Core")))
	case CustomProperties:
		kdef = fw.LookupPath(cue.MakePath(cue.Str("Custom")))
	case ComposableProperties:
		kdef = fw.LookupPath(cue.MakePath(cue.Str("Composable")))
	default:
		// unreachable so long as all the possibilities in KindProperties have switch branches
		panic("unreachable")
	}

	def.V = v.Unify(kdef)
	if def.V.Err() != nil {
		return def, errors.Wrap(errors.Promote(ErrValueNotAKind, ""), def.V.Err())
	}

	if err := def.V.Decode(props); err != nil {
		// Should only be reachable if CUE and Go framework types have diverged
		panic(errors.Details(err, nil))
	}
	def.Properties = *props
	return def, nil
}

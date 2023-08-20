package thema

import (
	"fmt"
	"path/filepath"
	"sync"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/load"
	"github.com/grafana/thema/internal/util"
)

var rtOnce sync.Once
var themaBI *build.Instance

func loadRuntime() *build.Instance {
	rtOnce.Do(func() {
		path := filepath.Join(util.Prefix, "github.com", "grafana", "thema")

		overlay := make(map[string]load.Source)
		if err := util.ToOverlay(path, CueJointFS, overlay); err != nil {
			// It's impossible for this to fail barring temporary bugs in filesystem
			// layout within the thema rt itself. These should be trivially
			// catchable during CI, so avoid forcing meaningless error handling on
			// dependers and prefer a panic.
			panic(err)
		}

		cfg := &load.Config{
			Overlay: overlay,
			Package: "thema",
			Module:  "github.com/grafana/thema",
			Dir:     path,
		}
		themaBI = load.Instances(nil, cfg)[0]

		// proactively check, so we don't have to do it when making a new library
		rt := cuecontext.New().BuildInstance(themaBI)
		if err := rt.Validate(cue.All()); err != nil {
			// As with the above, an error means that a problem exists in the
			// literal CUE code embedded in this version of package (that should
			// have trivially been caught with CI), so the caller can't fix anything
			// without changing the version of the thema Go library they're
			// depending on. It's a hard failure that should be unreachable outside
			// thema internal testing, so just panic.
			panic(errors.Details(err, nil))
		}
	})

	return themaBI
}

// Runtime holds the set of CUE constructs available in the Thema CUE package,
// allowing Thema's Go code to internally reuse the same native CUE functionality.
//
// Each Thema Runtime is bound to a single cue.Context, determined by the parameter
// passed to [NewRuntime].
type Runtime struct {
	// Value corresponds to loading the whole github.com/grafana/thema:thema
	// package.
	val cue.Value

	// Until CUE is safe for certain concurrent operations, keep a mutex to
	// help guard...at least somewhat.
	mut sync.RWMutex
}

// NewRuntime parses, loads and builds a full CUE instance/value representing
// all of the logic in the Thema CUE package (github.com/grafana/thema),
// and returns a Runtime instance ready for use.
//
// Building is performed using the provided cue.Context. Passing a nil context will panic.
//
// This function is the canonical way to make Thema logic callable from Go code.
func NewRuntime(ctx *cue.Context) *Runtime {
	if ctx == nil {
		panic("nil context provided")
	}
	rt := ctx.BuildInstance(loadRuntime())

	// FIXME preload all the known funcs into a map[string]cue.Value here to avoid runtime cost
	return &Runtime{
		val: rt,
	}
}

func (rt *Runtime) rl() {
	rt.mut.RLock()
	// rt.mut.Lock()
}

func (rt *Runtime) ru() {
	rt.mut.RUnlock()
	// rt.mut.Unlock()
}

func (rt *Runtime) l() {
	rt.mut.Lock()
}

func (rt *Runtime) u() {
	rt.mut.Unlock()
}

// Underlying returns the underlying cue.Value representing the whole Thema CUE
// library (github.com/grafana/thema).
func (rt *Runtime) Underlying() cue.Value {
	return rt.val
}

// Context returns the *cue.Context in which this runtime was built.
func (rt *Runtime) Context() *cue.Context {
	return rt.val.Context()
}

// Return the #Lineage definition (or panic)
//
// SURROUND CALLS TO THIS IN rl()/ru()
func (rt *Runtime) linDef() cue.Value {
	dlin := rt.val.LookupPath(cue.MakePath(cue.Def("#Lineage")))
	return dlin
}

type cueArgs map[string]interface{}

func (ca cueArgs) make(path string, rt *Runtime) (cue.Value, error) {
	rt.l()
	defer rt.u()

	var cpath cue.Path
	if path[0] == '_' {
		cpath = cue.MakePath(cue.Hid(path, "github.com/grafana/thema"))
	} else {
		cpath = cue.ParsePath(path)
	}
	cfunc := rt.val.LookupPath(cpath)
	if !cfunc.Exists() {
		panic(fmt.Sprintf("cannot call nonexistent CUE func %q", path))
	}
	if cfunc.Err() != nil {
		panic(cfunc.Err())
	}

	applic := []cue.Value{cfunc}
	for arg, val := range ca {
		p := cue.ParsePath(arg)
		step := applic[len(applic)-1]
		if !step.Allows(p.Selectors()[0]) {
			panic(fmt.Sprintf("CUE func %q does not take an argument named %q", path, arg))
		}
		applic = append(applic, step.FillPath(p, val))
	}
	last := applic[len(applic)-1]

	// Have to do the error check in a separate loop after all args are applied,
	// because args may depend on each other and erroneously error depending on
	// the order of application.
	for arg := range ca {
		argv := last.LookupPath(cue.ParsePath(arg))
		if argv.Err() != nil {
			return cue.Value{}, &errInvalidCUEFuncArg{
				cuefunc: path,
				argpath: arg,
				err:     argv.Err(),
			}
		}
	}
	return last, nil
}

func (ca cueArgs) call(path string, rt *Runtime) (cue.Value, error) {
	v, err := ca.make(path, rt)
	if err != nil {
		return cue.Value{}, err
	}
	rt.rl()
	rv := v.LookupPath(outpath)
	rt.ru()
	return rv, nil
}

type errInvalidCUEFuncArg struct {
	cuefunc string
	argpath string
	err     error
}

func (e *errInvalidCUEFuncArg) Error() string {
	return fmt.Sprintf("err on arg %q to CUE func %q: %s", e.argpath, e.cuefunc, errors.Details(e.err, nil))
}

var outpath = cue.MakePath(cue.Str("out"))

package load

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing/fstest"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/grafana/thema"
	"github.com/grafana/thema/internal/util"
	"github.com/yalue/merged_fs"
)

// ErrFSNotACueModule is a general error that wraps a particular error that explains
// why a particular fs.FS cannot be used as a CUE module FS as needed by InstancesWithThema.
type ErrFSNotACueModule struct {
	fserr error
}

func (e *ErrFSNotACueModule) Error() string {
	return fmt.Sprintf("provided fs.FS is not a valid CUE module: %q", e.fserr.Error())
}

func (e *ErrFSNotACueModule) Unwrap() error {
	return e.fserr
}

// AsModFS injects a cue.mod/module.cue file into the provided in fs.FS at its root,
// preparing it for use by
//
// If a cue.mod already exists in the provided fs.FS, this func is a no-op.
func AsModFS(in fs.FS, modname string) fs.FS {
	_, err := fs.Stat(in, filepath.Join("cue.mod", "module.cue"))
	if os.IsNotExist(err) {
		m := fstest.MapFS{
			// in general, Go virtual FS work only with slash separators internally
			"cue.mod/module.cue": &fstest.MapFile{Data: []byte(fmt.Sprintf(`module: "%s"`, modname))},
		}
		return merged_fs.NewMergedFS(m, in)
	}
	return in
}

var themamodpath = filepath.Join("cue.mod", "pkg", "github.com", "grafana", "thema")

// InstanceWithThema wraps CUE's [load.Instance] in order to allow
// loading .cue files that directly `import "github.com/grafana/thema"`, as
// lineages are expected to. This is accomplished by constructing a
// load.Config.Overlay with the Thema CUE files dynamically injected under
// cue.mod/pkg/, where CUE searches for mod-external imports.
//
// This loader is opinionated, preferring simple ease-of-use and fewer degrees
// of freedom to the completeness of load.Instances. Some reasonable use cases
// may not be achievable. Make your own as needed - all key component parts are
// exported from elsewhere in the Thema Go module.
//
// The modFS is expected to be an fs.FS containing the cue.mod module metadata,
// and any lineage(s) to be loaded.
//
// The root of the FS must be an importable CUE module with a path. That is,
// there must exist cue.mod/module.cue, which must contain a top-level field
// declaring the module name (aka import prefix/module path), e.g.:
//
//	module: "github.com/grafana/thema"
//
// The dir parameter must specify a directory containing .cue files with
// lineages to be loaded, relative to the module root directory. This is similar
// to load.Config.Dir, except:
//   - The package name of the root dir is the final element of the module name,
//     unless overridden with a [Package] option.
//   - "." and the empty string are a special value that will load the root
//     directory of the modFS.
//
// NOTE - this function is likely to be deprecated and removed in favor of a
// more generic dependency overlay loader.
func InstanceWithThema(modFS fs.FS, dir string, opts ...Option) (*build.Instance, error) {
	var modname string
	err := fs.WalkDir(modFS, "cue.mod", func(path string, d fs.DirEntry, err error) error {
		// fs.FS implementations tend to not use path separators as expected. Use a
		// normalized one for comparisons, but retain the original for calls back into modFS.
		normpath := filepath.FromSlash(path)
		if err != nil {
			return err
		}

		if d.IsDir() {
			switch normpath {
			case filepath.Join("cue.mod", "gen"), filepath.Join("cue.mod", "usr"):
				return fs.SkipDir
			case themamodpath:
				return fmt.Errorf("path %q already exists in modFS passed to InstancesWithThema, must be absent for dynamic dependency injection", themamodpath)
			}
			return nil
		} else if normpath == filepath.Join("cue.mod", "module.cue") {
			modf, err := modFS.Open(path)
			if err != nil {
				return err
			}
			defer modf.Close() // nolint: errcheck

			b, err := io.ReadAll(modf)
			if err != nil {
				return err
			}

			modname, err = cuecontext.New().CompileBytes(b).LookupPath(cue.MakePath(cue.Str("module"))).String()
			if err != nil {
				return err
			}
			if modname == "" {
				return fmt.Errorf("InstancesWithThema requires non-empty module name in modFS' cue.mod/module.cue")
			}
		}

		return nil
	})
	if err != nil {
		return nil, &ErrFSNotACueModule{fserr: err}
	}

	if modname == "" {
		return nil, &ErrFSNotACueModule{fserr: fmt.Errorf("cue.mod/module.cue modname was empty")}
	}

	modroot := filepath.FromSlash(filepath.Join(util.Prefix, modname))
	overlay := make(map[string]load.Source)
	if err := util.ToOverlay(modroot, modFS, overlay); err != nil {
		return nil, err
	}

	// Special case for when we're calling this loader with paths inside the thema module
	if modname == "github.com/grafana/thema" {
		if err := ToOverlay(modroot, thema.CueJointFS, overlay); err != nil {
			return nil, err
		}
	} else {
		if err := ToOverlay(filepath.Join(modroot, themamodpath), thema.CueFS, overlay); err != nil {
			return nil, err
		}
	}

	cfg := &load.Config{
		Overlay:    overlay,
		ModuleRoot: modroot,
		Module:     modname,
		Dir:        filepath.Join(modroot, dir),
		DataFiles:  true,
	}
	lc := &loadConfig{
		pkgname: filepath.Base(dir),
	}

	if dir == "" || dir == "." {
		lc.pkgname = filepath.Base(modroot)
		cfg.Dir = modroot
	}
	for _, opt := range opts {
		opt(lc)
	}

	cfg.Package = lc.pkgname

	inst := load.Instances(nil, cfg)[0]
	if inst.Err != nil {
		return nil, inst.Err
	}

	return inst, nil
}

// InstancesWithThema passes through to [InstanceWithThema].
// DEPRECATED: use InstanceWithThema.
func InstancesWithThema(modFS fs.FS, dir string, opts ...Option) (*build.Instance, error) {
	return InstanceWithThema(modFS, dir, opts...)
}

// Option defines optional parameters that may be passed to the loader.
type Option loadOption

// Internal representation of Option.
type loadOption func(c *loadConfig)

type loadConfig struct {
	pkgname string
}

// Package specifies a custom CUE package name use when loading CUE files.
// See ["cuelang.org/go/cue/load".Config.Package].
func Package(name string) Option {
	return func(c *loadConfig) {
		c.pkgname = name
	}
}

// ToOverlay maps the provided fs.FS into an Overlay for use in load.Config.
//
// An absolute path prefix must be provided.
func ToOverlay(prefix string, vfs fs.FS, overlay map[string]load.Source) error {
	return util.ToOverlay(prefix, vfs, overlay)
}

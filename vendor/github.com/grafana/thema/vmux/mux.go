package vmux

import (
	"encoding/json"
	"strings"

	"cuelang.org/go/cue"
	cjson "cuelang.org/go/encoding/json"
	"cuelang.org/go/encoding/yaml"
	pyaml "cuelang.org/go/pkg/encoding/yaml"
	"github.com/grafana/thema"
)

func allvstr(sch thema.Schema) string {
	var vl []string
	for isch := thema.SchemaP(sch.Lineage(), thema.SV(0, 0)); isch != nil; isch = isch.Successor() {
		vl = append(vl, isch.Version().String())
	}
	return strings.Join(vl, ", ")
}

func latest(lin thema.Lineage) thema.Schema {
	return thema.SchemaP(lin, thema.LatestVersion(lin))
}

// A Decoder can decode a []byte in a particular format (e.g. JSON, YAML) into a
// [cue.Value], readying it for a call to [thema.Schema.Validate].
type Decoder interface {
	Decode(ctx *cue.Context, b []byte) (cue.Value, error)
}

// An Encoder can encode a [thema.Instance] to a []byte in a particular format
// (e.g. JSON, YAML).
type Encoder interface {
	Encode(cue.Value) ([]byte, error)
}

// A Codec can decode a []byte in a particular format (e.g. JSON, YAML) into
// CUE, and decode from a [thema.Instance] back into a []byte.
//
// It is customary, but not necessary, that a Codec's input and output formats
// are the same.
type Codec interface {
	Decoder
	Encoder
}

type jsonCodec struct {
	path string
}

// NewJSONCodec creates a [Codec] that decodes from and encodes to a JSON []byte.
//
// The provided path is used as the CUE source path for each []byte input
// passed through the decoder. These paths do not affect behavior, but show up
// in error output (e.g. validation).
func NewJSONCodec(path string) Codec {
	return jsonCodec{
		path: path,
	}
}

func (e jsonCodec) Decode(ctx *cue.Context, data []byte) (cue.Value, error) {
	expr, err := cjson.Extract(e.path, data)
	if err != nil {
		return cue.Value{}, err
	}
	return ctx.BuildExpr(expr), nil
}

func (e jsonCodec) Encode(v cue.Value) ([]byte, error) {
	return json.Marshal(v)
}

type yamlCodec struct {
	path string
}

// NewYAMLCodec creates a [Codec] that decodes from and encodes to a YAML []byte.
//
// The provided path is used as the CUE source path for each []byte input
// passed through the decoder. These paths do not affect behavior, but show up
// in error output (e.g. validation).
func NewYAMLCodec(path string) Codec {
	return yamlCodec{
		path: path,
	}
}

func (e yamlCodec) Decode(ctx *cue.Context, data []byte) (cue.Value, error) {
	expr, err := yaml.Extract(e.path, data)
	if err != nil {
		return cue.Value{}, err
	}
	return ctx.BuildFile(expr), nil
}

func (e yamlCodec) Encode(v cue.Value) ([]byte, error) {
	s, err := pyaml.Marshal(v)
	return []byte(s), err
}

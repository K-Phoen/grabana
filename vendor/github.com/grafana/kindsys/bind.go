package kindsys

import (
	"cuelang.org/go/cue"
	"encoding/json"
	"github.com/grafana/kindsys/encoding"
	"github.com/grafana/thema"
	"github.com/grafana/thema/vmux"
)

type gis struct {
	// TODO
	Spec json.RawMessage `json:"spec"`
	// TODO
	Metadata json.RawMessage `json:"metadata"`
	// TODO
	//CustomMetadata json.RawMessage `json:"customMetadata"`
}

type withLineage interface {
	Kind
	Group() string
	Lineage() thema.Lineage
}

func bytesToAnyInstance(k withLineage, b []byte, codec Decoder) (*thema.Instance, error) {
	// Transform from k8s shape to intermediate grafana shape
	var gb encoding.GrafanaShapeBytes
	gb, err := codec.Decode(b)
	if err != nil {
		return nil, err
	}
	//if gb.Group != k.Group() || gb.Kind != k.Name() {
	//	return nil, fmt.Errorf("resource is %s.%s, not of kind %s.%s", gb.Group, gb.Kind, k.Group(), k.Name())
	//}
	// TODO make the intermediate type already look like this so we don't have to re-encode/decode
	gj := gis{
		Spec:     gb.Spec,
		Metadata: gb.Metadata,
		// TODO Status?
	}
	gjb, err := json.Marshal(gj)
	if err != nil {
		return nil, err
	}

	lin := k.Lineage()
	// reuse the cue context already attached to the underlying lineage
	ctx := lin.Runtime().Context()
	// decode JSON into a cue.Value
	cval, err := vmux.NewJSONCodec(k.MachineName()+".json").Decode(ctx, gjb)
	if err != nil {
		return nil, err
	}

	// TODO take advantage of apiVersion of object to pick the right schema to validate against
	sch, _ := lin.Schema(k.CurrentVersion()) // we verified at bind of this kind that this schema exists
	inst, curvererr := sch.Validate(cval)
	if curvererr != nil {
		for sch := lin.First(); sch != nil; sch = sch.Successor() {
			if sch.Version() == k.CurrentVersion() {
				continue
			}
			if inst, err = sch.Validate(cval); err == nil {
				curvererr = nil
				break
			}
		}
	}

	// TODO improve this once thema stacks all schema validation errors https://github.com/grafana/thema/issues/156
	return inst, curvererr
}

// TODO this is why we need to combine [Core] and [Custom]
type resourceKind interface {
	Kind
	Group() string
}

type grafanaShape struct {
	Kind       string         `json:"kind"`
	Group      string         `json:"group"`
	Version    string         `json:"apiVersion"`
	Spec       map[string]any `json:"spec"`
	CustomMeta map[string]any `json:"customMetadata"`
	//Metadata   map[string]any `json:"metadata"`
	Metadata map[string]any `json:"metadata"`
}

// FIXME this is a fugly temporary hack - make this go away when we have clarity on our different shapes and the types line up
func grafanaShapeToUnstructured(k resourceKind, inst *thema.Instance) (*UnstructuredResource, error) {
	gs := grafanaShape{}
	err := inst.Underlying().Decode(&gs)
	if err != nil {
		return nil, err
	}

	u := &UnstructuredResource{}

	u.StaticMeta.Group = k.Group()
	u.StaticMeta.Kind = k.Name()
	u.StaticMeta.Version = gs.Version
	// TODO what are we doing about namespace?
	if ns, has := gs.Metadata["namespace"]; has {
		u.StaticMeta.Namespace = ns.(string)
	}
	// TODO is this the place where we could universally handle name generation?
	if ns, has := gs.Metadata["name"]; has {
		u.StaticMeta.Name = ns.(string)
	}

	// Just re-decode directly into the CommonMetadata. (lol horrible duplication, kill all of this ASAP)
	err = inst.Underlying().LookupPath(cue.ParsePath("metadata")).Decode(&u.CommonMeta)
	if err != nil {
		return nil, err
	}
	// NOTE this doesn't populate anything right now
	u.CustomMeta = gs.CustomMeta

	return u, nil
}

package encoding

// KubernetesShapeBytes are a partially-encoded representation of a []byte that
// is in the standard Kubernetes shape.
type KubernetesShapeBytes struct {
	// Spec contains the marshaled SpecObject. It should be unmarshalable directly into the Resource-implementation's
	// Spec object using an unmarshaler of the appropriate WireFormat type
	Spec []byte
	// Metadata includes object-specific metadata, and may include CommonMetadata depending on implementation.
	// Clients must call SetCommonMetadata on the object after an Unmarshal if CommonMetadata is not provided in the bytes.
	Metadata []byte
	// Subresources contains a map of all subresources that are both part of the underlying Object implementation,
	// AND are supported by the Client implementation. Each entry should be unmarshalable directly into the
	// Object-implementation's relevant subresource using an unmarshaler of the appropriate WireFormat type
	Subresources map[string][]byte
}

// GrafanaShapeBytes is a collection of bytes encoded in some wire format which can be unmarshaled into
// each component of a kindsys.Resource. Also included are metadata about the kind and resource,
// which can be used to identify what kind and resource to use for unmarshaling.
type GrafanaShapeBytes struct {
	// Kind is the name of the kind for which these bytes were extracted from or could be composed into
	Kind string
	// Group is the group to which the kind belongs. Together with Kind this becomes a globally-unique identifier for the kind
	Group string
	// Version is the particular version of the kind these bytes are for.
	Version string
	// TODO
	Spec []byte
	// TODO
	Metadata []byte
	// TODO
	CustomMetadata []byte
	// TODO
	Subresources map[string][]byte
}

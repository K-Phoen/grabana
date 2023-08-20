package encoding

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO: remove--this is here to avoid import cycle (unexported to avoid accidental imports from users)
type commonMetadata struct {
	// UID is the unique ID of the object. This can be used to uniquely identify objects,
	// but is not guaranteed to be usable for lookups.
	UID string `json:"uid"`
	// ResourceVersion is a version string used to identify any and all changes to the object.
	// Any time the object changes in storage, the ResourceVersion will be changed.
	// This can be used to block updates if a change has been made to the object between when the object was
	// retrieved, and when the update was applied.
	ResourceVersion string `json:"resourceVersion"`
	// Labels are string key/value pairs attached to the object. They can be used for filtering,
	// or as additional metadata.
	Labels map[string]string `json:"labels,omitempty"`
	// CreationTimestamp indicates when the resource has been created.
	CreationTimestamp time.Time `json:"creationTimestamp"`
	// DeletionTimestamp indicates that the resource is pending deletion as of the provided time if non-nil.
	// Depending on implementation, this field may always be nil, or it may be a "tombstone" indicator.
	// It may also indicate that the system is waiting on some task to finish before the object is fully removed.
	DeletionTimestamp *time.Time `json:"deletionTimestamp,omitempty"`
	// Finalizers are a list of identifiers of interested parties for delete events for this resource.
	// Once a resource with finalizers has been deleted, the object should remain in the store,
	// DeletionTimestamp is set to the time of the "delete," and the resource will continue to exist
	// until the finalizers list is cleared.
	Finalizers []string `json:"finalizers,omitempty"`
	// UpdateTimestamp is the timestamp of the last update to the resource.
	UpdateTimestamp time.Time `json:"updateTimestamp"`
	// CreatedBy is a string which indicates the user or process which created the resource.
	// Implementations may choose what this indicator should be.
	CreatedBy string `json:"createdBy"`
	// UpdatedBy is a string which indicates the user or process which last updated the resource.
	// Implementations may choose what this indicator should be.
	UpdatedBy string `json:"updatedBy"`

	// ExtraFields stores implementation-specific metadata.
	// Not all Client implementations are required to honor all ExtraFields keys.
	// Generally, this field should be shied away from unless you know the specific
	// Client implementation you're working with and wish to track or mutate extra information.
	ExtraFields map[string]any `json:"extraFields"`
}

const annotationPrefix = "grafana.com/"

// KubernetesJSONEncoder is a kubernetes encoder for JSON wire format
type KubernetesJSONEncoder struct{}

// Encode accepts GrafanaShapedBytes which are in a JSON wire format,
// and produces a JSON-encoded kubernetes payload
func (k *KubernetesJSONEncoder) Encode(bytes GrafanaShapeBytes) ([]byte, error) {
	// Partly unmarshal the metadata in the GrafanaShapeBytes
	partial := make(map[string]json.RawMessage)
	err := json.Unmarshal(bytes.Metadata, &partial)
	if err != nil {
		return nil, fmt.Errorf("unable to parse metadata: %w", err)
	}
	// Move everything in extraFields to the top, then remove the extraFields key entirely
	if extraFieldsRaw, ok := partial["extraFields"]; ok {
		extraFields := make(map[string]json.RawMessage)
		err = json.Unmarshal(extraFieldsRaw, &extraFields)
		for key, val := range extraFields {
			partial[key] = val
		}
		delete(partial, "extraFields")
	}
	// Move every non-kubernetes key into the annotations
	annotations := make(map[string]string)
	if annotationsRaw, ok := partial["annotations"]; ok {
		err = json.Unmarshal(annotationsRaw, &annotations)
		if err != nil {
			return nil, err
		}
	}
	// TODO: make this dynamic instead of hard-coding field names
	annotations[annotationPrefix+"createdBy"] = rawToString(partial["createdBy"])
	delete(partial, "createdBy")
	annotations[annotationPrefix+"updatedBy"] = rawToString(partial["updatedBy"])
	delete(partial, "updatedBy")
	annotations[annotationPrefix+"updateTimestamp"] = rawToString(partial["updateTimestamp"])
	delete(partial, "updateTimestamp")

	if len(bytes.CustomMetadata) > 0 {
		custom := make(map[string]any)
		err = json.Unmarshal(bytes.CustomMetadata, &custom)
		if err != nil {
			return nil, fmt.Errorf("unable to parse custom metadata: %w", err)
		}
		for key, val := range custom {
			annotations[annotationPrefix+key] = anyToString(val)
		}
	}
	// Re-encode the annotations
	partial["annotations"], err = json.Marshal(annotations)
	if err != nil {
		return nil, err
	}

	// Package the bytes for kubernetes
	kube := make(map[string]json.RawMessage)
	kube["kind"], err = json.Marshal(bytes.Kind)
	if err != nil {
		return nil, err
	}
	kube["apiVersion"], err = json.Marshal(fmt.Sprintf("%s/%s", bytes.Group, bytes.Version))
	if err != nil {
		return nil, err
	}
	kube["metadata"], err = json.Marshal(partial)
	if err != nil {
		return nil, err
	}
	kube["spec"] = bytes.Spec
	for key, val := range bytes.Subresources {
		kube[key] = val
	}
	return json.Marshal(kube)
}

func rawToString(raw json.RawMessage) string {
	var val any
	json.Unmarshal(raw, &val)
	return anyToString(val)
}

func anyToString(val any) string {
	v := reflect.ValueOf(val)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.String, reflect.Int, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64, reflect.Bool:
		return fmt.Sprintf("%v", v.Interface())
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return "" // Invalid kind to encode
	default:
		bytes, _ := json.Marshal(val)
		return string(bytes)
	}
}

// KubernetesJSONDecoder is a kubernetes decoder for JSON wire format
type KubernetesJSONDecoder struct{}

// This is a bit hacky, but better than hard-coding keys, so it doesn't need to be updated if CommonMetadata changes
var (
	commonMetadataKeys              = make(map[string]any)
	marshaledEmptyCommonMetadata, _ = json.Marshal(commonMetadata{})
	_                               = json.Unmarshal(marshaledEmptyCommonMetadata, &commonMetadataKeys)
)

// Decode accepts JSON-encoded bytes of a kubernetes object,
// and returns JSON-encoded GrafanaShapeBytes of that object
func (k *KubernetesJSONDecoder) Decode(bytes []byte) (GrafanaShapeBytes, error) {
	// Decode the bytes into a KubeObject
	// This will pass-through the bytes for the spec and all subresources, but unmarshal the metadata
	// We have to unmarshal the metadata so we can properly translate it, then re-encode
	partial := make(map[string]json.RawMessage)
	err := json.Unmarshal(bytes, &partial)
	if err != nil {
		return GrafanaShapeBytes{}, err
	}
	res := GrafanaShapeBytes{
		Subresources: make(map[string][]byte),
	}
	// TODO(IfSentient): there is a more efficient way to do this by again only partially unmarshaling,
	// but we need a good way to determine types of fields for the CommonMetadata (and eventually, CustomMetadata)
	// when extracting from annotations (right now, we hard-code type conversion for CommonMetadata keys later on)
	kubeMeta := metav1.ObjectMeta{}
	for key, val := range partial {
		switch key {
		case "metadata":
			err = json.Unmarshal(val, &kubeMeta)
			if err != nil {
				return res, fmt.Errorf("unable to decode kubernetes metadata: %w", err)
			}
		case "spec":
			res.Spec = val
		case "apiVersion":
			s := ""
			err = json.Unmarshal(val, &s)
			if err != nil {
				return res, err
			}
			apiVersion := strings.Split(s, "/")
			res.Group = apiVersion[0]
			res.Version = apiVersion[1]
		case "kind":
			s := ""
			err = json.Unmarshal(val, &s)
			if err != nil {
				return res, err
			}
			res.Kind = s
		default:
			res.Subresources[key] = val
		}
	}

	// Break up the metadata into common and custom, then re-encode both
	cmd := commonMetadata{
		UID:               string(kubeMeta.UID),
		ResourceVersion:   kubeMeta.ResourceVersion,
		Labels:            kubeMeta.Labels,
		CreationTimestamp: kubeMeta.CreationTimestamp.Time.UTC(),
		Finalizers:        kubeMeta.Finalizers,
	}
	if kubeMeta.OwnerReferences != nil && len(kubeMeta.OwnerReferences) > 0 {
		cmd.ExtraFields["ownerReferences"] = kubeMeta.OwnerReferences
	}
	if kubeMeta.ManagedFields != nil && len(kubeMeta.ManagedFields) > 0 {
		cmd.ExtraFields["managedFields"] = kubeMeta.ManagedFields
	}
	// deletionTimestamp can be nil, but is of a kubernetes time type, so we have to cast if non-nil
	if kubeMeta.DeletionTimestamp != nil {
		cmd.DeletionTimestamp = &kubeMeta.DeletionTimestamp.Time
	}
	// Other common metadata keys are in annotations
	cmd.CreatedBy = kubeMeta.Annotations[annotationPrefix+"createdBy"]
	cmd.UpdatedBy = kubeMeta.Annotations[annotationPrefix+"updatedBy"]
	cmd.UpdateTimestamp, _ = time.Parse(time.RFC3339, kubeMeta.Annotations[annotationPrefix+"updateTimestamp"])

	// TODO deletions commented for now, but the question about whether to leave them is kinda fundamental to converting between kubernetes and grafana shapes
	//delete(kubeMeta.Annotations, annotationPrefix+"updatedBy")
	//delete(kubeMeta.Annotations, annotationPrefix+"createdBy")
	//delete(kubeMeta.Annotations, annotationPrefix+"updateTimestamp")

	// For all other annotation keys which begin with annotationPrefix (grafana.com/), strip the prefix and put them in custom metadata
	customMeta := make(map[string]any)
	for key, val := range kubeMeta.Annotations {
		if len(key) <= len(annotationPrefix) || key[:len(annotationPrefix)] != annotationPrefix {
			// Keep going if the key doesn't start with annotationPrefix
			continue
		}
		tkey := key[len(annotationPrefix):]
		if _, ok := commonMetadataKeys[tkey]; ok {
			// We've already handled this one
			continue
		}
		//delete(kubeMeta.Annotations, key)
		customMeta[tkey] = val
	}
	// With annotations keys trimmed out of the original, we can add it to extra fields in common metadata
	cmd.ExtraFields = map[string]any{
		"generation":  kubeMeta.Generation,
		"annotations": kubeMeta.Annotations,
	}

	res.Metadata, err = json.Marshal(cmd)
	if err != nil {
		return res, err
	}
	res.CustomMetadata, err = json.Marshal(customMeta)
	return res, err
}

package kindsys

var _ Resource = &UnstructuredResource{}

// UnstructuredResource is an untyped representation of [Resource]. In the same
// way that map[string]any can represent any JSON []byte, UnstructuredResource
// can represent a [Resource] for any [Core] or [Custom] kind. But it is not
// strongly typed, and lacks any user-defined methods that may exist on a
// kind-specific struct that implements [Resource].
type UnstructuredResource struct {
	BasicMetadataObject
	Spec   map[string]any `json:"spec,omitempty"`
	Status map[string]any `json:"status,omitempty"`
}

func (u *UnstructuredResource) SpecObject() any {
	return u.Spec
}

func (u *UnstructuredResource) Subresources() map[string]any {
	return map[string]any{
		"status": u.Status,
	}
}

func (u *UnstructuredResource) Copy() Resource {
	com := CommonMetadata{
		UID:               u.CommonMeta.UID,
		ResourceVersion:   u.CommonMeta.ResourceVersion,
		CreationTimestamp: u.CommonMeta.CreationTimestamp.UTC(),
		UpdateTimestamp:   u.CommonMeta.UpdateTimestamp.UTC(),
		CreatedBy:         u.CommonMeta.CreatedBy,
		UpdatedBy:         u.CommonMeta.UpdatedBy,
	}

	copy(u.CommonMeta.Finalizers, com.Finalizers)
	if u.CommonMeta.DeletionTimestamp != nil {
		*com.DeletionTimestamp = *(u.CommonMeta.DeletionTimestamp)
	}
	for k, v := range u.CommonMeta.Labels {
		com.Labels[k] = v
	}
	com.ExtraFields = mapcopy(u.CommonMeta.ExtraFields)

	cp := UnstructuredResource{
		Spec:   mapcopy(u.Spec),
		Status: mapcopy(u.Status),
	}

	cp.CommonMeta = com
	cp.CustomMeta = mapcopy(u.CustomMeta)
	return &cp
}

func mapcopy(m map[string]any) map[string]any {
	cp := make(map[string]any)
	for k, v := range m {
		if vm, ok := v.(map[string]any); ok {
			cp[k] = mapcopy(vm)
		} else {
			cp[k] = v
		}
	}

	return cp
}

package packages

import (
	"fmt"
)

type Manifests []Manifest

type NormalizedPackage struct {
	Metadata  PackageMetadata `json:"metadata"`
	Status    PackageStatus   `json:"status"`
	Manifests Manifests       `json:"manifests"`
}

func (pkg NormalizedPackage) Descriptor() PackageDescriptor {
	return PackageDescriptor{
		Metadata: pkg.Metadata,
		Status:   pkg.Status,
	}
}

func (pkg NormalizedPackage) Locate(ref Reference) (Manifest, error) {
	return pkg.Manifests.Locate(ref)
}

type Manifest struct {
	APIVersion  string                 `json:"apiVersion"`
	Kind        string                 `json:"kind"`
	Metadata    map[string]string      `json:"metadata"`
	Annotations map[string]string      `json:"annotations"`
	Spec        map[string]interface{} `json:"spec"`
}

func (manifests Manifests) Locate(ref Reference) (Manifest, error) {
	for _, manifest := range manifests {
		if manifest.Kind == ref.ReferredKind && manifest.Metadata["name"] == ref.ReferredName {
			return manifest, nil
		}
	}
	return Manifest{}, fmt.Errorf("'%s' not found", ref.String())
}

type Reference struct {
	ReferredKind string `json:"referredKind"`
	ReferredName string `json:"referredName"`
}

func (ref Reference) Valid() bool {
	return ref.ReferredKind != "" && ref.ReferredName != ""
}

func (ref Reference) String() string {
	return fmt.Sprintf("%s:%s", ref.ReferredKind, ref.ReferredName)
}

type RefManifest struct {
	APIVersion  string            `json:"apiVersion"`
	Kind        string            `json:"kind"`
	Metadata    map[string]string `json:"metadata"`
	Annotations map[string]string `json:"annotations"`
	Spec        Reference         `json:"spec"`
}

func PanelRef(manifestName string) Reference {
	return Reference{
		ReferredKind: "Panel",
		ReferredName: manifestName,
	}
}

func TargetRef(manifestName string) Reference {
	return Reference{
		ReferredKind: "Target",
		ReferredName: manifestName,
	}
}

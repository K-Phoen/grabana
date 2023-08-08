package packages

type PackageMetadata struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Tags        []string `json:"tags"`
	Icon        string   `json:"icon"`

	Homepage string `json:"homepage"`
}

type PackageStatus struct {
	Installed bool `json:"installed"`
	Unpacked  bool `json:"unpacked"`
}

type Resource struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Spec        map[string]interface{} `json:"spec"`
}

type PackageSpec struct {
	Dashboards []Resource `json:"dashboards"`
	Panels     []Resource `json:"panels"`
	Queries    []Resource `json:"queries"`
	Alerts     []Resource `json:"alerts"`
}

type Package struct {
	Metadata PackageMetadata `json:"metadata"`
	Status   PackageStatus   `json:"status"`
	Spec     PackageSpec     `json:"spec"`
}

type PackageDescriptor struct {
	Metadata PackageMetadata `json:"metadata"`
	Status   PackageStatus   `json:"status"`
}

func (pkg Package) Descriptor() PackageDescriptor {
	return PackageDescriptor{
		Metadata: pkg.Metadata,
		Status:   pkg.Status,
	}
}

type Descriptor interface {
	Descriptor() PackageDescriptor
}

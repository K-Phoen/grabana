package packages

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

func isResource(input interface{}) bool {
	asMap := input.(map[string]interface{})

	return asMap["spec"] != nil
}

func DenormalizePanel(panel Manifest, normalized NormalizedPackage) (Resource, error) {
	denormalized := Resource{
		Title:       guessManifestTitle(panel),
		Description: guessManifestDescription(panel),
	}

	if err := mapstructure.Decode(panel.Spec, &denormalized.Spec); err != nil {
		return denormalized, err
	}

	if panel.Spec["targets"] == nil {
		return denormalized, nil
	}

	targets := panel.Spec["targets"].([]interface{})
	newTargets := make([]map[string]interface{}, 0, len(targets))
	for _, target := range targets {
		if !isResource(target) {
			newTargets = append(newTargets, target.(map[string]interface{}))
			continue
		}

		ref, err := referenceFromAny(target)
		if err != nil {
			return denormalized, err
		}

		if !ref.Valid() {
			continue
		}

		if ref.ReferredKind != "Target" {
			return denormalized, fmt.Errorf("panel '%s' references a '%s' within its 'targets' field", denormalized.Title, ref.ReferredKind)
		}

		referredTarget, err := normalized.Locate(ref)
		if err != nil {
			return denormalized, err
		}

		denormalizedTarget, err := DenormalizeQuery(referredTarget, normalized)
		if err != nil {
			return denormalized, err
		}

		newTargets = append(newTargets, denormalizedTarget.Spec)
	}

	denormalized.Spec["targets"] = newTargets

	return denormalized, nil
}

func DenormalizeQuery(query Manifest, _ NormalizedPackage) (Resource, error) {
	newQuery := Resource{
		Title:       guessManifestTitle(query),
		Description: guessManifestDescription(query),
	}

	if err := mapstructure.Decode(query.Spec, &newQuery.Spec); err != nil {
		return newQuery, err
	}

	return newQuery, nil
}

func referenceFromAny(refResource interface{}) (Reference, error) {
	resource := RefManifest{}
	err := mapstructure.Decode(refResource, &resource)

	return resource.Spec, err
}

func guessManifestTitle(manifest Manifest) string {
	if title, ok := manifest.Spec["title"]; ok {
		return title.(string)
	}

	return manifest.Metadata["name"]
}

func guessManifestDescription(manifest Manifest) string {
	if desc, ok := manifest.Annotations["packages.grafana.com/docs"]; ok {
		return desc
	}

	if desc, ok := manifest.Spec["description"]; ok {
		return desc.(string)
	}

	return manifest.Metadata["description"]
}

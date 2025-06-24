package functions

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// MergeManifests merges two lists of manifest objects by overwriting fields from the second list onto the first when apiVersion, kind, and metadata.name match.
func MergeManifests(_ context.Context, manifests1, manifests2 []map[string]interface{}) ([]map[string]interface{}, diag.Diagnostics) {
	merged := make([]map[string]interface{}, 0, len(manifests1)+len(manifests2))
	index := make(map[string]int)

	// Helper to create a unique key for each manifest
	manifestKey := func(m map[string]interface{}) string {
		apiVersion, _ := m["apiVersion"].(string)
		kind, _ := m["kind"].(string)
		metadata, _ := m["metadata"].(map[string]interface{})
		name := ""
		if metadata != nil {
			name, _ = metadata["name"].(string)
		}
		return fmt.Sprintf("%s|%s|%s", apiVersion, kind, name)
	}

	for i, m := range manifests1 {
		merged = append(merged, m)
		index[manifestKey(m)] = i
	}

	for _, m2 := range manifests2 {
		key := manifestKey(m2)
		if i, found := index[key]; found {
			// Overwrite fields from m2 onto m1 (shallow merge)
			for k, v := range m2 {
				merged[i][k] = v
			}
		} else {
			merged = append(merged, m2)
		}
	}

	return merged, nil
}

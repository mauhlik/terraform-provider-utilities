package functions

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func MergeManifests(_ context.Context, manifests1, manifests2 []map[string]interface{}) ([]map[string]interface{}, diag.Diagnostics) {
	merged := make([]map[string]interface{}, 0, len(manifests1)+len(manifests2))
	index := make(map[string]int)

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
		merged = append(merged, DeepCopyMap(m))
		index[manifestKey(m)] = i
	}

	for _, m2 := range manifests2 {
		key := manifestKey(m2)
		if i, found := index[key]; found {
			merged[i] = DeepMerge(merged[i], m2)
		} else {
			merged = append(merged, DeepCopyMap(m2))
		}
	}

	return merged, nil
}

// DeepMerge recursively merges src into dst.
// - For maps: merges keys recursively.
// - For slices: replaces the slice (Kubernetes semantics).
// - For other values: src overwrites dst.
func DeepMerge(dst, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		if vMap, ok := v.(map[string]interface{}); ok {
			if dstMap, found := dst[k].(map[string]interface{}); found {
				dst[k] = DeepMerge(dstMap, vMap)
			} else {
				dst[k] = DeepCopyMap(vMap)
			}
		} else if vSlice, ok := v.([]interface{}); ok {
			// Kubernetes usually replaces slices, not merges them
			dst[k] = DeepCopySlice(vSlice)
		} else {
			dst[k] = v
		}
	}
	return dst
}

// DeepCopyMap creates a deep copy of a map[string]interface{}.
func DeepCopyMap(src map[string]interface{}) map[string]interface{} {
	copy := make(map[string]interface{}, len(src))
	for k, v := range src {
		switch val := v.(type) {
		case map[string]interface{}:
			copy[k] = DeepCopyMap(val)
		case []interface{}:
			copy[k] = DeepCopySlice(val)
		default:
			copy[k] = val
		}
	}
	return copy
}

// DeepCopySlice creates a deep copy of a []interface{}.
func DeepCopySlice(src []interface{}) []interface{} {
	copy := make([]interface{}, len(src))
	for i, v := range src {
		switch val := v.(type) {
		case map[string]interface{}:
			copy[i] = DeepCopyMap(val)
		case []interface{}:
			copy[i] = DeepCopySlice(val)
		default:
			copy[i] = val
		}
	}
	return copy
}

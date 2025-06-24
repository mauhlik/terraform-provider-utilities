package functions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func NewMergeManifests() function.Function {
	return &mergeManifestsFunction{}
}

type mergeManifestsFunction struct{}

func (f *mergeManifestsFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "merge_manifests"
}

func (f *mergeManifestsFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Merge two lists/tuples of Kubernetes manifests",
		Description:         "Merges two lists/tuples of Kubernetes manifests. Objects are merged when apiVersion, kind, and metadata.name match.",
		MarkdownDescription: "Merges two lists/tuples of Kubernetes manifests. Objects are merged when apiVersion, kind, and metadata.name match.",
		Parameters: []function.Parameter{
			function.DynamicParameter{
				Name:        "manifests1",
				Description: "First list/tuple of manifest objects.",
			},
			function.DynamicParameter{
				Name:        "manifests2",
				Description: "Second list/tuple of manifest objects.",
			},
		},
		Return: function.DynamicReturn{},
	}
}

func (f *mergeManifestsFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var dv1, dv2 types.Dynamic
	if err := req.Arguments.GetArgument(ctx, 0, &dv1); err != nil {
		resp.Error = function.NewFuncError(fmt.Sprintf("failed to get manifests1: %v", err))
		return
	}
	if err := req.Arguments.GetArgument(ctx, 1, &dv2); err != nil {
		resp.Error = function.NewFuncError(fmt.Sprintf("failed to get manifests2: %v", err))
		return
	}

	// Check for null/unknown
	if dv1.IsNull() || dv1.IsUnknown() || dv2.IsNull() || dv2.IsUnknown() {
		resp.Error = function.NewFuncError("input manifests cannot be null or unknown")
		return
	}

	// Extract underlying value
	uv1 := dv1.UnderlyingValue()
	uv2 := dv2.UnderlyingValue()

	// Handle both tuples and lists
	manifests1, err := extractManifests(uv1)
	if err != nil {
		resp.Error = function.NewFuncError(fmt.Sprintf("manifests1: %v", err))
		return
	}
	manifests2, err := extractManifests(uv2)
	if err != nil {
		resp.Error = function.NewFuncError(fmt.Sprintf("manifests2: %v", err))
		return
	}

	// Merge manifests
	merged, _ := MergeManifests(ctx, manifests1, manifests2)

	// Convert merged manifests to JSON strings
	mergedElems := make([]attr.Value, len(merged))
	for i, manifest := range merged {
		jsonBytes, err := json.Marshal(manifest)
		if err != nil {
			resp.Error = function.NewFuncError(fmt.Sprintf("failed to marshal manifest: %v", err))
			return
		}
		mergedElems[i] = types.StringValue(string(jsonBytes))
	}

	// Create list of JSON strings
	mergedList := types.ListValueMust(types.StringType, mergedElems)

	// Return as dynamic value
	dynamicResult := types.DynamicValue(mergedList)
	if err := resp.Result.Set(ctx, dynamicResult); err != nil {
		resp.Error = err
	}
}

// extractManifests handles both lists and tuples
func extractManifests(val attr.Value) ([]map[string]interface{}, error) {
	switch v := val.(type) {
	case basetypes.ListValue:
		return extractFromList(v)
	case basetypes.TupleValue:
		return extractFromTuple(v)
	default:
		return nil, fmt.Errorf("expected list or tuple, got %T", val)
	}
}

func extractFromList(list basetypes.ListValue) ([]map[string]interface{}, error) {
	elements := list.Elements()
	result := make([]map[string]interface{}, len(elements))
	for i, elem := range elements {
		obj, ok := elem.(basetypes.ObjectValue)
		if !ok {
			return nil, fmt.Errorf("element %d is not an object", i)
		}
		result[i] = attrObjectToGo(obj)
	}
	return result, nil
}

func extractFromTuple(tuple basetypes.TupleValue) ([]map[string]interface{}, error) {
	elements := tuple.Elements()
	result := make([]map[string]interface{}, len(elements))
	for i, elem := range elements {
		obj, ok := elem.(basetypes.ObjectValue)
		if !ok {
			return nil, fmt.Errorf("tuple element %d is not an object", i)
		}
		result[i] = attrObjectToGo(obj)
	}
	return result, nil
}

// Helper: recursively convert attr.ObjectValue to Go map
func attrObjectToGo(obj basetypes.ObjectValue) map[string]interface{} {
	goMap := make(map[string]interface{}, len(obj.Attributes()))
	for k, v := range obj.Attributes() {
		goMap[k] = attrToGo(v)
	}
	return goMap
}

// Helper: recursively convert attr.Value to Go value
func attrToGo(val attr.Value) interface{} {
	switch v := val.(type) {
	case basetypes.StringValue:
		return v.ValueString()
	case basetypes.Int64Value:
		return v.ValueInt64()
	case basetypes.BoolValue:
		return v.ValueBool()
	case basetypes.Float64Value:
		return v.ValueFloat64()
	case basetypes.ListValue:
		result := make([]interface{}, len(v.Elements()))
		for i, e := range v.Elements() {
			result[i] = attrToGo(e)
		}
		return result
	case basetypes.ObjectValue:
		return attrObjectToGo(v)
	case basetypes.MapValue:
		result := make(map[string]interface{}, len(v.Elements()))
		for k, a := range v.Elements() {
			result[k] = attrToGo(a)
		}
		return result
	case basetypes.TupleValue:
		result := make([]interface{}, len(v.Elements()))
		for i, e := range v.Elements() {
			result[i] = attrToGo(e)
		}
		return result
	default:
		return nil
	}
}

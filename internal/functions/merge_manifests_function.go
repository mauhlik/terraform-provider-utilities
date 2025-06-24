package functions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

// NewMergeManifests returns a new instance of the merge_manifests function
func NewMergeManifests() function.Function {
	return mergeManifestsFunction{}
}

type mergeManifestsFunction struct{}

func (f mergeManifestsFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "merge_manifests"
}

func (f mergeManifestsFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Merge two lists of Kubernetes manifests, overwriting fields from the second list when apiVersion, kind, and metadata.name match.",
		Description:         "Merges two lists of Kubernetes manifests. If apiVersion, kind, and metadata.name match, fields from the second list overwrite the first.",
		MarkdownDescription: "Merges two lists of Kubernetes manifests. If apiVersion, kind, and metadata.name match, fields from the second list overwrite the first.",
		Parameters: []function.Parameter{
			function.DynamicParameter{
				Name:        "manifests1",
				Description: "First list of manifest objects.",
			},
			function.DynamicParameter{
				Name:        "manifests2",
				Description: "Second list of manifest objects.",
			},
		},
		Return: function.DynamicReturn{},
	}
}

func (f mergeManifestsFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var manifests1, manifests2 []map[string]interface{}
	if err := req.Arguments.GetArgument(ctx, 0, &manifests1); err != nil {
		resp.Error = err
		return
	}
	if err := req.Arguments.GetArgument(ctx, 1, &manifests2); err != nil {
		resp.Error = err
		return
	}
	merged, _ := MergeManifests(ctx, manifests1, manifests2)
	if err := resp.Result.Set(ctx, merged); err != nil {
		resp.Error = err
	}
}

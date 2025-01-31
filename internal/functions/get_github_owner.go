package functions

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = GetGithubOwnerVariable{}

func NewGetGithubOwnerVariable() function.Function {
	return GetGithubOwnerVariable{}
}

type GetGithubOwnerVariable struct{}

func (r GetGithubOwnerVariable) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "get_github_owner"
}

func (r GetGithubOwnerVariable) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Get repository owner",
		Description:         "Get the repository owner.",
		MarkdownDescription: "Get the repository owner.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "name",
				Description:         "The name of an repository",
				MarkdownDescription: "The name of an repository",
			},
		},
		Return: function.StringReturn{},
	}
}

func (r GetGithubOwnerVariable) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var repoName string

	resp.Error = req.Arguments.Get(ctx, &repoName)
	if resp.Error != nil {
		return
	}

	parts := strings.Split(repoName, "/")
	if len(parts) != 2 {
		resp.Error = function.NewFuncError("Invalid repository name, expected format is owner/repo")
		return
	}

	owner := parts[0]

	resp.Error = resp.Result.Set(ctx, owner)
}

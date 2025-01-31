package functions

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = GetGithubRepoNameVariable{}

func NewGetGithubRepoNameVariable() function.Function {
	return GetGithubRepoNameVariable{}
}

type GetGithubRepoNameVariable struct{}

func (r GetGithubRepoNameVariable) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "get_github_repo_name"
}

func (r GetGithubRepoNameVariable) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Get repository name",
		Description:         "Get the repository name.",
		MarkdownDescription: "Get the repository name.",
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

func (r GetGithubRepoNameVariable) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
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

	repo := parts[1]

	resp.Error = resp.Result.Set(ctx, repo)
}

package provider

import (
	"context"

	"github.com/mauhlik/terraform-provider-utilities/internal/functions"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = &UtilitiesFunctionsProvider{}

type UtilitiesFunctionsProvider struct{}

func NewUtilitiesFunctionsProvider() func() provider.Provider {
	return func() provider.Provider {
		return &UtilitiesFunctionsProvider{}
	}
}

func (p *UtilitiesFunctionsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "utilities_functions"
}

func (p *UtilitiesFunctionsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *UtilitiesFunctionsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *UtilitiesFunctionsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *UtilitiesFunctionsProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

func (p *UtilitiesFunctionsProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		functions.NewGetEnvironmentVariable,
		functions.NewGetGithubOwnerVariable,
		functions.NewGetGithubRepoNameVariable,
		functions.NewDelayValue,
		functions.NewMergeManifests,
	}
}

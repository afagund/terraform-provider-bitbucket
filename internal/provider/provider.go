package provider

import (
	"context"
	"os"

	"github.com/afagund/terraform-provider-bitbucket/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ provider.Provider = &bitbucketProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &bitbucketProvider{
			version: version,
		}
	}
}

type bitbucketProvider struct {
	version string
}

type bitbucketProviderModel struct {
	Host      types.String `tfsdk:"host"`
	Workspace types.String `tfsdk:"workspace"`
	Token     types.String `tfsdk:"token"`
}

func (p *bitbucketProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "bitbucket"
	resp.Version = p.version
}

func (p *bitbucketProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"workspace": schema.StringAttribute{
				Optional: true,
			},
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *bitbucketProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Bitbucket client")

	var config bitbucketProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Bitbucket API Host",
			"The provider cannot create the Bitbucket API client as there is an unknown configuration value for the Bitbucket API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BITBUCKET_HOST environment variable.",
		)
	}

	if config.Workspace.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("workspace"),
			"Unknown Bitbucket API Workspace",
			"The provider cannot create the Bitbucket API client as there is an unknown configuration value for the Bitbucket API workspace. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BITBUCKET_WORKSPACE environment variable.",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Bitbucket API Token",
			"The provider cannot create the Bitbucket API client as there is an unknown configuration value for the Bitbucket API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BITBUCKET_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("BITBUCKET_HOST")
	workspace := os.Getenv("BITBUCKET_WORKSPACE")
	token := os.Getenv("BITBUCKET_TOKEN")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Workspace.IsNull() {
		workspace = config.Workspace.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Bitbucket API Host",
			"The provider cannot create the Bitbucket API client as there is a missing or empty value for the Bitbucket API host. "+
				"Set the host value in the configuration or use the BITBUCKET_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if workspace == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("workspace"),
			"Missing Bitbucket API Workspace",
			"The provider cannot create the Bitbucket API client as there is a missing or empty value for the Bitbucket API workspace. "+
				"Set the workspace value in the configuration or use the BITBUCKET_WORKSPACE environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Bitbucket API Token",
			"The provider cannot create the Bitbucket API client as there is a missing or empty value for the Bitbucket API token. "+
				"Set the token value in the configuration or use the BITBUCKET_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "bitbucket_host", host)
	ctx = tflog.SetField(ctx, "bitbucket_workspace", workspace)
	ctx = tflog.SetField(ctx, "bitbucket_token", token)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "bitbucket_token")

	tflog.Debug(ctx, "Creating Bitbucket client")

	client, err := client.NewClient(&host, &workspace, &token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Bitbucket API Client",
			"An unexpected error occurred when creating the Bitbucket API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Bitbucket Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Bitbucket client", map[string]any{"success": true})
}

func (p *bitbucketProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewRepositoryDataSource,
	}
}

func (p *bitbucketProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRepositoryResource,
		NewGroupPermissionResource,
		NewBranchRestrictionResource,
	}
}

package provider

import (
	"context"
	"fmt"

	"github.com/afagund/terraform-provider-bitbucket/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &repositoryDataSource{}
	_ datasource.DataSourceWithConfigure = &repositoryDataSource{}
)

type repositoryDataSourceModel struct {
	Slug      types.String  `tfsdk:"slug"`
	IsPrivate types.Bool    `tfsdk:"is_private"`
	Scm       types.String  `tfsdk:"scm"`
	Project   *projectModel `tfsdk:"project"`
	Website   types.String  `tfsdk:"website"`
}

func NewRepositoryDataSource() datasource.DataSource {
	return &repositoryDataSource{}
}

type repositoryDataSource struct {
	client *client.Client
}

func (d *repositoryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository"
}

func (d *repositoryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"slug": schema.StringAttribute{
				Required: true,
			},
			"is_private": schema.BoolAttribute{
				Computed: true,
			},
			"scm": schema.StringAttribute{
				Computed: true,
			},
			"project": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"key": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"website": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (d *repositoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state repositoryDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	repository, err := d.client.GetRepository(state.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Bitbucket Repository",
			"Could not read Bitbucket repository "+state.Slug.ValueString()+": "+err.Error(),
		)
		return
	}

	state.mapFrom(repository)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *repositoryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (m *repositoryDataSourceModel) mapFrom(c *client.Repository) {
	var project projectModel
	project.Key = types.StringValue(c.Project.Key)

	m.Slug = types.StringValue(c.Slug)
	m.IsPrivate = types.BoolValue(c.IsPrivate)
	m.Scm = types.StringValue(c.Scm)
	m.Project = &project
	m.Website = types.StringPointerValue(c.Website)
}

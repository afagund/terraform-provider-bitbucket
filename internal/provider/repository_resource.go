package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/afagund/terraform-provider-bitbucket/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &repositoryResource{}
	_ resource.ResourceWithConfigure   = &repositoryResource{}
	_ resource.ResourceWithImportState = &repositoryResource{}
)

type projectModel struct {
	Key types.String `tfsdk:"key"`
}

type repositoryResourceModel struct {
	Slug      types.String  `tfsdk:"slug"`
	IsPrivate types.Bool    `tfsdk:"is_private"`
	Scm       types.String  `tfsdk:"scm"`
	Project   *projectModel `tfsdk:"project"`
	Website   types.String  `tfsdk:"website"`
}

func NewRepositoryResource() resource.Resource {
	return &repositoryResource{}
}

type repositoryResource struct {
	client *client.Client
}

func (r *repositoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository"
}

func (r *repositoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"slug": schema.StringAttribute{
				Required: true,
			},
			"is_private": schema.BoolAttribute{
				Required: true,
			},
			"scm": schema.StringAttribute{
				Required: true,
			},
			"project": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"key": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"website": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (r *repositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan repositoryResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var newRepository client.Repository
	plan.mapTo(&newRepository)

	repository, err := r.client.CreateRepository(plan.Slug.ValueString(), newRepository)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating repository",
			"Could not create repository, unexpected error: "+err.Error(),
		)
		return
	}

	plan.mapFrom(repository)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *repositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state repositoryResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	repository, err := r.client.GetRepository(state.Slug.ValueString())
	if err != nil {
		var statusErr client.StatusErr
		if errors.As(err, &statusErr) {
			if statusErr.StatusCode == 404 {
				resp.State.RemoveResource(ctx)
				return
			}
		}

		resp.Diagnostics.AddError(
			"Error Reading Bitbucket Repository",
			"Could not read Bitbucket repository "+state.Slug.ValueString()+": "+err.Error(),
		)
		return
	}

	state.mapFrom(repository)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *repositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan repositoryResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var newRepository client.Repository
	plan.mapTo(&newRepository)

	repository, err := r.client.UpdateRepository(plan.Slug.ValueString(), newRepository)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Bitbucket Repository",
			"Could not update repository, unexpected error: "+err.Error(),
		)
		return
	}

	plan.mapFrom(repository)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *repositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state repositoryResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRepository(state.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Bitbucket Repository",
			"Could not delete repository, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *repositoryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

func (r *repositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("slug"), req, resp)
}

func (m *repositoryResourceModel) mapTo(c *client.Repository) {
	c.IsPrivate = m.IsPrivate.ValueBool()
	c.Scm = m.Scm.ValueString()
	c.Project.Key = m.Project.Key.ValueString()
	c.Website = m.Website.ValueStringPointer()
}

func (m *repositoryResourceModel) mapFrom(c *client.Repository) {
	var project projectModel
	project.Key = types.StringValue(c.Project.Key)

	m.Slug = types.StringValue(c.Slug)
	m.IsPrivate = types.BoolValue(c.IsPrivate)
	m.Scm = types.StringValue(c.Scm)
	m.Project = &project
	m.Website = types.StringPointerValue(c.Website)
}

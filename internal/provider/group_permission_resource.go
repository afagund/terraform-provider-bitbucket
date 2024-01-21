package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/afagund/terraform-provider-bitbucket/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &groupPermissionResource{}
	_ resource.ResourceWithConfigure   = &groupPermissionResource{}
	_ resource.ResourceWithImportState = &groupPermissionResource{}
)

type groupPermissionResourceModel struct {
	RepositorySlug types.String `tfsdk:"repository_slug"`
	GroupSlug      types.String `tfsdk:"group_slug"`
	Permission     types.String `tfsdk:"permission"`
}

func NewGroupPermissionResource() resource.Resource {
	return &groupPermissionResource{}
}

type groupPermissionResource struct {
	client *client.Client
}

func (r *groupPermissionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_permission"
}

func (r *groupPermissionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"repository_slug": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group_slug": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"permission": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *groupPermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan groupPermissionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var newGroupPermission client.GroupPermission
	plan.mapTo(&newGroupPermission)

	groupPermission, err := r.client.CreateGroupPermission(plan.RepositorySlug.ValueString(), plan.GroupSlug.ValueString(), newGroupPermission)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating group permission",
			"Could not create group permission, unexpected error: "+err.Error(),
		)
		return
	}

	plan.mapFrom(groupPermission)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *groupPermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupPermissionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupPermission, err := r.client.GetGroupPermission(state.RepositorySlug.ValueString(), state.GroupSlug.ValueString())
	if err != nil {
		var statusErr client.StatusErr
		if errors.As(err, &statusErr) {
			if statusErr.StatusCode == 404 {
				resp.State.RemoveResource(ctx)
				return
			}
		}

		resp.Diagnostics.AddError(
			"Error Reading Bitbucket Group Permission",
			"Could not read Bitbucket group permission for repository "+state.RepositorySlug.ValueString()+": "+err.Error(),
		)
		return
	}

	state.mapFrom(groupPermission)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *groupPermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan groupPermissionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var newGroupPermission client.GroupPermission
	plan.mapTo(&newGroupPermission)

	groupPermission, err := r.client.UpdateGroupPermission(plan.RepositorySlug.ValueString(), plan.GroupSlug.ValueString(), newGroupPermission)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Bitbucket Group Permission",
			"Could not update group permission, unexpected error: "+err.Error(),
		)
		return
	}

	plan.mapFrom(groupPermission)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *groupPermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state groupPermissionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteGroupPermission(state.RepositorySlug.ValueString(), state.GroupSlug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Bitbucket Group Permission",
			"Could not delete group permission, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *groupPermissionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *groupPermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: repository_slug,group_slug. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("repository_slug"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_slug"), idParts[1])...)
}

func (m *groupPermissionResourceModel) mapTo(c *client.GroupPermission) {
	c.Permission = m.Permission.ValueString()
}

func (m *groupPermissionResourceModel) mapFrom(c *client.GroupPermission) {
	m.Permission = types.StringValue(c.Permission)
}

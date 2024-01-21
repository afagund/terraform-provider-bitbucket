package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/afagund/terraform-provider-bitbucket/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &branchRestrictionResource{}
	_ resource.ResourceWithConfigure   = &branchRestrictionResource{}
	_ resource.ResourceWithImportState = &branchRestrictionResource{}
)

type userModel struct {
	Uuid types.String `tfsdk:"uuid"`
}

type groupModel struct {
	Slug types.String `tfsdk:"slug"`
}

type branchRestrictionResourceModel struct {
	ID              types.Int64  `tfsdk:"id"`
	RepositorySlug  types.String `tfsdk:"repository_slug"`
	Kind            types.String `tfsdk:"kind"`
	BranchMatchKind types.String `tfsdk:"branch_match_kind"`
	Pattern         types.String `tfsdk:"pattern"`
	Users           []userModel  `tfsdk:"users"`
	Groups          []groupModel `tfsdk:"groups"`
}

func NewBranchRestrictionResource() resource.Resource {
	return &branchRestrictionResource{}
}

type branchRestrictionResource struct {
	client *client.Client
}

func (r *branchRestrictionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branch_restriction"
}

func (r *branchRestrictionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"repository_slug": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"kind": schema.StringAttribute{
				Required: true,
			},
			"branch_match_kind": schema.StringAttribute{
				Required: true,
			},
			"pattern": schema.StringAttribute{
				Required: true,
			},
			"users": schema.ListNestedAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
			"groups": schema.ListNestedAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"slug": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
		},
	}
}

func (r *branchRestrictionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan branchRestrictionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var newBranchRestriction client.BranchRestriction
	plan.mapTo(&newBranchRestriction)

	branchRestriction, err := r.client.CreateBranchRestriction(plan.RepositorySlug.ValueString(), newBranchRestriction)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating branch restriction",
			"Could not create branch restriction, unexpected error: "+err.Error(),
		)
		return
	}

	plan.mapFrom(branchRestriction)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *branchRestrictionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state branchRestrictionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	branchRestriction, err := r.client.GetBranchRestriction(state.RepositorySlug.ValueString(), int(state.ID.ValueInt64()))
	if err != nil {
		var statusErr client.StatusErr
		if errors.As(err, &statusErr) {
			if statusErr.StatusCode == 404 {
				resp.State.RemoveResource(ctx)
				return
			}
		}

		resp.Diagnostics.AddError(
			"Error Reading Bitbucket Branch Restriction",
			"Could not read Bitbucket branch restriction id "+state.ID.String()+": "+err.Error(),
		)
		return
	}

	state.mapFrom(branchRestriction)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *branchRestrictionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan branchRestrictionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var newBranchRestriction client.BranchRestriction
	plan.mapTo(&newBranchRestriction)

	branchRestriction, err := r.client.UpdateBranchRestriction(plan.RepositorySlug.ValueString(), int(plan.ID.ValueInt64()), newBranchRestriction)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Bitbucket Branch Restriction",
			"Could not update branch restriction, unexpected error: "+err.Error(),
		)
		return
	}

	plan.mapFrom(branchRestriction)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *branchRestrictionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state branchRestrictionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteBranchRestriction(state.RepositorySlug.ValueString(), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Bitbucket Branch Restriction",
			"Could not delete branch restriction, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *branchRestrictionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *branchRestrictionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: repository_slug,id. Got: %q", req.ID),
		)
		return
	}

	id, err := strconv.Atoi(idParts[1])
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting ID",
			"Could not convert resource id, unexpected error: "+err.Error(),
		)
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("repository_slug"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (m *branchRestrictionResourceModel) mapTo(c *client.BranchRestriction) {
	var users []client.User
	for _, user := range m.Users {
		users = append(users, client.User{
			Uuid: user.Uuid.ValueString(),
		})
	}

	var groups []client.Group
	for _, group := range m.Groups {
		groups = append(groups, client.Group{
			Slug: group.Slug.ValueString(),
		})
	}

	c.Kind = m.Kind.ValueString()
	c.BranchMatchKind = m.BranchMatchKind.ValueString()
	c.Pattern = m.Pattern.ValueString()
	c.Users = users
	c.Groups = groups
}

func (m *branchRestrictionResourceModel) mapFrom(c *client.BranchRestriction) {
	m.ID = types.Int64Value(int64(c.ID))
	m.Kind = types.StringValue(c.Kind)
	m.BranchMatchKind = types.StringValue(c.BranchMatchKind)
	m.Pattern = types.StringValue(c.Pattern)

	m.Users = []userModel{}
	for _, user := range c.Users {
		m.Users = append(m.Users, userModel{
			Uuid: types.StringValue(user.Uuid),
		})
	}

	m.Groups = []groupModel{}
	for _, group := range c.Groups {
		m.Groups = append(m.Groups, groupModel{
			Slug: types.StringValue(group.Slug),
		})
	}
}

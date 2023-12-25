package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/williamokano/litmus-chaos-thin-client/pkg/client"
)

var (
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

type userDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Username types.String `tfsdk:"username"`
	Role     types.String `tfsdk:"role"`
	Name     types.String `tfsdk:"name"`
	Email    types.String `tfsdk:"email"`
}

type userDataSource struct {
	client *client.LitmusClient
}

func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

func (d *userDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	litmusClient, ok := req.ProviderData.(*client.LitmusClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *litmus.LitmusClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = litmusClient
}

func (d *userDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *userDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "User ID",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "User username",
				Required:    true,
			},
			"email": schema.StringAttribute{
				Description: "User email",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "User name",
				Computed:    true,
			},
			"role": schema.StringAttribute{
				Description: "User role",
				Computed:    true,
			},
		},
	}
}

func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state userDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := d.client.FindUserByUsername(state.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Litmus Chaos User",
			"Could not read Litmus Chaos User with username "+state.Username.ValueString()+": "+err.Error(),
		)
		return
	}

	state.ID = types.StringValue(user.ID)
	state.Name = types.StringValue(user.Name)
	state.Role = types.StringValue(string(user.Role))

	if user.Name != "" {
		state.Name = types.StringValue(user.Name)
	}

	if user.Email != "" {
		state.Email = types.StringValue(user.Email)
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

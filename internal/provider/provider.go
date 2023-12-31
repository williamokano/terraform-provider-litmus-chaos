package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/williamokano/litmus-chaos-thin-client/pkg/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &litmusChaosProvider{}
)

// litmusChaosProviderModel maps provider schema data to a Go type.
type litmusChaosProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Token    types.String `tfsdk:"token"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &litmusChaosProvider{
			version: version,
		}
	}
}

// litmusChaosProvider is the provider implementation.
type litmusChaosProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *litmusChaosProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "litmus-chaos"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *litmusChaosProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Litmus Chaos Control Plane.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URI for Litmus Chaos Control Plane. May also be provided via `LITMUS_CHAOS_HOST` environment variable.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username for Litmus Chaos server. May also be provided via `LITMUS_CHAOS_USERNAME` environment variable.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for Litmus Chaos server. May also be provided via `LITMUS_CHAOS_PASSWORD` environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"token": schema.StringAttribute{
				Description: "API Token for Litmus Chaos Control Plane API. May also be provided via `LITMUS_CHAOS_TOKEN` environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a Litmus Chaos client for data sources and resources.
func (p *litmusChaosProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Litmus Chaos client")

	// Retrieve provider data from configuration
	var config litmusChaosProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Litmus Chaos Host",
			"The provider cannot create the Litmus Chaos client as there is an unknown configuration value for the Litmus Chaos Control Plane host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LITMUS_CHAOS_HOST environment variable.",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Litmus Chaos API Token",
			"The provider cannot create the Litmus Chaos API client as there is an unknown configuration value for the Litmus Chaos API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LITMUS_CHAOS_TOKEN environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Litmus Chaos API Username",
			"The provider cannot create the Litmus Chaos API client as there is an unknown configuration value for the Litmus Chaos API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LITMUS_CHAOS_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Litmus Chaos API Password",
			"The provider cannot create the Litmus Chaos API client as there is an unknown configuration value for the Litmus Chaos API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LITMUS_CHAOS_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	var (
		host     = os.Getenv("LITMUS_CHAOS_HOST")
		token    = os.Getenv("LITMUS_CHAOS_TOKEN")
		username = os.Getenv("LITMUS_CHAOS_USERNAME")
		password = os.Getenv("LITMUS_CHAOS_PASSWORD")
	)

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Litmus Chaos Control Plane Host",
			"The provider cannot create the Litmus Chaos API client as there is a missing or empty value for the Litmus Chaos API host. "+
				"Set the host value in the configuration or use the LITMUS_CHAOS_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	// If token is not provided, fallback to username and pass. Messages related to token
	// must be preferential as username/password is not secure
	if token == "" {
		// No fallback
		if username == "" && password == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("token"),
				"Missing Litmus Chaos API Token",
				"The provider cannot create the Litmus Chaos client as there is a missing or empty value for the Litmus Chaos API token. "+
					"Set the token value in the configuration or use the LITMUS_CHAOS_TOKEN environment variable. "+
					"If either is already set, ensure the value is not empty. "+
					"This can be empty if username and password are BOTH provided. ",
			)
		} else {

			// Missing one of the fields
			if username == "" {
				resp.Diagnostics.AddAttributeError(
					path.Root("username"),
					"Missing Litmus Chaos API Username",
					"The provider cannot create the Litmus Chaos client as there is a missing or empty value for the Litmus Chaos API Username. "+
						"Set the username value in the configuration or use the LITMUS_CHAOS_USERNAME environment variable. "+
						"If either is already set, ensure the value is not empty. "+
						"This can be empty if token is provided. ",
				)
			}

			if password == "" {
				resp.Diagnostics.AddAttributeError(
					path.Root("password"),
					"Missing Litmus Chaos API Password",
					"The provider cannot create the Litmus Chaos client as there is a missing or empty value for the Litmus Chaos API Password. "+
						"Set the password value in the configuration or use the LITMUS_CHAOS_PASSWORD environment variable. "+
						"If either is already set, ensure the value is not empty. "+
						"This can be empty if token is provided. ",
				)
			}
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "litmus_chaos_host", host)
	ctx = tflog.SetField(ctx, "litmus_chaos_token", token)
	ctx = tflog.SetField(ctx, "litmus_chaos_username", username)
	ctx = tflog.SetField(ctx, "litmus_chaos_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "litmus_chaos_token", "litmus_chaos_password")

	tflog.Debug(ctx, "Creating Litmus Chaos client")

	// Create a new Litmus Chaos client using the configuration values
	litmusClient, err := client.NewClientFromCredentials(host, client.LitmusCredentials{
		Username: username,
		Password: password,
		Token:    token,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Litmus Chaos Client",
			"An unexpected error occurred when creating the Litmus Chaos client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Litmus Chaos Client Error: "+err.Error(),
		)
		return
	}

	// Make the Litmus Chaos client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = litmusClient
	resp.ResourceData = litmusClient

	tflog.Info(ctx, "Configured Litmus Chaos client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *litmusChaosProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewUserDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *litmusChaosProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectResource,
	}
}

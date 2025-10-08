// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func NewPopulateAction() action.Action {
	return &populateAction{}
}

type populateAction struct {
	client *Client
}

var (
	_ action.Action              = (*populateAction)(nil)
	_ action.ActionWithConfigure = (*populateAction)(nil)
)

// Configure adds the provider configured client to the resource.
func (a *populateAction) Configure(_ context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *provider.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	a.client = client
}

type populateActionModel struct {
	DatabaseId types.String        `tfsdk:"id"`
	Items      []databaseItemModel `tfsdk:"items"`
}

// orderItemModel maps order item data.
type databaseItemModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func (a *populateAction) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_population_action"
}

func (a *populateAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Invokes an AWS Lambda function with the specified payload. This action allows for imperative invocation of Lambda functions with full control over invocation parameters.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The name, ARN, or partial ARN of the Lambda function to invoke. You can specify a function name (e.g., my-function), a qualified function name (e.g., my-function:PROD), or a partial ARN (e.g., 123456789012:function:my-function).",
				Required:    true,
			},
			"items": schema.ListNestedAttribute{
				Description: "The JSON payload to send to the Lambda function. This should be a valid JSON string that represents the event data for your function.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Required: true,
						},
						"value": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
		},
	}
}

func (a *populateAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var config populateActionModel
	tflog.Info(ctx, "Dumping client", map[string]any{
		"host_uri": string(a.client.HostURL),
		"token":    string(a.client.Token),
	})
	// Parse configuration
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Invoking Populator", map[string]any{
		"database_id": string(config.DatabaseId.ValueString()),
	})
	time.Sleep(10 * time.Second)
	for _, item := range config.Items {
		itemRequest := CreateDatabaseItemRequest{
			Key:   string(item.Key.ValueString()),
			Value: string(item.Value.ValueString()),
		}

		_, err := a.client.CreateDatabaseItem(itemRequest, string(config.DatabaseId.ValueString()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating order",
				"Could not create order, unexpected error: "+err.Error(),
			)
			return
		}
		time.Sleep(1 * time.Second)
		// tflog.Info(ctx, "Lambda function invocation action completed successfully", map[string]any{
		// 	"database_id": string(config.DatabaseId.ValueString()),
		// 	"key":         string(create_item_response.Key),
		// 	"value":       string(create_item_response.Value),
		// })
	}

}

package model

import "github.com/hashicorp/terraform-plugin-framework/types"

type PrivilegeApplication struct {
	Id          types.String   `tfsdk:"id"`
	Name        types.String   `tfsdk:"name"`
	Description types.String   `tfsdk:"description"`
	Actions     []types.String `tfsdk:"actions"`
	Domain      types.String   `tfsdk:"domain"`
}

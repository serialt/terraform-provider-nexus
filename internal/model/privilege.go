package model

import "github.com/hashicorp/terraform-plugin-framework/types"

type PrivilegeApplication struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Actions     types.List   `tfsdk:"actions"`
	Domain      types.String `tfsdk:"domain"`
}

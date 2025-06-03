package model

import "github.com/hashicorp/terraform-plugin-framework/types"

type SoftQuotaModel struct {
	Limit types.Int64  `tfsdk:"limit"`
	Type  types.String `tfsdk:"type"`
}

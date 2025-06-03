package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BlobStoreFileModel struct {
	Id                    types.String    `tfsdk:"id"`
	Name                  types.String    `tfsdk:"name"`
	Path                  types.String    `tfsdk:"path"`
	BlobCount             types.Int64     `tfsdk:"blob_count"`
	AvailableSpaceInBytes types.Int64     `tfsdk:"available_space_in_bytes"`
	TotalSizeInBytes      types.Int64     `tfsdk:"total_size_in_bytes"`
	SoftQuota             *SoftQuotaModel `tfsdk:"soft_quota"`
}

type BlobStoreGroupModel struct {
	Id                    types.String    `tfsdk:"id"`
	Name                  types.String    `tfsdk:"name"`
	AvailableSpaceInBytes types.Int64     `tfsdk:"available_space_in_bytes"`
	BlobCount             types.Int64     `tfsdk:"blob_count"`
	FillPolicy            types.String    `tfsdk:"fill_policy"`
	Members               []types.String  `tfsdk:"members"`
	TotalSizeInBytes      types.Int64     `tfsdk:"total_size_in_bytes"`
	SoftQuota             *SoftQuotaModel `tfsdk:"soft_quota"`
}

type ComponentModel struct {
	ProprietaryComponents types.Bool `tfsdk:"proprietary_components"`
}

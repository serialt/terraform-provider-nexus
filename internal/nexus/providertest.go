package nexus

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerConfig = `
provider "nexus" {
  insecure = true
  url      = "http://nexus.local.io:8081"
  username = "admin"
  password = "sugar"
`
)

var (
	TestNexusProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)
)

func init() {
	TestNexusProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"nexus": providerserver.NewProtocol6WithError(New()),
	}

}

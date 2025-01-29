package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/mauhlik/terraform-provider-utilities/internal/provider"
)

func main() {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/mauhlik/terraform-provider-utilities",
		Debug:   false,
	}

	err := providerserver.Serve(context.Background(), provider.NewUtilitiesFunctionsProvider(), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}

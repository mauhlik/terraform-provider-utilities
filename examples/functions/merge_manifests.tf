terraform {
  required_providers {
    utilities = {
      source = "mauhlik/utilities"
      version = ">= 0.1.0"
    }
  }
}

locals {
  manifests1 = [
    {
      apiVersion = "v1"
      kind       = "ConfigMap"
      metadata   = { name = "my-config" }
      data       = { foo = "bar" }
    }
  ]
  manifests2 = [
    {
      apiVersion = "v1"
      kind       = "ConfigMap"
      metadata   = { name = "my-config" }
      data       = { foo = "baz", extra = "value" }
    }
  ]
  merged = provider::utilities::merge_manifests(local.manifests1, local.manifests2)
}

output "merged" {
  value = local.merged
}

# Merge Manifests Function

## Description

The `merge_manifests` function merges two lists of Kubernetes manifest objects. If two manifests have the same `apiVersion`, `kind`, and `metadata.name`, the fields from the second list will overwrite those in the first (shallow merge).

## Usage

```
merge_manifests(manifests1, manifests2)
```

- `manifests1`: First list of manifest objects.
- `manifests2`: Second list of manifest objects.

Returns: The merged list of manifest objects.

## Example

```
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
```

Result:

```
[
  {
    apiVersion = "v1"
    kind       = "ConfigMap"
    metadata   = { name = "my-config" }
    data       = { foo = "baz", extra = "value" }
  }
]
```

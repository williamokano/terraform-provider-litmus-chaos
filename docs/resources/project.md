---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "litmus-chaos_project Resource - terraform-provider-litmus-chaos"
subcategory: ""
description: |-
  
---

# litmus-chaos_project (Resource)



## Example Usage

```terraform
# Create a new project
resource "litmus-chaos_project" "main_project" {
  name = "Main Project"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the Project

### Read-Only

- `id` (String) Project ID
- `last_updated` (String) Date of last modification

## Import

Import is supported using the following syntax:

```shell
# Project can be imported by specifying the uuid identifier.
terraform import litmus_chaos_project.main_project "96fc3e7d-33ad-4891-8be2-ae4a30ba76cf"
```

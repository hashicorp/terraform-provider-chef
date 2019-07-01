---
layout: "chef"
page_title: "Chef: chef_search"
sidebar_current: "docs-chef-data-source-search"
description: |-
  Searches for and returns data from Chef Server.
---

# Data Source: chef_search

Use this data source to make a [search](https://docs.chef.io/chef_search.html) and return data from
Chef Server.

## Example Usage

```hcl
data "chef_search" "node_with_role" {
    query = "chef_environment:test AND role:some_role"
    filter  {
        name = "a"
        value = ["a_cookbook", "a"]
    }
    filter {
        name = "b"
        value = ["a_cookbook", "b"]
    }
    filter {
        name = "host"
        value = ["ipaddress"]
    }
}

data "chef_search" "some_role" {
    index = "role"
    query = "name:some_role"
    filter  {
        name = "a"
        value = ["override_attributes", "a_cookbook", "a"]
    }
    filter {
        name = "b"
        value = ["override_attributes", "a_cookbook", "b"]
    }
    unique = true
}

data "chef_search" "some_check" {
    index = "nagios_services"
    query = "id:some_check"
    unique = true
}

output "node_with_role_a" {
    value = "${data.chef_search.node_with_role.result.a}"
}

output "node_with_role_host" {
    value = "${data.chef_search.some_role.result.host}"
}

output "some_role_b" {
    value = "${data.chef_search.some_role.result.b}"
}

output "command" {
    value = "${data.chef_search.some_check.result.command_line}"
}
```

## Argument reference

The following arguments are supported:

* `ìndex` - (Optional, Defaults to `node`) The name of the index on the Chef server against which
  the search query will run: `client`, *`data_bag_name`*, `environment`, `node` or `role`.

* `query` - (Required) A valid search query against an object on the Chef server.

* `filter` - (Optional) One or more name/value pairs that specify what will be returned in the
  result. The value is a list indicating a path to the attribute to return. The value is a Note that
  every value in the result needs to be a string. Non strings will be converted to strings. This is
  due to limitations in terraform.

* `unique` - (Optional) If `true` an error will be raised if query didn't result in exactly one result.

## Attributes Reference

The following attributes are exported:

* `total_num` - The number of results that matched the query.

* `result` - If the query gave any result the first one will be put in result.

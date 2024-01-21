terraform {
  required_providers {
    bitbucket = {
      source = "hashicorp.com/afagund/bitbucket"
    }
  }
}

provider "bitbucket" {
  host      = "https://api.bitbucket.org/2.0"
  workspace = "afagund"
}

data "bitbucket_repository" "this" {
  slug = "demo"
}

output "bitbucket_repository" {
  value = data.bitbucket_repository.this
}

resource "bitbucket_repository" "demo" {
  slug       = "demo"
  is_private = true
  scm        = "git"

  project = {
    key : "INT"
  }
}

resource "bitbucket_group_permission" "admins" {
  repository_slug = "demo"
  group_slug      = "admins"
  permission      = "write"

  depends_on = [
    bitbucket_repository.demo
  ]
}

resource "bitbucket_group_permission" "users" {
  repository_slug = "demo"
  group_slug      = "users"
  permission      = "write"

  depends_on = [
    bitbucket_repository.demo
  ]
}

resource "bitbucket_branch_restriction" "demo" {
  repository_slug   = "demo"
  kind              = "push"
  branch_match_kind = "glob"
  pattern           = "main"

  users = [
    {
      uuid = "{115d8671-697c-4a7e-8958-c948613e3a79}"
    }
  ]

  groups = [
    {
      slug = "admins"
    },
    {
      slug = "users"
    }
  ]

  depends_on = [
    bitbucket_repository.demo
  ]
}

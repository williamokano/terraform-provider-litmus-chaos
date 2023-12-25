# Configuration-based authentication

terraform {
  required_providers {
    litmus-chaos = {
      source = "williamokano/litmus-chaos"
    }
  }
}

provider "litmus-chaos" {

}

# Create a new project
data "litmus-chaos_user" "user_foo" {
  username = "william.okano@deliveryhero.com"
}

output "user_foo_id" {
  value = data.litmus-chaos_user.user_foo.id
}
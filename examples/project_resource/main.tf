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
resource "litmus-chaos_project" "main_project" {
  name = "Main Project"
}

# Configuration-based authentication
provider "litmus-chaos" {
  host  = "https://litmus-chaos.control.plane.url"
  token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZXh0Ijoibm90aGluZyB0byBzZWUgaGVyZSJ9.SRUjK3EShU1vwJ3kokJEez25GmmzuU1-NF2iMDXBh8c"
}

# Configuration example based on username and password
# never use admin user for automation, always prefer token
provider "litmus-chaos" {
  host     = "https://litmus-chaos.control.plane.url"
  username = "admin"
  password = "litmus"
  alias    = "litmus-user-pass"
}
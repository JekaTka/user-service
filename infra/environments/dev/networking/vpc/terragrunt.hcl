terraform {
  source = "../../../../modules//networking/vpc"
}

include "root" {
  path   = find_in_parent_folders()
  expose = true
}

inputs = {
  cidr                   = "10.15.0.0"
  single_nat_gateway     = true
  one_nat_gateway_per_az = false
}
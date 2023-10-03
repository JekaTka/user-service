terraform {
  source = "../../../modules//kms"
}

include "root" {
  path   = find_in_parent_folders()
  expose = true
}

inputs = {
  aws_iam_service_linked_role_autoscaling = false
}
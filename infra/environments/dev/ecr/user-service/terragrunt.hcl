terraform {
  source = "../../../../modules//ecr"
}

include "root" {
  path   = find_in_parent_folders()
  expose = true
}

dependency "kms" {
  config_path                             = "../../kms/"
  mock_outputs_allowed_terraform_commands = ["init", "validate", "plan", "providers", "terragrunt-info", "show"]
  mock_outputs = {
    kms_deployment_key_arn = "arn:aws:kms:us-east-1:123456789000:key/fake-kms-key"
  }
}

inputs = {
  ecr_name             = "dev/user-service"
  image_tag_mutability = "IMMUTABLE"
  kms_key_id           = dependency.kms.outputs.kms_deployment_key_arn

  create_ecr_cache_repository         = false
  number_of_days_to_keep_cache_images = 14

  create_ecr_repository_policy = false
}
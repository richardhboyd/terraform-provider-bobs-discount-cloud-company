terraform {
  required_providers {
    bobsdiscountcloudco = {
      source = "hashicorp.com/edu/bobsdiscountcloudco"
    }
  }
}

provider "bobsdiscountcloudco" {
  host     = "https://api.us-east-1.whybobs.com"
  api_key = "SOME_API_KEY"
}

data "bobsdiscountcloudco_databases" "example" {}

resource "bobsdiscountcloudco_database" "example" {
  name = "steven"
  lifecycle {
    action_trigger {
      events    = [after_create]
      actions   = [action.bobsdiscountcloudco_population_action.message]
    }
  }
}

action "bobsdiscountcloudco_population_action" "message" {
  config {
    id = resource.bobsdiscountcloudco_database.example.id
    items = [{
      key = "richard"
      value = "Boyd"
    }]
  }
}
output "edu_databases" {
  value = resource.bobsdiscountcloudco_database.example.id
}

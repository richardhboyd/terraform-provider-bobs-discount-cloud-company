terraform {
  required_providers {
    hashicups = {
      source = "hashicorp.com/edu/hashicups"
    }
  }
}

provider "hashicups" {
  host     = "https://api.us-east-1.whybobs.com"
  api_key = "SOME_API_KEY"
}

data "hashicups_databases" "example" {}

resource "hashicups_database" "example" {
  name = "steven"
  lifecycle {
    action_trigger {
      events    = [after_create]
      actions   = [action.hashicups_population_action.message]
    }
  }
}

action "hashicups_population_action" "message" {
  config {
    id = resource.hashicups_database.example.id
    items = [{
      key = "richard"
      value = "Boyd"
    }]
  }
}
output "edu_databases" {
  value = resource.hashicups_database.example.id
}

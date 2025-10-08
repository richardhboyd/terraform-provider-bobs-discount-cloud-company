#!/bin/bash

version="0.2.2"
provider="bobs-discount-cloud-company"
provider_binary="terraform-provider-$provider_v$version"

# Function to create terraform.rc file
create_terraform_rc() {
    cat <<EOF > ~/.terraformrc
provider_installation {
  filesystem_mirror {
    path = "$PWD/.terraform/plugins"
  }
  direct {
    exclude = ["registry.terraform.io/*/*"]
  }
}
EOF
    echo "Created terraform.rc at ~/.terraformrc"
}

# Function to create provider.tf file
create_provider_config() {
    cat <<EOF > provider.tf
terraform {
  required_providers {
    $provider = {
      source = "hashicorp/$provider"
      version = "$version"
    }
  }
}
EOF
    echo "Created provider.tf configuration"
}

# Main setup
if [ ! -f "$provider_binary" ]; then
    echo "Error: Provider binary not found in current directory: $provider_binary" >&2
    exit 1
fi

# Create directory structure
provider_dir=".terraform/plugins/registry.terraform.io/hashicorp/$provider/$version/linux_amd64"
mkdir -p "$provider_dir"
echo "Created provider directory structure"

# Copy provider binary
cp "$provider_binary" "$provider_dir/"
echo "Copied provider binary to: $provider_dir/$provider_binary"

# Create terraform.rc file
create_terraform_rc

# Create provider configuration
create_provider_config

echo -e "\nSetup completed successfully!\n"
echo "Next steps:"
echo "1. Run twtrito initialize your working directory"
echo "2. Create your terraform configuration files (.tf)"
echo "3. Run 'terraform plan' to verify the setup"
#!/bin/bash

# Function to initialize a Go module in a directory
initialize_module() {
    local dir=$1
    local module_name=$2

    echo "Initializing module in $dir"
    cd "$dir" || exit 1
    rm -f go.mod go.sum
    go mod init "$module_name"
    go mod tidy
    cd - >/dev/null
}

# Remove all existing go.mod and go.sum files
find . -name "go.mod" -o -name "go.sum" -delete

# Initialize the root module
initialize_module "." "http_server"

# Initialize modules for each service
services=("auth-service" "api-gateway" "media-service" "post-service" "user-service" "shared")

for service in "${services[@]}"; do
    if [ -d "$service" ]; then
        initialize_module "$service" "http_server/$service"
    fi
done

echo "Module initialization complete!"

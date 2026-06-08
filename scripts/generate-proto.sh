#!/bin/bash

# Script to generate protobuf files

echo "Generating protobuf files..."

for dir in proto/*; do
    if [ -d "$dir" ]; then
        echo "Generating $dir..."
        protoc --go_out=. --go_opt=paths=source_relative \
            --go-grpc_out=. --go-grpc_opt=paths=source_relative \
            $dir/*.proto
    fi
done

echo "✓ Proto generation complete"

#!/bin/bash
# Helper script to run Claude in Docker sandbox for Ralph

set -e

# Check if ANTHROPIC_API_KEY is set
if [ -z "$ANTHROPIC_API_KEY" ]; then
    echo "Error: ANTHROPIC_API_KEY environment variable not set"
    echo "Export your key first: export ANTHROPIC_API_KEY='your-key'"
    exit 1
fi

# Build the image if needed
echo "Building Ralph sandbox image..."
docker-compose -f docker-compose.ralph.yml build

# Start the container
echo "Starting Ralph sandbox..."
docker-compose -f docker-compose.ralph.yml run --rm ralph bash

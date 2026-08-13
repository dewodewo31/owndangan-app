#!/bin/bash
set -euo pipefail

# Rollback script for Owndangan platform

ENVIRONMENT=${1:-staging}
PREVIOUS_VERSION=${2:-}

echo "========================================"
echo "Rolling back: ${ENVIRONMENT}"
echo "========================================"

# Validate environment
if [[ "$ENVIRONMENT" != "staging" && "$ENVIRONMENT" != "production" ]]; then
    echo "Error: Environment must be 'staging' or 'production'"
    exit 1
fi

# Load environment variables
if [[ -f ".env.${ENVIRONMENT}" ]]; then
    export $(cat .env.${ENVIRONMENT} | xargs)
fi

cd docker

if [[ -n "$PREVIOUS_VERSION" ]]; then
    echo "Rolling back to version: ${PREVIOUS_VERSION}"
    VERSION=${PREVIOUS_VERSION} docker-compose -f docker-compose.prod.yml up -d
else
    echo "Rolling back to previous version..."
    docker-compose -f docker-compose.prod.yml down
    docker-compose -f docker-compose.prod.yml up -d
fi

echo "Rollback complete!"
echo "Running health checks..."

# Health check
sleep 10
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ Backend health check passed"
else
    echo "❌ Backend health check failed"
    exit 1
fi

if curl -f http://localhost:3000 > /dev/null 2>&1; then
    echo "✅ Frontend health check passed"
else
    echo "❌ Frontend health check failed"
    exit 1
fi

echo "========================================"
echo "Rollback successful!"
echo "========================================"

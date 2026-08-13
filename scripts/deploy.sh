#!/bin/bash
set -euo pipefail

# Deployment script for Owndangan platform

ENVIRONMENT=${1:-staging}
VERSION=${2:-latest}

echo "========================================"
echo "Deploying to: ${ENVIRONMENT}"
echo "Version: ${VERSION}"
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

# Deploy
cd docker

case $ENVIRONMENT in
    staging)
        echo "Deploying to staging..."
        VERSION=${VERSION} docker-compose -f docker-compose.prod.yml up -d
        ;;
    production)
        echo "Deploying to production..."
        VERSION=${VERSION} docker-compose -f docker-compose.prod.yml up -d
        ;;
esac

echo "Deployment complete!"
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
echo "Deployment successful!"
echo "========================================"

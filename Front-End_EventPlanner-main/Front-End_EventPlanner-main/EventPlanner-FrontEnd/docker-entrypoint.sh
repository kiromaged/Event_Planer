#!/bin/ash

# Event Planner Frontend Docker Entrypoint Script
# Purpose: Initialize frontend environment at runtime

set -e

# Default backend URL configuration
BACKEND_URL="${BACKEND_URL:-http://backend:8080}"

echo "========================================="
echo "Event Planner Frontend - Initialization"
echo "========================================="
echo "Backend URL: $BACKEND_URL"
echo "Node Environment: ${NODE_ENV:-production}"
echo "========================================="

# Optional: If you need to inject environment variables into Angular's config at runtime,
# you could update config files here. For now, this is mainly a pass-through.

# Log startup information
echo "[$(date +'%Y-%m-%d %H:%M:%S')] Frontend container starting..."

# Execute the main command (nginx)
exec "$@"

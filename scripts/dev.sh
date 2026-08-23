#!/usr/bin/env bash
set -eo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

COMMAND="${1:-status}"
SERVICE="${2:-}"

case "$COMMAND" in
  start|up)
    echo "🚀 Starting containers..."
    if [ -n "$SERVICE" ]; then
      docker compose up -d "$SERVICE"
    else
      docker compose up -d
    fi
    docker compose ps --format "table {{.Name}}\t{{.Status}}\t{{.Ports}}"
    ;;
  stop|down)
    echo "🛑 Stopping containers..."
    docker compose down
    ;;
  restart)
    if [ -n "$SERVICE" ]; then
      echo "🔄 Restarting $SERVICE..."
      docker compose restart "$SERVICE"
    else
      echo "🔄 Restarting all containers..."
      docker compose restart
    fi
    docker compose ps --format "table {{.Name}}\t{{.Status}}\t{{.Ports}}"
    ;;
  status|ps)
    docker compose ps --format "table {{.Name}}\t{{.Status}}\t{{.Ports}}"
    ;;
  logs)
    if [ -n "$SERVICE" ]; then
      docker compose logs --tail=50 "$SERVICE"
    else
      docker compose logs --tail=30
    fi
    ;;
  build)
    if [ -n "$SERVICE" ]; then
      echo "🔨 Building $SERVICE..."
      docker compose build "$SERVICE"
    else
      echo "🔨 Building all containers..."
      docker compose build
    fi
    ;;
  *)
    echo "Usage: $0 {start|stop|restart|status|logs|build} [service_name]"
    exit 1
    ;;
esac

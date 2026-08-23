#!/usr/bin/env bash
set -eo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"

echo "📖 Gerando documentação Swagger..."
(cd "$BACKEND_DIR" && go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go)
echo "✓ [Swagger] Documentação gerada com sucesso em backend/docs/"

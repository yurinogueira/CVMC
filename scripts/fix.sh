#!/usr/bin/env bash
set -eo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "🔧 Auto-formatting backend..."
(cd "$PROJECT_ROOT/backend" && go fmt ./...)

echo "🔧 Auto-formatting & lint-fixing frontend..."
if [ ! -d "$PROJECT_ROOT/frontend/node_modules" ]; then
  (cd "$PROJECT_ROOT/frontend" && npm ci)
fi
(cd "$PROJECT_ROOT/frontend" && ./node_modules/.bin/prettier --write . --log-level warn && ./node_modules/.bin/eslint . --fix)

if command -v terraform &> /dev/null; then
  echo "🔧 Auto-formatting terraform..."
  (cd "$PROJECT_ROOT" && terraform fmt -recursive terraform)
fi

echo "✓ [CVMC] Formatting and fixes applied successfully!"

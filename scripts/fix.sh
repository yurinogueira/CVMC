#!/usr/bin/env bash
set -eo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "🔧 Auto-formatting backend..."
(cd "$PROJECT_ROOT/backend" && go fmt ./...)

echo "🔧 Auto-formatting & lint-fixing frontend..."
(cd "$PROJECT_ROOT/frontend" && npx prettier --write . --log-level warn && npm run lint -- --fix)

echo "✓ [CVMC] Formatting and fixes applied successfully!"

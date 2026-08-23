#!/usr/bin/env bash
set -eo pipefail

# Root directory of CVMC project
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

check_backend() {
  local backend_dir="$PROJECT_ROOT/backend"
  local output=""
  local status=0

  # 1. Vet
  if ! output=$(cd "$backend_dir" && go vet ./... 2>&1); then
    echo "❌ [Backend] go vet failed:"
    echo "$output"
    return 1
  fi

  # 2. Test (suppress '? cvmc/... [no test files]' noise)
  output=$(cd "$backend_dir" && go test ./... 2>&1) || status=$?
  if [ $status -ne 0 ]; then
    echo "❌ [Backend] go test failed:"
    echo "$output" | grep -v "^? "
    return 1
  fi

  # Count tested packages
  local passed_count
  passed_count=$(echo "$output" | grep -c "^ok " || true)
  echo "✓ [Backend] OK (vet passed, ${passed_count} test packages passed)"
  return 0
}

check_frontend() {
  local frontend_dir="$PROJECT_ROOT/frontend"
  local output=""

  # 1. TypeCheck
  if ! output=$(cd "$frontend_dir" && npx tsc -b 2>&1); then
    echo "❌ [Frontend] Typecheck (tsc) failed:"
    echo "$output"
    return 1
  fi

  # 2. Lint
  if ! output=$(cd "$frontend_dir" && npm run lint --silent 2>&1); then
    echo "❌ [Frontend] ESLint failed:"
    echo "$output"
    return 1
  fi

  # 3. Format
  if ! output=$(cd "$frontend_dir" && npm run format --silent 2>&1); then
    echo "❌ [Frontend] Prettier format check failed (run scripts/fix.sh to auto-fix):"
    echo "$output"
    return 1
  fi

  # 4. Vitest (concise reporter)
  if ! output=$(cd "$frontend_dir" && npx vitest run --reporter=dot 2>&1); then
    echo "❌ [Frontend] Vitest failed:"
    echo "$output"
    return 1
  fi

  echo "✓ [Frontend] OK (typecheck, lint, format & tests passed)"
  return 0
}

TARGET="${1:-all}"
FAILED=0

case "$TARGET" in
  backend)
    check_backend || FAILED=1
    ;;
  frontend)
    check_frontend || FAILED=1
    ;;
  all)
    check_backend || FAILED=1
    check_frontend || FAILED=1
    if [ $FAILED -eq 0 ]; then
      echo "🚀 [CVMC] All checks passed successfully!"
    fi
    ;;
  *)
    echo "Usage: $0 [backend|frontend|all]"
    exit 1
    ;;
esac

exit $FAILED

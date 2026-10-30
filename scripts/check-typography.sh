#!/usr/bin/env bash
set -eo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FRONTEND_DIR="$PROJECT_ROOT/frontend"

echo "🔍 [Typography] Verifying visual hierarchy, PageHeader contract & Topbar route synchronization..."

PAGES=(
  "src/features/dashboard/pages/DashboardPage.tsx"
  "src/features/cars/pages/VehiclesPage.tsx"
  "src/features/cars/pages/VehicleDetailsPage.tsx"
  "src/features/maintenance/pages/MaintenancePage.tsx"
  "src/features/maintenance/pages/RegisterMaintenancePage.tsx"
  "src/features/reports/pages/ReportsPage.tsx"
  "src/features/profile/pages/ProfilePage.tsx"
)

ERRORS=0

for page in "${PAGES[@]}"; do
  PAGE_PATH="$FRONTEND_DIR/$page"
  if [ ! -f "$PAGE_PATH" ]; then
    echo "❌ [Typography] Expected protected page not found: $page"
    ERRORS=$((ERRORS + 1))
    continue
  fi

  # 1. Check PageHeader import (multiline-safe)
  if ! grep -zoP "import\s*\{[^}]*PageHeader[^}]*\}\s*from" "$PAGE_PATH" > /dev/null 2>&1; then
    echo "❌ [Typography] Page $page does not import PageHeader from shared"
    ERRORS=$((ERRORS + 1))
  fi

  # 2. Check <PageHeader rendering
  if ! grep -q "<PageHeader" "$PAGE_PATH"; then
    echo "❌ [Typography] Page $page does not render <PageHeader"
    ERRORS=$((ERRORS + 1))
  fi

  # 3. Check for forbidden ad-hoc component="h1" outside PageHeader
  if grep -q -E "<Typography[^>]*component=\"h1\"" "$PAGE_PATH"; then
    echo "❌ [Typography] Page $page contains forbidden ad-hoc <Typography component=\"h1\">. Use PageHeader instead."
    ERRORS=$((ERRORS + 1))
  fi

  # 4. Check for forbidden ad-hoc variant="h4" page titles
  if grep -q -E "<Typography[^>]*variant=\"h4\"" "$PAGE_PATH"; then
    echo "❌ [Typography] Page $page contains forbidden ad-hoc <Typography variant=\"h4\">. Use PageHeader instead."
    ERRORS=$((ERRORS + 1))
  fi
done

# 5. Verify Topbar route mappings
TOPBAR_PATH="$FRONTEND_DIR/src/layouts/Topbar.tsx"
MANDATORY_ROUTES=(
  "/dashboard"
  "/vehicles"
  "/maintenance"
  "/reports"
  "/profile"
)

for route in "${MANDATORY_ROUTES[@]}"; do
  if ! grep -q "\"$route\"" "$TOPBAR_PATH"; then
    echo "❌ [Typography] Route $route is missing from staticRouteTitles in Topbar.tsx"
    ERRORS=$((ERRORS + 1))
  fi
done

if [ $ERRORS -ne 0 ]; then
  echo "❌ [Typography] Failed with $ERRORS visual consistency violations."
  echo "📖 Refer to docs/frontend/ui-components.md for the canonical PageHeader contract."
  exit 1
fi

echo "✓ [Typography] OK (PageHeader contract and Topbar routes validated)"
exit 0

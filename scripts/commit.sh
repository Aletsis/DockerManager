#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MSG="${1:-}"

if [ -z "$MSG" ]; then
  echo "Error: Debes proporcionar un mensaje de commit."
  echo "Ejemplo: ./scripts/commit.sh \"feat: agregar soporte de volumenes\""
  exit 1
fi

# Bump version based on commit message
"${ROOT_DIR}/scripts/bump-version.sh" auto "$MSG"

# Stage version files
git add \
  "${ROOT_DIR}/version.json" \
  "${ROOT_DIR}/internal/version/version.go" \
  "${ROOT_DIR}/frontend/src/shared/version.ts" \
  "${ROOT_DIR}/frontend/package.json" \
  "${ROOT_DIR}/wails.json"

# Execute commit
GIT_AMENDING=1 git commit -m "$MSG"

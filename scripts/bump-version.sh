#!/usr/bin/env bash
set -euo pipefail

# Root directory of DockerManager
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION_JSON="${ROOT_DIR}/version.json"
VERSION_GO="${ROOT_DIR}/internal/version/version.go"
VERSION_TS="${ROOT_DIR}/frontend/src/shared/version.ts"
PACKAGE_JSON="${ROOT_DIR}/frontend/package.json"
WAILS_JSON="${ROOT_DIR}/wails.json"

# Read current version from version.json
get_current_version() {
  if [ -f "$VERSION_JSON" ]; then
    grep -o '"version"[[:space:]]*:[[:space:]]*"[^"]*"' "$VERSION_JSON" | sed -E 's/.*"([^"]+)".*/\1/'
  else
    echo "0.1.0"
  fi
}

CURRENT_VERSION="$(get_current_version)"
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"

# Default if any part missing
MAJOR=${MAJOR:-0}
MINOR=${MINOR:-1}
PATCH=${PATCH:-0}

update_files() {
  local NEW_VERSION="$1"
  echo "Actualizando versión a v${NEW_VERSION}..."

  # 1. version.json
  cat << JSON_EOF > "$VERSION_JSON"
{
  "version": "${NEW_VERSION}"
}
JSON_EOF

  # 2. internal/version/version.go
  mkdir -p "$(dirname "$VERSION_GO")"
  cat << GO_EOF > "$VERSION_GO"
package version

// Current holds the current application version of DockerManager (SemVer format X.Y.Z)
const Current = "${NEW_VERSION}"
GO_EOF

  # 3. frontend/src/shared/version.ts
  mkdir -p "$(dirname "$VERSION_TS")"
  cat << TS_EOF > "$VERSION_TS"
// Single source of truth mirrored for frontend
export const APP_VERSION = '${NEW_VERSION}';
TS_EOF

  # 4. frontend/package.json
  if [ -f "$PACKAGE_JSON" ]; then
    sed -i -E "s/(\"version\"[[:space:]]*:[[:space:]]*)\"[^\"]*\"/\1\"${NEW_VERSION}\"/" "$PACKAGE_JSON"
  fi

  # 5. wails.json
  if [ -f "$WAILS_JSON" ]; then
    sed -i -E "s/(\"productVersion\"[[:space:]]*:[[:space:]]*)\"[^\"]*\"/\1\"${NEW_VERSION}\"/" "$WAILS_JSON"
  fi

  echo "✓ Archivos sincronizados exitosamente con la versión ${NEW_VERSION}."
}

# Determine bump type from commit message
# Rules:
# - BREAKING CHANGE or <type>!: -> major
# - feat: / feat(...): -> minor
# - fix:, perf:, refactor:, chore: -> patch
# - docs:, test:, etc. -> none (mantiene versión)
detect_bump_type() {
  local MSG="$1"

  if echo "$MSG" | grep -Eq 'BREAKING CHANGE|^[a-zA-Z0-9_-]+(\([^\)]+\))?!:'; then
    echo "major"
  elif echo "$MSG" | grep -Eq '^feat(\([^\)]+\))?:'; then
    echo "minor"
  elif echo "$MSG" | grep -Eq '^(fix|perf|refactor|chore)(\([^\)]+\))?:'; then
    echo "patch"
  elif echo "$MSG" | grep -Eq '^(docs|test|ci|style|build)(\([^\)]+\))?:'; then
    echo "none"
  else
    # Si no tiene prefijo convencional, se mantiene la versión
    echo "none"
  fi
}

ACTION="${1:-help}"

case "$ACTION" in
  get)
    echo "$CURRENT_VERSION"
    ;;
  major)
    NEW_MAJOR=$((MAJOR + 1))
    NEW_VERSION="${NEW_MAJOR}.0.0"
    update_files "$NEW_VERSION"
    ;;
  minor)
    NEW_MINOR=$((MINOR + 1))
    NEW_VERSION="${MAJOR}.${NEW_MINOR}.0"
    update_files "$NEW_VERSION"
    ;;
  patch)
    NEW_PATCH=$((PATCH + 1))
    NEW_VERSION="${MAJOR}.${MINOR}.${NEW_PATCH}"
    update_files "$NEW_VERSION"
    ;;
  set)
    if [ -z "${2:-}" ]; then
      echo "Error: Debe especificar la versión X.Y.Z (ej. ./scripts/bump-version.sh set 0.2.0)" >&2
      exit 1
    fi
    NEW_VERSION="${2#v}"
    update_files "$NEW_VERSION"
    ;;
  auto)
    MSG="${2:-}"
    if [ -z "$MSG" ]; then
      echo "Error: Debe proporcionar el mensaje del commit para el modo auto." >&2
      exit 1
    fi
    BUMP_TYPE="$(detect_bump_type "$MSG")"
    case "$BUMP_TYPE" in
      major)
        echo "Detectado cambio mayor / breaking change. Incrementando Mayor (X)..."
        NEW_MAJOR=$((MAJOR + 1))
        update_files "${NEW_MAJOR}.0.0"
        ;;
      minor)
        echo "Detectada nueva funcionalidad (feat). Incrementando Menor (Y)..."
        NEW_MINOR=$((MINOR + 1))
        update_files "${MAJOR}.${NEW_MINOR}.0"
        ;;
      patch)
        echo "Detectado parche/mantenimiento (fix/perf/refactor/chore). Incrementando Parche (Z)..."
        NEW_PATCH=$((PATCH + 1))
        update_files "${MAJOR}.${MINOR}.${NEW_PATCH}"
        ;;
      none)
        echo "Detectado commit sin incremento (docs/test/otro). Se mantiene la versión ${CURRENT_VERSION}."
        ;;
    esac
    ;;
  *)
    echo "Uso: $0 {get|patch|minor|major|set <X.Y.Z>|auto <mensaje-commit>}"
    exit 1
    ;;
esac

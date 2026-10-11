#!/usr/bin/env bash
set -eo pipefail

# Root directory of CVMC project
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DOCS_DIR="$PROJECT_ROOT/docs"

if [ ! -d "$DOCS_DIR" ]; then
  echo "❌ [Docs] Diretório 'docs/' não encontrado em $DOCS_DIR"
  exit 1
fi

ALLOWED_TYPES="index|log|architecture|domain|frontend|operations|adr|runbook|concept"

VERIFY_SYNC=false
for arg in "$@"; do
  if [ "$arg" = "--verify-sync" ]; then
    VERIFY_SYNC=true
  fi
done

ERRORS=0
FILE_COUNT=0
LINK_COUNT=0

log_error() {
  local file="$1"
  local msg="$2"
  echo "❌ [Docs] $(realpath --relative-to="$PROJECT_ROOT" "$file"): $msg"
  ERRORS=$((ERRORS + 1))
}

# 1. Validar Frontmatter e Metadados de cada arquivo .md
while IFS= read -r file; do
  FILE_COUNT=$((FILE_COUNT + 1))

  # Linha 1 deve ser ---
  first_line=$(head -n 1 "$file" || true)
  if [ "$first_line" != "---" ]; then
    log_error "$file" "Frontmatter ausente ou não inicia com '---' na linha 1"
    continue
  fi

  # Deve existir um segundo --- fechando o frontmatter
  closing_line=$(awk 'NR>1 && /^---$/ { print NR; exit }' "$file")
  if [ -z "$closing_line" ]; then
    log_error "$file" "Frontmatter não possui delimitador de fechamento '---'"
    continue
  fi

  # Extrair bloco frontmatter
  frontmatter=$(sed -n "2,$((closing_line - 1))p" "$file")

  # Extrair type
  doc_type=$(echo "$frontmatter" | grep -E '^type:' | head -n 1 | awk '{print $2}' | tr -d '"' | tr -d "'" || true)
  if [ -z "$doc_type" ]; then
    log_error "$file" "Campo obrigatório 'type:' ausente no frontmatter"
  elif ! echo "$doc_type" | grep -Eq "^($ALLOWED_TYPES)$"; then
    log_error "$file" "Tipo '$doc_type' inválido. Tipos permitidos: $ALLOWED_TYPES"
  fi

  # Extrair title
  if ! echo "$frontmatter" | grep -Eq '^title:'; then
    log_error "$file" "Campo obrigatório 'title:' ausente no frontmatter"
  fi

  # Extrair description
  if ! echo "$frontmatter" | grep -Eq '^description:'; then
    log_error "$file" "Campo obrigatório 'description:' ausente no frontmatter"
  fi

  # 2. Validar Links Markdown locais no arquivo
  dir_name=$(dirname "$file")
  
  # Extrair alvos dos links no formato [label](target)
  # Ignora links http://, https://, mailto:, e âncoras puras #
  while IFS= read -r link_target; do
    [ -z "$link_target" ] && continue

    # Remover eventual âncora (#...)
    clean_target="${link_target%%#*}"
    [ -z "$clean_target" ] && continue

    LINK_COUNT=$((LINK_COUNT + 1))

    # Resolver caminho absoluto
    if [[ "$clean_target" = /* ]]; then
      resolved="$PROJECT_ROOT$clean_target"
    else
      resolved="$dir_name/$clean_target"
    fi

    # Normalizar caminho
    resolved_normalized=$(cd "$dir_name" && realpath -m "$clean_target" 2>/dev/null || true)

    if [ ! -f "$resolved_normalized" ]; then
      log_error "$file" "Link quebrado para '$link_target' (resolvido: '$resolved_normalized')"
    fi
  done < <(grep -o -E '\[[^]]+\]\(([^)]+)\)' "$file" | sed -E 's/.*\]\(([^)]+)\)/\1/' | grep -v -E '^(https?://|mailto:|#)' || true)

done < <(find "$DOCS_DIR" -type f -name "*.md" | sort)

# 3. Validar se todos os documentos de docs/ estão indexados em docs/index.md ou docs/log.md
INDEX_FILE="$DOCS_DIR/index.md"
LOG_FILE="$DOCS_DIR/log.md"
if [ ! -f "$INDEX_FILE" ]; then
  log_error "$DOCS_DIR" "Catálogo central 'docs/index.md' não encontrado"
else
  while IFS= read -r file; do
    # Ignorar index.md e log.md da exigência de auto-indexação
    rel_path=$(realpath --relative-to="$DOCS_DIR" "$file")
    if [ "$rel_path" = "index.md" ] || [ "$rel_path" = "log.md" ]; then
      continue
    fi

    if echo "$rel_path" | grep -q "^logs/"; then
      if ! grep -Fq "$rel_path" "$LOG_FILE"; then
        log_error "$file" "Arquivo diário de log não está referenciado no catálogo 'docs/log.md'"
      fi
    else
      if ! grep -Fq "$rel_path" "$INDEX_FILE"; then
        log_error "$file" "Documento não está referenciado no catálogo central 'docs/index.md'"
      fi
    fi
  done < <(find "$DOCS_DIR" -type f -name "*.md" | sort)
fi

# 4. Validar integridade e estrutura cronológica de docs/log.md e docs/logs/
LOG_FILE="$DOCS_DIR/log.md"
if [ ! -f "$LOG_FILE" ]; then
  log_error "$DOCS_DIR" "Arquivo de histórico 'docs/log.md' não encontrado"
else
  # Validar formato ISO YYYY-MM-DD no frontmatter timestamp
  log_fm_ts=$(grep -E '^timestamp:' "$LOG_FILE" | head -n 1 | awk '{print $2}' | tr -d '"' | tr -d "'" || true)
  if [ -z "$log_fm_ts" ] || ! echo "$log_fm_ts" | grep -Eq '^[0-9]{4}-[0-9]{2}-[0-9]{2}$'; then
    log_error "$LOG_FILE" "Campo 'timestamp:' em docs/log.md deve estar no formato ISO AAAA-MM-DD (ex: 2026-10-10)"
  fi

  # Validar presença de seções datadas
  first_entry_date=$(grep -E '^## 📅 \[?[0-9]{4}-[0-9]{2}-[0-9]{2}' "$LOG_FILE" | head -n 1 | sed -E 's/^## 📅 \[?([0-9]{4}-[0-9]{2}-[0-9]{2}).*/\1/' || true)
  if [ -z "$first_entry_date" ]; then
    log_error "$LOG_FILE" "docs/log.md não contém nenhuma seção no formato obrigatório '## 📅 [AAAA-MM-DD](logs/AAAA-MM-DD.md) — <Título>' ou '## 📅 AAAA-MM-DD — <Título>'"
  elif [ -n "$log_fm_ts" ] && [ "$log_fm_ts" != "$first_entry_date" ]; then
    log_error "$LOG_FILE" "Divergência em docs/log.md: o frontmatter timestamp ('$log_fm_ts') não coincide com a data da entrada mais recente ('$first_entry_date')"
  fi
fi

# Validar arquivos diários sob docs/logs/
if [ -d "$DOCS_DIR/logs" ]; then
  for daily_log in "$DOCS_DIR/logs"/*.md; do
    [ -e "$daily_log" ] || continue
    daily_base=$(basename "$daily_log")
    if ! echo "$daily_base" | grep -Eq '^[0-9]{4}-[0-9]{2}-[0-9]{2}\.md$'; then
      log_error "$daily_log" "Nome de arquivo diário de log inválido. Padrão obrigatório: AAAA-MM-DD.md (ex: 2026-10-10.md)"
    fi
  done
fi

# 5. Se --verify-sync ativo, validar se mudanças em código/infra possuem correspondente em docs/
if [ "$VERIFY_SYNC" = true ]; then
  if command -v git &>/dev/null && git -C "$PROJECT_ROOT" rev-parse --is-inside-work-tree &>/dev/null; then
    diff_target=""
    if git -C "$PROJECT_ROOT" rev-parse --verify origin/main &>/dev/null; then
      diff_target="origin/main"
    elif git -C "$PROJECT_ROOT" rev-parse --verify main &>/dev/null; then
      diff_target="main"
    elif git -C "$PROJECT_ROOT" rev-parse --verify HEAD~1 &>/dev/null; then
      diff_target="HEAD~1"
    fi

    changed_files=$( { [ -n "$diff_target" ] && git -C "$PROJECT_ROOT" diff --name-only "$diff_target" 2>/dev/null || true; git -C "$PROJECT_ROOT" status --porcelain -uall 2>/dev/null | sed -E 's/^...//; s/.* -> //'; } | sort -u )

    if [ -n "$changed_files" ]; then
      code_patterns="^(backend/|frontend/|terraform/|deploy/|scripts/|\.github/|\.agents/)"
      code_changed=$(echo "$changed_files" | grep -E "$code_patterns" || true)

      if [ -n "$code_changed" ]; then
        docs_changed=$(echo "$changed_files" | grep -E "^docs/" || true)
        log_changed=$(echo "$changed_files" | grep -E "^docs/(log\.md|logs/)" || true)

        if [ -z "$docs_changed" ]; then
          log_error "$DOCS_DIR" "Zero-Divergence Violation: Arquivos de código/infraestrutura foram alterados, mas nenhum documento sob 'docs/' foi atualizado."
        elif [ -z "$log_changed" ]; then
          log_error "$DOCS_DIR/log.md" "Zero-Divergence Violation: Modificações de código e docs/ detectadas, mas nem 'docs/log.md' nem arquivos em 'docs/logs/' foram atualizados com o registro da intervenção."
        fi
      fi
    fi
  else
    echo "ℹ️ [Docs] Ambiente Git não detectado; verificação de sincronização ignorada."
  fi
fi

if [ $ERRORS -ne 0 ]; then
  echo "❌ [Docs] Falha na validação da documentação: $ERRORS erro(s) encontrado(s)."
  exit 1
fi

sync_msg=""
if [ "$VERIFY_SYNC" = true ]; then
  sync_msg=" [Zero-Divergence verificado]"
fi

echo "✓ [Docs] OK ($FILE_COUNT arquivos e $LINK_COUNT links locais validados com sucesso$sync_msg)"
exit 0

#!/usr/bin/env bash
set -euo pipefail

# Uso: ./scripts/sync-env.sh [IP_OU_HOST_OCI] [USUARIO_SSH (default: ubuntu)] [CAMINHO_SSH_KEY]
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
TF_DIR="$ROOT_DIR/terraform"

echo "==> 1. Extraindo outputs do Terraform..."
if [ ! -d "$TF_DIR" ]; then
    echo "Erro: Diretório terraform não encontrado em $TF_DIR"
    exit 1
fi

TF_OUTPUTS=$(cd "$TF_DIR" && terraform output -json 2>/dev/null || true)

if [ -z "$TF_OUTPUTS" ]; then
    echo "Erro: Não foi possível obter os outputs do Terraform."
    exit 1
fi

PARSED=$(python3 -c '
import json, sys
data = json.loads(sys.stdin.read())
def get_val(key, default=""):
    return data.get(key, {}).get("value", default)

ip = get_val("instance_public_ip")
mongo_uri = get_val("mongodb_uri")
bucket = get_val("object_storage_bucket", "project-files")
namespace = get_val("object_storage_namespace", "")
print(f"{ip}\t{mongo_uri}\t{bucket}\t{namespace}")
' <<< "$TF_OUTPUTS")

TF_HOST=$(echo "$PARSED" | cut -f1)
MONGO_URI=$(echo "$PARSED" | cut -f2)
BUCKET_NAME=$(echo "$PARSED" | cut -f3)
BUCKET_NAMESPACE=$(echo "$PARSED" | cut -f4)

HOST="${1:-$TF_HOST}"
USER="${2:-ubuntu}"
KEY="${3:-}"

if [ -z "$HOST" ]; then
    echo "Erro: Host/IP da instância não encontrado."
    exit 1
fi

if [ -z "$MONGO_URI" ]; then
    echo "Erro: mongodb_uri não encontrado no output do Terraform."
    exit 1
fi

SSH_OPTS=("-o" "StrictHostKeyChecking=accept-new")
if [ -n "$KEY" ]; then
    SSH_OPTS+=("-i" "$KEY")
fi

echo "==> 2. Gerando segredos JWT e configurações para $USER@$HOST..."
JWT_SECRET=$(openssl rand -base64 32)
JWT_REFRESH_SECRET=$(openssl rand -base64 32)

ENV_CONTENT="# Configurações de Produção CVMC (Gerado automaticamente)
PORT=8080
LOG_LEVEL=info
STORAGE_PROVIDER=local
UPLOAD_PATH=/opt/cvmc/data/uploads

# Domínios e Segurança
ALLOWED_ORIGINS=https://cvmc.yurinogueira.dev.br
COOKIE_DOMAIN=.yurinogueira.dev.br
COOKIE_SECURE=true

# MongoDB Atlas
MONGO_URI=${MONGO_URI}
MONGO_DATABASE=cvmc

# OCI Object Storage Always Free
OCI_STORAGE_BUCKET=${BUCKET_NAME}
OCI_STORAGE_NAMESPACE=${BUCKET_NAMESPACE}
OCI_STORAGE_REGION=sa-saopaulo-1

# JWT Secrets
JWT_SECRET=${JWT_SECRET}
JWT_REFRESH_SECRET=${JWT_REFRESH_SECRET}
"

echo "==> 3. Enviando /etc/cvmc/backend.env para o servidor OCI..."
ssh "${SSH_OPTS[@]}" "$USER@$HOST" "sudo mkdir -p /etc/cvmc && sudo tee /etc/cvmc/backend.env > /dev/null && sudo chmod 600 /etc/cvmc/backend.env && sudo chown root:cvmc /etc/cvmc/backend.env" <<< "$ENV_CONTENT"

echo "==> backend.env configurado com sucesso no servidor!"

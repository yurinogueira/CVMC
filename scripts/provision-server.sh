#!/usr/bin/env bash
set -euo pipefail

# Uso: ./scripts/provision-server.sh <IP_OU_HOST_OCI> [USUARIO_SSH (default: ubuntu)] [CAMINHO_SSH_KEY]
HOST="${1:-}"
USER="${2:-ubuntu}"
KEY="${3:-}"

if [ -z "$HOST" ]; then
    echo "Uso: $0 <IP_OU_HOST_OCI> [USUARIO_SSH] [CHAVE_SSH]"
    echo "Exemplo: $0 203.0.113.10 ubuntu"
    exit 1
fi

SSH_OPTS=("-o" "StrictHostKeyChecking=accept-new")
if [ -n "$KEY" ]; then
    SSH_OPTS+=("-i" "$KEY")
fi

echo "==> Conectando em $USER@$HOST para provisionar o servidor..."

ssh "${SSH_OPTS[@]}" "$USER@$HOST" << 'EOF'
set -euo pipefail

echo "==> 1. Verificando Swap..."
if [ ! -f /swapfile ]; then
    echo "Configurando Swap de 2GB..."
    sudo fallocate -l 2G /swapfile || sudo dd if=/dev/zero of=/swapfile bs=1M count=2048
    sudo chmod 600 /swapfile
    sudo mkswap /swapfile
    sudo swapon /swapfile
    echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
    sudo sysctl vm.swappiness=10
    echo 'vm.swappiness=10' | sudo tee /etc/sysctl.d/99-swap.conf
else
    echo "Swap já configurado."
fi

echo "==> 2. Instalando Caddy e utilitários..."
export DEBIAN_FRONTEND=noninteractive
sudo apt-get update -y
sudo apt-get install -y debian-keyring debian-archive-keyring apt-transport-https curl iptables-persistent

if ! command -v caddy &> /dev/null; then
    echo "Instalando pacote oficial do Caddy..."
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
    sudo apt-get update -y
    sudo apt-get install -y caddy
else
    echo "Caddy já instalado: $(caddy version)"
fi

echo "==> 3. Configurando usuário e diretórios CVMC..."
id -u cvmc &>/dev/null || sudo useradd -r -s /bin/false -d /opt/cvmc cvmc
sudo mkdir -p /opt/cvmc/data/uploads /etc/cvmc /var/log/caddy
sudo chown -R cvmc:cvmc /opt/cvmc
sudo chown -R root:cvmc /etc/cvmc
sudo chmod 750 /etc/cvmc
sudo chown -R caddy:caddy /var/log/caddy

echo "==> 4. Configurando Firewall (Portas 80 e 443)..."
sudo iptables -I INPUT 5 -m state --state NEW -p tcp --dport 80 -j ACCEPT 2>/dev/null || sudo iptables -A INPUT -p tcp --dport 80 -j ACCEPT
sudo iptables -I INPUT 5 -m state --state NEW -p tcp --dport 443 -j ACCEPT 2>/dev/null || sudo iptables -A INPUT -p tcp --dport 443 -j ACCEPT
sudo netfilter-persistent save || true

sudo systemctl enable caddy
sudo systemctl start caddy || true

echo "==> Servidor provisionado e pronto para deploy!"
EOF

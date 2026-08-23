#!/usr/bin/env bash
set -euo pipefail

# 1. Configurar Swap de 2GB (Essencial para VM de 1GB de RAM)
if [ ! -f /swapfile ]; then
    echo "=== Configurando Swap de 2GB ==="
    fallocate -l 2G /swapfile || dd if=/dev/zero of=/swapfile bs=1M count=2048
    chmod 600 /swapfile
    mkswap /swapfile
    swapon /swapfile
    echo '/swapfile none swap sw 0 0' >> /etc/fstab
    sysctl vm.swappiness=10
    echo 'vm.swappiness=10' >> /etc/sysctl.d/99-swap.conf
fi

# 2. Atualizar pacotes base e instalar utilitários
export DEBIAN_FRONTEND=noninteractive
apt-get update -y
apt-get install -y debian-keyring debian-archive-keyring apt-transport-https curl ufw iptables-persistent

# 3. Instalar Caddy (Repositório Oficial)
if ! command -v caddy &> /dev/null; then
    echo "=== Instalando Caddy ==="
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | tee /etc/apt/sources.list.d/caddy-stable.list
    apt-get update -y
    apt-get install -y caddy
fi

# 4. Criar usuário e diretórios da aplicação CVMC
echo "=== Configurando usuário e diretórios ==="
id -u cvmc &>/dev/null || useradd -r -s /bin/false -d /opt/cvmc cvmc
mkdir -p /opt/cvmc /etc/cvmc
chown -R cvmc:cvmc /opt/cvmc
chown -R root:cvmc /etc/cvmc
chmod 750 /etc/cvmc

# 5. Configurar regras de firewall na instância OCI (Ubuntu)
echo "=== Configurando regras de Firewall (Portas 80 e 443) ==="
iptables -I INPUT 5 -m state --state NEW -p tcp --dport 80 -j ACCEPT || true
iptables -I INPUT 5 -m state --state NEW -p tcp --dport 443 -j ACCEPT || true
netfilter-persistent save || true

systemctl enable caddy
systemctl restart caddy

echo "=== Provisionamento concluído com sucesso ==="

#!/bin/bash

set -e

# ===========================================
# piunter installer
# https://github.com/joaomjbraga/piunter
# ===========================================

REPO="joaomjbraga/piunter"
BINARY_NAME="piunter-linux-amd64"
INSTALL_DIR="/usr/local/bin"
USER_INSTALL_DIR="$HOME/.local/bin"

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info() {
    echo -e "${GREEN}[*]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[!]${NC} $1"
}

error() {
    echo -e "${RED}[x]${NC} $1"
}

# Verifica se é Linux
if [[ "$(uname)" != "Linux" ]]; then
    error "Este script funciona apenas no Linux."
    exit 1
fi

# Detecta arquitetura
ARCH=$(uname -m)
case "$ARCH" in
    x86_64)
        ARCH_NAME="amd64"
        ;;
    aarch64|arm64)
        ARCH_NAME="arm64"
        ;;
    *)
        error "Arquitetura não suportada: $ARCH"
        exit 1
        ;;
esac
BINARY_NAME="piunter-linux-${ARCH_NAME}"

info "Arquitetura detectada: $ARCH_NAME"

# Determina diretório de instalação
if [[ -d "$USER_INSTALL_DIR" ]] && [[ ! -w "/usr/local/bin" ]]; then
    INSTALL_DIR="$USER_INSTALL_DIR"
fi

info "Instalando em: $INSTALL_DIR"

# Verifica curl
if ! command -v curl &> /dev/null; then
    error "curl é necessário. Instale com: sudo apt install curl"
    exit 1
fi

# Pega versão mais recente
info "Buscando última versão..."
VERSION=$(curl -sL "https://api.github.com/repos/${REPO}/releases/latest" | grep -o '"tag_name":.*' | cut -d'"' -f4)

if [[ -z "$VERSION" ]]; then
    warn "Não foi possível buscar versão. Usando 'latest'."
    VERSION="latest"
fi

info "Versão: $VERSION"

# URL de download
if [[ "$VERSION" == "latest" ]]; then
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}"
else
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}"
fi

info "Baixando de: $DOWNLOAD_URL"

# Download
TEMP_FILE=$(mktemp)
trap "rm -f $TEMP_FILE" EXIT

curl -fsL "$DOWNLOAD_URL" -o "$TEMP_FILE" || {
    error "Falha ao baixar. Verifique se a release existe."
    exit 1
}

# Verifica se o arquivo é válido (não é HTML de erro)
if file "$TEMP_FILE" | grep -q "HTML"; then
    error "Download retornou HTML. A release pode não existir."
    exit 1
fi

# Instala
info "Instalando..."

# Verifica se precisa de sudo
if [[ ! -w "$INSTALL_DIR" ]]; then
    warn "Precisa de privilégios de administrador."
    sudo cp "$TEMP_FILE" "${INSTALL_DIR}/piunter"
    sudo chmod 755 "${INSTALL_DIR}/piunter"
else
    cp "$TEMP_FILE" "${INSTALL_DIR}/piunter"
    chmod 755 "${INSTALL_DIR}/piunter"
fi

info "Instalado com sucesso!"

# Confirma
if command -v piunter &> /dev/null; then
    info "Executando piunter --version:"
    piunter --version 2>/dev/null || piunter --help | head -1 || true
else
    warn "Adicione $INSTALL_DIR ao PATH se necessário."
fi

echo
info "Pronto! Use 'piunter --help' para começar."
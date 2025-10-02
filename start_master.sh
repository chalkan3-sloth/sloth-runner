#!/bin/bash
# Script para iniciar o master server do Sloth Runner

set -e

# Configuração
MASTER_IP="192.168.1.29"
MASTER_PORT="50053"
DAEMON=true

# Cores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Iniciando Sloth Runner Master Server${NC}"
echo ""

# Verificar se já está rodando
if pgrep -f "sloth-runner master" > /dev/null; then
    echo -e "${RED}⚠️  Master server já está rodando!${NC}"
    echo ""
    echo "Para parar o master existente, execute:"
    echo "  pkill -f 'sloth-runner master'"
    echo ""
    exit 1
fi

# Verificar se o executável existe
if [ ! -f "./sloth-runner" ]; then
    echo -e "${RED}❌ Executável sloth-runner não encontrado!${NC}"
    echo ""
    echo "Certifique-se de estar no diretório correto ou instale o sloth-runner."
    exit 1
fi

# Iniciar master
echo -e "${BLUE}📡 Configuração:${NC}"
echo "  IP: ${MASTER_IP}"
echo "  Porta: ${MASTER_PORT}"
echo "  Daemon: ${DAEMON}"
echo ""

if [ "$DAEMON" = true ]; then
    ./sloth-runner master start \
        --port "${MASTER_PORT}" \
        --bind-address "${MASTER_IP}" \
        --daemon
    
    sleep 2
    
    # Verificar se iniciou
    if pgrep -f "sloth-runner master" > /dev/null; then
        echo -e "${GREEN}✅ Master server iniciado com sucesso em modo daemon!${NC}"
        echo ""
        echo "Para verificar agents registrados:"
        echo "  ./sloth-runner agent list --master ${MASTER_IP}:${MASTER_PORT}"
        echo ""
        echo "Para parar o master:"
        echo "  pkill -f 'sloth-runner master'"
    else
        echo -e "${RED}❌ Falha ao iniciar master server${NC}"
        exit 1
    fi
else
    echo -e "${BLUE}🔍 Iniciando em modo foreground (pressione Ctrl+C para parar)${NC}"
    echo ""
    ./sloth-runner master start \
        --port "${MASTER_PORT}" \
        --bind-address "${MASTER_IP}"
fi

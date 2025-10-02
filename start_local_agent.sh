#!/bin/bash
# Script para iniciar um agent local do Sloth Runner

set -e

# Configuração
AGENT_NAME="${1:-local-agent}"
AGENT_IP="${2:-192.168.1.29}"
AGENT_PORT="${3:-50051}"
MASTER_IP="192.168.1.29"
MASTER_PORT="50053"
DAEMON=true

# Cores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}🤖 Iniciando Sloth Runner Agent${NC}"
echo ""

# Verificar argumentos
if [ -z "$1" ]; then
    echo -e "${YELLOW}ℹ️  Uso:${NC}"
    echo "  $0 <agent-name> [agent-ip] [agent-port]"
    echo ""
    echo "Exemplos:"
    echo "  $0 local-agent                           # Agent local na porta 50051"
    echo "  $0 test-agent 192.168.1.29 50052        # Agent local na porta 50052"
    echo ""
    echo -e "${BLUE}Usando valores padrão...${NC}"
    echo ""
fi

# Verificar se já está rodando
if pgrep -f "sloth-runner agent start.*--name ${AGENT_NAME}" > /dev/null; then
    echo -e "${RED}⚠️  Agent '${AGENT_NAME}' já está rodando!${NC}"
    echo ""
    echo "Para parar este agent, execute:"
    echo "  pkill -f 'sloth-runner agent start.*--name ${AGENT_NAME}'"
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

# Verificar se o master está rodando
echo -e "${BLUE}🔍 Verificando master server...${NC}"
if ! nc -z "${MASTER_IP}" "${MASTER_PORT}" 2>/dev/null; then
    echo -e "${RED}❌ Master server não está acessível em ${MASTER_IP}:${MASTER_PORT}${NC}"
    echo ""
    echo "Inicie o master server primeiro:"
    echo "  ./start_master.sh"
    echo ""
    exit 1
fi
echo -e "${GREEN}✅ Master server está acessível${NC}"
echo ""

# Iniciar agent
echo -e "${BLUE}📡 Configuração:${NC}"
echo "  Nome: ${AGENT_NAME}"
echo "  IP: ${AGENT_IP}"
echo "  Porta: ${AGENT_PORT}"
echo "  Master: ${MASTER_IP}:${MASTER_PORT}"
echo "  Daemon: ${DAEMON}"
echo ""

if [ "$DAEMON" = true ]; then
    ./sloth-runner agent start \
        --name "${AGENT_NAME}" \
        --port "${AGENT_PORT}" \
        --master "${MASTER_IP}:${MASTER_PORT}" \
        --bind-address "${AGENT_IP}" \
        --daemon
    
    sleep 2
    
    # Verificar se iniciou
    if pgrep -f "sloth-runner agent start.*--name ${AGENT_NAME}" > /dev/null; then
        echo -e "${GREEN}✅ Agent '${AGENT_NAME}' iniciado com sucesso!${NC}"
        echo ""
        echo "Para verificar agents registrados:"
        echo "  ./sloth-runner agent list --master ${MASTER_IP}:${MASTER_PORT}"
        echo ""
        echo "Para parar este agent:"
        echo "  pkill -f 'sloth-runner agent start.*--name ${AGENT_NAME}'"
    else
        echo -e "${RED}❌ Falha ao iniciar agent${NC}"
        exit 1
    fi
else
    echo -e "${BLUE}🔍 Iniciando em modo foreground (pressione Ctrl+C para parar)${NC}"
    echo ""
    ./sloth-runner agent start \
        --name "${AGENT_NAME}" \
        --port "${AGENT_PORT}" \
        --master "${MASTER_IP}:${MASTER_PORT}" \
        --bind-address "${AGENT_IP}"
fi

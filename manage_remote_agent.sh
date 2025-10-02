#!/bin/bash
# Script helper para gerenciar agents remotos do Sloth Runner via SSH

set -e

# Configuração
MASTER_IP="192.168.1.29"
MASTER_PORT="50053"
REMOTE_BIN_PATH="$HOME/.local/bin"

# Cores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Função para mostrar uso
show_usage() {
    echo -e "${BLUE}🤖 Sloth Runner Remote Agent Manager${NC}"
    echo ""
    echo "Uso: $0 <comando> [argumentos]"
    echo ""
    echo "Comandos:"
    echo "  start <user@host> <agent-name> <agent-ip> [port]"
    echo "      Inicia um agent remoto via SSH"
    echo ""
    echo "  stop <user@host> <agent-name>"
    echo "      Para um agent remoto via SSH"
    echo ""
    echo "  status <user@host> <agent-name>"
    echo "      Verifica status de um agent remoto"
    echo ""
    echo "  install <user@host>"
    echo "      Instala sloth-runner no host remoto"
    echo ""
    echo "Exemplos:"
    echo "  $0 start chalkan3@192.168.1.16 ladyguica 192.168.1.16"
    echo "  $0 stop chalkan3@192.168.1.16 ladyguica"
    echo "  $0 status chalkan3@192.168.1.16 ladyguica"
    echo "  $0 install chalkan3@192.168.1.16"
    echo ""
}

# Função para iniciar agent remoto
start_remote_agent() {
    local remote_host="$1"
    local agent_name="$2"
    local agent_ip="$3"
    local agent_port="${4:-50051}"
    
    if [ -z "$remote_host" ] || [ -z "$agent_name" ] || [ -z "$agent_ip" ]; then
        echo -e "${RED}❌ Argumentos insuficientes${NC}"
        echo ""
        show_usage
        exit 1
    fi
    
    echo -e "${BLUE}🚀 Iniciando agent '${agent_name}' em ${remote_host}${NC}"
    echo ""
    echo -e "${BLUE}📡 Configuração:${NC}"
    echo "  Host: ${remote_host}"
    echo "  Nome: ${agent_name}"
    echo "  IP: ${agent_ip}"
    echo "  Porta: ${agent_port}"
    echo "  Master: ${MASTER_IP}:${MASTER_PORT}"
    echo ""
    
    # Verificar se o host está acessível
    if ! ping -c 1 "${agent_ip}" &> /dev/null; then
        echo -e "${RED}❌ Host ${agent_ip} não está acessível${NC}"
        exit 1
    fi
    
    # Iniciar agent via SSH
    ssh "${remote_host}" "cd ${REMOTE_BIN_PATH} && ./sloth-runner agent start \
        --name '${agent_name}' \
        --port ${agent_port} \
        --master '${MASTER_IP}:${MASTER_PORT}' \
        --bind-address '${agent_ip}' \
        --daemon"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Agent '${agent_name}' iniciado com sucesso!${NC}"
        echo ""
        echo "Para verificar:"
        echo "  ./sloth-runner agent list --master ${MASTER_IP}:${MASTER_PORT}"
    else
        echo -e "${RED}❌ Falha ao iniciar agent${NC}"
        exit 1
    fi
}

# Função para parar agent remoto
stop_remote_agent() {
    local remote_host="$1"
    local agent_name="$2"
    
    if [ -z "$remote_host" ] || [ -z "$agent_name" ]; then
        echo -e "${RED}❌ Argumentos insuficientes${NC}"
        echo ""
        show_usage
        exit 1
    fi
    
    echo -e "${BLUE}🛑 Parando agent '${agent_name}' em ${remote_host}${NC}"
    echo ""
    
    # Parar agent via SSH
    ssh "${remote_host}" "pkill -f 'sloth-runner agent start.*--name ${agent_name}'"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Agent '${agent_name}' parado com sucesso!${NC}"
    else
        echo -e "${YELLOW}⚠️  Agent pode não estar rodando${NC}"
    fi
}

# Função para verificar status de agent remoto
status_remote_agent() {
    local remote_host="$1"
    local agent_name="$2"
    
    if [ -z "$remote_host" ] || [ -z "$agent_name" ]; then
        echo -e "${RED}❌ Argumentos insuficientes${NC}"
        echo ""
        show_usage
        exit 1
    fi
    
    echo -e "${BLUE}🔍 Verificando status do agent '${agent_name}' em ${remote_host}${NC}"
    echo ""
    
    # Verificar status via SSH
    ssh "${remote_host}" "pgrep -f 'sloth-runner agent start.*--name ${agent_name}' > /dev/null"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Agent '${agent_name}' está rodando${NC}"
        
        # Mostrar processo
        echo ""
        echo "Processo:"
        ssh "${remote_host}" "ps aux | grep 'sloth-runner agent start.*--name ${agent_name}' | grep -v grep"
    else
        echo -e "${RED}❌ Agent '${agent_name}' não está rodando${NC}"
    fi
}

# Função para instalar sloth-runner no host remoto
install_remote() {
    local remote_host="$1"
    
    if [ -z "$remote_host" ]; then
        echo -e "${RED}❌ Host remoto não especificado${NC}"
        echo ""
        show_usage
        exit 1
    fi
    
    echo -e "${BLUE}📦 Instalando sloth-runner em ${remote_host}${NC}"
    echo ""
    
    # Copiar script de instalação e executar
    if [ -f "./install.sh" ]; then
        echo "Copiando script de instalação..."
        scp install.sh "${remote_host}:/tmp/"
        
        echo "Executando instalação remota..."
        ssh "${remote_host}" "bash /tmp/install.sh"
        
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✅ Instalação concluída com sucesso!${NC}"
        else
            echo -e "${RED}❌ Falha na instalação${NC}"
            exit 1
        fi
    else
        echo -e "${RED}❌ Script install.sh não encontrado${NC}"
        exit 1
    fi
}

# Processar comando
COMMAND="$1"

case "$COMMAND" in
    start)
        shift
        start_remote_agent "$@"
        ;;
    stop)
        shift
        stop_remote_agent "$@"
        ;;
    status)
        shift
        status_remote_agent "$@"
        ;;
    install)
        shift
        install_remote "$@"
        ;;
    *)
        show_usage
        exit 1
        ;;
esac

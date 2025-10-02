#!/bin/bash
# 🎬 DEMO: Execução Remota via CMD no Sloth Runner
# Este script demonstra a solução funcional de execução remota

set -e

MASTER="192.168.1.29:50053"
COLOR_GREEN='\033[0;32m'
COLOR_BLUE='\033[0;34m'
COLOR_YELLOW='\033[1;33m'
COLOR_RESET='\033[0m'

echo "═══════════════════════════════════════════════════════════════════"
echo "    🎬 DEMONSTRAÇÃO: Execução Remota via CMD                       "
echo "═══════════════════════════════════════════════════════════════════"
echo ""

# Função para mostrar cabeçalho
show_header() {
    echo ""
    echo -e "${COLOR_BLUE}═══════════════════════════════════════════════════════════════════${COLOR_RESET}"
    echo -e "${COLOR_BLUE}    $1${COLOR_RESET}"
    echo -e "${COLOR_BLUE}═══════════════════════════════════════════════════════════════════${COLOR_RESET}"
    echo ""
}

# Função para mostrar sucesso
show_success() {
    echo -e "${COLOR_GREEN}✅ $1${COLOR_RESET}"
}

# Função para mostrar info
show_info() {
    echo -e "${COLOR_YELLOW}ℹ️  $1${COLOR_RESET}"
}

# Passo 1: Verificar agentes
show_header "Passo 1: Verificando Agentes Disponíveis"
show_info "Executando: ./sloth-runner agent list --master $MASTER"
echo ""
./sloth-runner agent list --master $MASTER
show_success "Agentes verificados!"

# Aguardar
echo ""
read -p "Pressione ENTER para continuar com o Exemplo 1 (Hello World)..."

# Passo 2: Exemplo 1 - Hello World
show_header "Passo 2: Exemplo 1 - Hello World Remoto"
show_info "Executando: ./sloth-runner run -f examples/agents/hello_remote_cmd.sloth hello_remote"
echo ""
./sloth-runner run -f examples/agents/hello_remote_cmd.sloth hello_remote 2>&1 | grep -E "✅|📋|Execution Summary|hello|Task execution completed"
show_success "Exemplo 1 concluído!"

# Aguardar
echo ""
read -p "Pressione ENTER para continuar com o Exemplo 2 (Funcional Completo)..."

# Passo 3: Exemplo 2 - Funcional
show_header "Passo 3: Exemplo 2 - Exemplo Funcional Completo"
show_info "Executando: ./sloth-runner run -f examples/agents/functional_cmd_example.sloth remote_via_cmd"
echo ""
./sloth-runner run -f examples/agents/functional_cmd_example.sloth remote_via_cmd 2>&1 | grep -E "✅|📋|Execution Summary|Task execution completed"
show_success "Exemplo 2 concluído!"

# Aguardar
echo ""
read -p "Pressione ENTER para continuar com o Exemplo 3 (Pipeline Completo)..."

# Passo 4: Exemplo 3 - Pipeline Completo
show_header "Passo 4: Exemplo 3 - Pipeline de Infraestrutura Completo"
show_info "Executando: ./sloth-runner run -f examples/agents/complete_infrastructure_check.sloth distributed_infrastructure_check"
echo ""
./sloth-runner run -f examples/agents/complete_infrastructure_check.sloth distributed_infrastructure_check 2>&1 | grep -E "✅|===|RELATÓRIO|🎉|Execution Summary|Task execution completed"
show_success "Exemplo 3 concluído!"

# Finalização
show_header "🎉 DEMONSTRAÇÃO CONCLUÍDA"
echo ""
show_success "Todos os exemplos executados com sucesso!"
echo ""
show_info "Documentação disponível em:"
echo "  • examples/agents/QUICK_START.md"
echo "  • examples/agents/README_CMD_FUNCIONAL.md"
echo "  • examples/agents/INDEX.md"
echo ""
echo "═══════════════════════════════════════════════════════════════════"
echo ""

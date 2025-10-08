#!/bin/bash

# Watcher System Test Suite Runner
# Executa todos os testes de watchers sequencialmente

set -e

SLOTH_RUNNER="${HOME}/.local/bin/sloth-runner"
TEST_DIR="$(dirname "$0")"
LOG_DIR="/tmp/watcher_tests"

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Criar diretório de logs
mkdir -p "$LOG_DIR"

echo ""
echo "╔═══════════════════════════════════════════════════════╗"
echo "║     🧪 Watcher System Test Suite                     ║"
echo "║     Testing: Watchers, Events, Hooks                 ║"
echo "╚═══════════════════════════════════════════════════════╝"
echo ""

# Array de testes
TESTS=(
    "01_file_watcher_test.sloth:File Watcher (Create/Change/Delete)"
    "02_cpu_watcher_test.sloth:CPU Watcher (Threshold Detection)"
    "03_memory_watcher_test.sloth:Memory Watcher (Threshold Detection)"
    "04_process_watcher_test.sloth:Process Watcher (Start/Stop)"
    "05_port_watcher_test.sloth:Port Watcher (Listen Detection)"
    "06_watcher_with_hook_test.sloth:Watcher + Hook Integration"
    "07_complete_e2e_test.sloth:Complete End-to-End Test"
)

PASSED=0
FAILED=0
TOTAL=${#TESTS[@]}

# Função para executar um teste
run_test() {
    local test_file=$1
    local test_name=$2
    local test_num=$3

    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}Test ${test_num}/${TOTAL}: ${test_name}${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""

    local workflow_name=$(basename "$test_file" .sloth)
    local log_file="${LOG_DIR}/${workflow_name}.log"

    # Executar teste
    if $SLOTH_RUNNER run "$workflow_name" --file "${TEST_DIR}/${test_file}" --yes > "$log_file" 2>&1; then
        echo -e "${GREEN}✅ PASSED${NC}"
        PASSED=$((PASSED + 1))

        # Mostrar resumo do log
        echo ""
        echo "Last 10 lines of output:"
        tail -10 "$log_file" | sed 's/^/  /'
    else
        echo -e "${RED}❌ FAILED${NC}"
        FAILED=$((FAILED + 1))

        # Mostrar erro
        echo ""
        echo "Error log:"
        tail -20 "$log_file" | sed 's/^/  /'
    fi

    echo ""
    echo "Full log: $log_file"

    # Pausa entre testes para permitir processamento de eventos
    echo ""
    echo -e "${YELLOW}⏳ Waiting 10 seconds before next test...${NC}"
    sleep 10
}

# Executar todos os testes
for i in "${!TESTS[@]}"; do
    IFS=':' read -r test_file test_name <<< "${TESTS[$i]}"
    run_test "$test_file" "$test_name" $((i + 1))
done

# Resultados finais
echo ""
echo "╔═══════════════════════════════════════════════════════╗"
echo "║              📊 Test Results Summary                  ║"
echo "╚═══════════════════════════════════════════════════════╝"
echo ""
echo "Total tests: $TOTAL"
echo -e "${GREEN}Passed: $PASSED${NC}"
if [ $FAILED -gt 0 ]; then
    echo -e "${RED}Failed: $FAILED${NC}"
else
    echo "Failed: 0"
fi
echo ""

# Informações adicionais
echo "📁 Test logs location: $LOG_DIR"
echo "🔍 Agent logs: ssh chalkan3@192.168.1.16 'cat agent.log'"
echo "📊 Event database: ssh chalkan3@192.168.1.16 'sqlite3 /etc/sloth-runner/events.db'"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✅ ALL TESTS PASSED!${NC}"
    echo "   Watcher system is working correctly"
    exit 0
else
    echo -e "${RED}❌ SOME TESTS FAILED${NC}"
    echo "   Review logs for details"
    exit 1
fi

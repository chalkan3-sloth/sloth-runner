#!/bin/bash

echo "🔧 DEMONSTRAÇÃO: exec.run vs io.popen"
echo "===================================="
echo ""

echo "📋 1. Listando exemplos disponíveis:"
echo "---"
ls -la examples/agents/*.sloth | grep -E "(exec|popen|delegate)" | awk '{print "• " $9 " (" $5 " bytes)"}'
echo ""

echo "🚀 2. Testando exec.run method:"
echo "---"
sloth-runner run -f examples/agents/exec_vs_popen.sloth simple_exec_test 2>&1 | grep -E "(exec.run|Host:|Directory listing|SUCCESS|FAILED)" | head -10
echo ""

echo "📊 3. Comparação dos métodos:"
echo "---"
echo "✅ exec.run (RECOMENDADO):"
echo "   • Integração nativa com sloth-runner"
echo "   • Controle de ambiente e working directory"
echo "   • Tratamento estruturado de erros"
echo "   • Separação clara de stdout/stderr"
echo "   • Logs automáticos no sistema"
echo ""
echo "⚠️  io.popen (MÉTODO ANTIGO):"
echo "   • Funcional mas limitado"
echo "   • Menos controle sobre execução"
echo "   • Tratamento manual de erros"
echo "   • Sem logs estruturados"
echo ""

echo "💡 4. Sintaxe Recomendada:"
echo "---"
cat << 'EOF'
-- ✅ RECOMENDADO: Usar exec.run
local stdout, stderr, failed = exec.run("ls -la $HOME")
if not failed then
    log.info("Sucesso: " .. stdout)
else
    log.error("Erro: " .. stderr)
end

-- ❌ EVITAR: io.popen (método antigo)
local handle = io.popen("ls -la $HOME")
local output = handle:read("*a")
local success = handle:close()
EOF

echo ""
echo "🎯 5. Arquivos de Exemplo Atualizados:"
echo "---"
echo "• exec_vs_popen.sloth          → Comparação dos métodos"
echo "• ls_delegate_exec.sloth       → Usando exec.run com delegate_to"
echo "• ls_delegate_simple.sloth     → Versão com io.popen (para comparação)"
echo ""

echo "✅ RECOMENDAÇÃO FINAL:"
echo "Use sempre exec.run em vez de io.popen para melhor integração!"
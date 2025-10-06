#!/bin/bash
# Script para ajudar na extração de comandos do main.go para estrutura modular
# Uso: ./scripts/extract-command.sh COMMAND_NAME PARENT_COMMAND

set -e

COMMAND_NAME="$1"
PARENT_COMMAND="$2"

if [ -z "$COMMAND_NAME" ]; then
    echo "Uso: $0 COMMAND_NAME [PARENT_COMMAND]"
    echo "Exemplo: $0 list agent"
    echo "Exemplo: $0 version"
    exit 1
fi

# Diretório de destino
if [ -z "$PARENT_COMMAND" ]; then
    DIR="cmd/sloth-runner/commands"
    PACKAGE="commands"
else
    DIR="cmd/sloth-runner/commands/${PARENT_COMMAND}"
    PACKAGE="${PARENT_COMMAND}"
fi

# Criar diretório se não existir
mkdir -p "$DIR"

FILE="${DIR}/${COMMAND_NAME}.go"

if [ -f "$FILE" ]; then
    echo "❌ Arquivo já existe: $FILE"
    exit 1
fi

# Template básico
cat > "$FILE" << 'EOF'
package PACKAGE

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewCOMMAND_NAMECommand creates the COMMAND_NAME command
func NewCOMMAND_NAMECommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "COMMAND_NAME",
		Short: "TODO: Add short description",
		Long:  `TODO: Add long description`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement command logic
			// 1. Extract flags
			// 2. Create service if needed
			// 3. Execute operation
			// 4. Display result
			return nil
		},
	}
}
EOF

# Substituir placeholders
sed -i '' "s/PACKAGE/${PACKAGE}/g" "$FILE"
sed -i '' "s/COMMAND_NAME/${COMMAND_NAME}/g" "$FILE"

# Capitalizar primeira letra para função
COMMAND_CAPITALIZED="$(tr '[:lower:]' '[:upper:]' <<< ${COMMAND_NAME:0:1})${COMMAND_NAME:1}"
sed -i '' "s/NewCOMMAND_NAME/New${COMMAND_CAPITALIZED}/g" "$FILE"

echo "✅ Comando criado: $FILE"
echo ""
echo "Próximos passos:"
echo "1. Encontre o comando no main.go:"
echo "   grep -A 50 'var ${COMMAND_NAME}.*Cmd = &cobra.Command' cmd/sloth-runner/main.go"
echo ""
echo "2. Copie a lógica do RunE para o novo arquivo"
echo ""
echo "3. Adicione o comando ao parent:"
if [ -z "$PARENT_COMMAND" ]; then
    echo "   Edite cmd/sloth-runner/commands/root.go"
else
    echo "   Edite cmd/sloth-runner/commands/${PARENT_COMMAND}/${PARENT_COMMAND}.go"
fi
echo ""
echo "4. Teste o comando:"
echo "   go build -o sloth-runner-test ./cmd/sloth-runner"
echo "   ./sloth-runner-test ${COMMAND_NAME} --help"
